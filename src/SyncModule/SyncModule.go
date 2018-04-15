package syncmodule

import (
	."types"
	"time"
	"fmt"
)

func SyncModule (localElevatorID string, 
		peerUpdateCh <-chan PeerUpdate,
		networkRx <-chan SyncArray,
		networkTx chan<- SyncArray,
		sendSyncArrayToCost chan<- SyncArray,
		localElevatorCh <-chan Elevator, 
		receivedButtonPress <-chan ButtonEvent) {
	
	var isAlone bool
	var initialized bool = false
	localSyncArray := initLocalSyncArray(localElevatorID)
	ticker := time.NewTicker(200*time.Millisecond)
	time.Sleep(500*time.Millisecond) //la til denne 
	for {
		select { 
		case newPeerList := <- peerUpdateCh:
			fmt.Printf("Case: new peer list: %+v\n", newPeerList)
			// Deleting lost Elevators from local sync array
			//fjerne alle ackd fra LocalSyncArray hvis peer er lost

			for iter := 0; iter < len(newPeerList.Lost); iter++ {
				lostId := newPeerList.Lost[iter]
                delete(localSyncArray.AllElevators, lostId)

				//delete acks from disconnected elevators
				for floors := 0; floors < N_FLOORS; floors++ {
					for btn := 0; btn < N_BUTTONS-1; btn++ {
						delete(localSyncArray.AckHallStates[floors][btn], lostId)
					}
				}
			}

			if len(newPeerList.Peers) == 1 {
				// Set flag if this elevator is disconnected from all others
				isAlone = true 
				// Set non-confirmed Hallorders to unknown state.
				for floor := 0; floor < N_FLOORS; floor++ {
					for btn := 0; btn < N_BUTTONS-1; btn++ {
						if (localSyncArray.HallStates[floor][btn] == Hall_none) || (localSyncArray.HallStates[floor][btn] == Hall_unconfirmed) {
							localSyncArray.HallStates[floor][btn] = Hall_unknown
						}
					}
				}

				tempSyncArray := initLocalSyncArray(localElevatorID)
				localSyncArray.AckHallStates = tempSyncArray.AckHallStates 
			} else {
				isAlone = false
			}

		case recievedSyncArray := <- networkRx:
            /*for k, v := range recievedSyncArray.AllElevators {
                fmt.Printf("  Received Elevator %+v  :  %+v\n", k, v)
            }*/
			fmt.Printf("\n      Received sync array from: %+v\n\n", recievedSyncArray.OwnerId)
			if initialized { //kan dette løses med peerTxenable? ELLER KAN VI FLYTTE  ?? 
				//the owner of a SyncArray will always have it's own Elevator struct at index 0 in .AllElevators
				fmt.Println("Lengde på recievedSyncArray.AllElevators:", len(recievedSyncArray.AllElevators))
				
                // Add/update received elevator states
                for k, v := range recievedSyncArray.AllElevators {
                    if k != localElevatorID {
                        fmt.Println("remote key", k)
                        localSyncArray.AllElevators[k] = v
                    }
                }

				localSyncArray = updateHallStates(recievedSyncArray, localSyncArray, localElevatorID) 

				sendSyncArrayToCost <- localSyncArray 
			}

		case newBtnEvent := <- receivedButtonPress: 
			if initialized {
				localSyncArray = addOrders(newBtnEvent, localSyncArray, localElevatorID, isAlone)
				sendSyncArrayToCost <- localSyncArray
			}

		case recievedLocalElev := <- localElevatorCh: 
			//fmt.Printf("Case: new local state: %+v\n", recievedLocalElev)
        
			//KAN VI FLYTTE DENNE TIL OVER FOR SELECTEN??????????????????????????????????????????????????????????
            
            localSyncArray.AllElevators[localElevatorID] = recievedLocalElev
			if !initialized { 
				initialized = true
			}

				// only send recieved elevator to cost if an order is completed - in order to disable hall lights.
			temp := localSyncArray
			localSyncArray = completeOrders(recievedLocalElev, localSyncArray, localElevatorID, isAlone)
			if temp.HallStates != localSyncArray.HallStates {
				sendSyncArrayToCost <- localSyncArray
			}

		case <-ticker.C:
			if initialized {
                fmt.Printf("      Sending localSyncArray. This peer knows of: ")
                for k, _ := range localSyncArray.AllElevators {
                    fmt.Printf("  %+v ", k)
                }
                fmt.Printf("\n")
                
				networkTx <- localSyncArray 
				//fmt.Println(localSyncArray.HallStates)
				//fmt.Println("Lengde på AllElevators:", len(localSyncArray.AllElevators))
				//fmt.Println("Lengde på AckHallStates", len(localSyncArray.AckHallStates[1][B_HallDown]))
				//printSyncArray(localSyncArray)

			}
		}
	}
}

func initLocalSyncArray(owner string) SyncArray { 
	localSyncArray := SyncArray{} 
	localSyncArray.OwnerId = owner
	localSyncArray.AllElevators = make(map[string]Elevator)
	//set all hall states to unknown in case of reboot
	for floors := 0; floors < N_FLOORS; floors++ {
		for btn := 0; btn < N_BUTTONS-1; btn++ {
			localSyncArray.HallStates[floors][btn] = Hall_unknown
			localSyncArray.AckHallStates[floors][btn] = make(map[string]bool) 
		}
	}

	return localSyncArray  
} 

func updateHallStates(recievedSyncArray SyncArray, localSyncArray SyncArray, localElevatorID string) SyncArray {
	senderID := recievedSyncArray.OwnerId

	for floors := 0; floors < N_FLOORS; floors++ {
		for btn := 0; btn < N_BUTTONS-1; btn++ {
			if recAck := recievedSyncArray.AckHallStates[floors][btn][senderID]; recAck == true {
				localSyncArray.AckHallStates[floors][btn][senderID] = recAck
			}

			switch recievedSyncArray.HallStates[floors][btn] {
			case Hall_none: 
				if localSyncArray.HallStates[floors][btn] == Hall_confirmed {
					localSyncArray.HallStates[floors][btn] = Hall_none
					//Delete entire Ack list for that floor and button
				    localSyncArray.AckHallStates[floors][btn] = make(map[string]bool)

				} else if localSyncArray.HallStates[floors][btn] == Hall_unknown {
					localSyncArray.HallStates[floors][btn] = Hall_none
				}

			case Hall_unconfirmed: 
				if localSyncArray.HallStates[floors][btn] == Hall_unknown {
					localSyncArray.HallStates[floors][btn] = Hall_unconfirmed
					localSyncArray.AckHallStates[floors][btn][localElevatorID] = true

				} else if localSyncArray.HallStates[floors][btn] == Hall_none {
					localSyncArray.HallStates[floors][btn] = Hall_unconfirmed
					localSyncArray.AckHallStates[floors][btn][localElevatorID] = true

				} else if localSyncArray.HallStates[floors][btn] == Hall_unconfirmed {
					if allAliveAckd(acksKeys(localSyncArray.AckHallStates[floors][btn]), elevatorsKeys(localSyncArray.AllElevators)) {
						localSyncArray.HallStates[floors][btn] = Hall_confirmed
					}
				} 

			case Hall_confirmed:
				if localSyncArray.HallStates[floors][btn] == Hall_unknown {
					localSyncArray.HallStates[floors][btn] = Hall_confirmed
					localSyncArray.AckHallStates[floors][btn][localElevatorID] = true

				} else if localSyncArray.HallStates[floors][btn] == Hall_unconfirmed {
					localSyncArray.HallStates[floors][btn] = Hall_confirmed
				}

			case Hall_unknown: 
			}
		}
	}
	return localSyncArray 
}


// plsfix these two with reflect (also give me templates damnit)
func acksKeys(acks map[string]bool) []string {
    keys := make([]string, len(acks))
    i := 0
    for k := range acks {
        keys[i] = k
        i++
    }
    return keys
}
func elevatorsKeys(acks map[string]Elevator) []string {
    keys := make([]string, len(acks))
    i := 0
    for k := range acks {
        keys[i] = k
        i++
    }
    return keys
}


func allAliveAckd(acks []string, elevators []string) bool {
	if (len(acks) >= len(elevators)) && (elevators != nil) {
		result := true
		for _, i := range elevators {
			r := false
			for _, j := range acks {
				if i == j {
					r = true
				}
			}
			if !r {
				result = false
			}
		}
		return result
	}
	return false
}



func addOrders(newBtnEvent ButtonEvent, localSyncArray SyncArray, localElevatorID string, isAlone bool) SyncArray {  //dumt navn 
	if newBtnEvent.Button == B_Cab {
        e := localSyncArray.AllElevators[localElevatorID]
		e.Requests[newBtnEvent.Floor][B_Cab] = true 
        localSyncArray.AllElevators[localElevatorID] = e
	} else if !isAlone { 
		localSyncArray.HallStates[newBtnEvent.Floor][newBtnEvent.Button] = Hall_unconfirmed //skal være unconfirmed 
		localSyncArray.AckHallStates[newBtnEvent.Floor][newBtnEvent.Button][localElevatorID] = true
	}
	return localSyncArray 
}


func completeOrders(updatedLocalElev Elevator, localSyncArray SyncArray, localElevatorID string, isAlone bool) SyncArray {  //dumt navn
	for floors := 0; floors < N_FLOORS; floors++ {
		//remove completed cab requests from local Elevator in localSyncArray
		if updatedLocalElev.CompletedReq[floors][B_Cab] == true {        
            e := localSyncArray.AllElevators[localElevatorID]
            e.Requests[floors][B_Cab] = false 
            localSyncArray.AllElevators[localElevatorID] = e
		}

		//remove completed hall orders
		for btn := 0; btn < N_BUTTONS-1; btn++ { 
			if updatedLocalElev.CompletedReq[floors][btn] {
				if isAlone {
					localSyncArray.HallStates[floors][btn] = Hall_unknown
				} else {
					localSyncArray.HallStates[floors][btn] = Hall_none
				}
				
				//Delete entire Ack list for that floor and button
				localSyncArray.AckHallStates[floors][btn] = make(map[string]bool)
			} 
		}
	}
	return localSyncArray
}



func remove(e []Elevator, i int) []Elevator {
    e[len(e)-1], e[i] = e[i], e[len(e)-1]
    return e[:len(e)-1]
}

func printSyncArray(localSyncArray SyncArray) {
	fmt.Println("vvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvv")
	for i := range localSyncArray.AllElevators {
		fmt.Println("Elev8r index:",i)
		fmt.Println("Floor       :",localSyncArray.AllElevators[i].Floor)
		fmt.Println("Direction   :",localSyncArray.AllElevators[i].Direction)
		fmt.Println("Requests    :",localSyncArray.AllElevators[i].Requests)
		fmt.Println("CompletedReq:",localSyncArray.AllElevators[i].CompletedReq)
		fmt.Println("Behaviour   :",localSyncArray.AllElevators[i].Behaviour)
		fmt.Println("ID          :",localSyncArray.AllElevators[i].Id)
		fmt.Println("-------------------------------------------------------")
	}
	fmt.Println("HallStates: ")
	fmt.Println(localSyncArray.HallStates)
	fmt.Println("AckHallStates: ")
	fmt.Println(localSyncArray.AckHallStates)
	fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
}