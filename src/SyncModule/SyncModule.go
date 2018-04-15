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
	localSyncArray := initLocalSyncArray()
	ticker := time.NewTicker(3000*time.Millisecond)
	time.Sleep(500*time.Millisecond) //la til denne 
	for {
		select { 
		case newPeerList := <- peerUpdateCh:
			fmt.Println("Case: new peer list")
			// Deleting lost Elevators from local sync array
			//fjerne alle ackd fra LocalSyncArray hvis peer er lost

			for iter := 0; iter < len(newPeerList.Lost); iter++ {
				lostId := newPeerList.Lost[iter]
				if remIndex := Sync_ElevIndexFinder(localSyncArray,lostId); remIndex != -1 {
					localSyncArray.AllElevators = remove(localSyncArray.AllElevators, remIndex)
				}

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

				tempSyncArray := initLocalSyncArray()
				localSyncArray.AckHallStates = tempSyncArray.AckHallStates 
			} else {
				isAlone = false
			}

		case recievedSyncArray := <- networkRx:
			if initialized { //kan dette løses med peerTxenable? ELLER KAN VI FLYTTE  ?? 
				t := time.Now()
				fmt.Println(t.String()," - Case: mottok sync array paa net")
				fmt.Println("******************************************************")
				fmt.Println(recievedSyncArray)
				fmt.Println("******************************************************\n")
				//Update/overwrite the sender's elevator struct in localSyncArray 

				
				//the owner of a SyncArray will always have it's own Elevator struct at index 0 in .AllElevators
				if index := Sync_ElevIndexFinder(localSyncArray,recievedSyncArray.AllElevators[0].Id); index != -1 {
					localSyncArray.AllElevators[index] = recievedSyncArray.AllElevators[0]
				} else {
					localSyncArray.AllElevators = append(localSyncArray.AllElevators, recievedSyncArray.AllElevators[0])
				}

				localSyncArray = updateHallStates(recievedSyncArray, localSyncArray, localElevatorID) 

				sendSyncArrayToCost <- localSyncArray 
			}

		case newBtnEvent := <- receivedButtonPress: 
			if initialized {
				fmt.Println("case: mottok knapp")
				localSyncArray = addOrders(newBtnEvent, localSyncArray, localElevatorID, isAlone)
				//add new cab request to local Elevator in localSyncArray
				//fmt.Println("knapp lagt til i sync array")
				//networkTx <- localSyncArray
				fmt.Println("Sync array som sendes til cost: ")
				printSyncArray(localSyncArray)
				sendSyncArrayToCost <- localSyncArray
			}

		case recievedLocalElev := <- localElevatorCh: 
			
			fmt.Println("Case: mottok local Elevator")
			
			//KAN VI FLYTTE DENNE TIL OVER FOR SELECTEN??????????????????????????????????????????????????????????
			if !initialized { 
				//localSyncArray.AllElevators = append(localSyncArray.AllElevators, recievedLocalElev)
				localSyncArray.AllElevators[0] = recievedLocalElev
				localSyncArray.AllElevators[0].Id = localElevatorID
				initialized = true
			} 
			// local elevator is always at index 0 of localSyncArray.AllElevators
			localSyncArray.AllElevators[0].Floor     = recievedLocalElev.Floor
			localSyncArray.AllElevators[0].Direction = recievedLocalElev.Direction
			localSyncArray.AllElevators[0].Behaviour = recievedLocalElev.Behaviour

				// only send recieved elevator to cost if an order is completed - in order to disable hall lights.
			temp := localSyncArray
			localSyncArray = completeOrders(recievedLocalElev, localSyncArray, isAlone)
			if temp.HallStates != localSyncArray.HallStates {
				sendSyncArrayToCost <- localSyncArray
			}

		case <-ticker.C:
			if initialized {
				t := time.Now()
				fmt.Println(t.String()," - Case: Ticker, send localSyncArray to other peers:")
				printSyncArray(localSyncArray)
				networkTx <- localSyncArray 
			}
		}
		
	}

}

func initLocalSyncArray() SyncArray { //tok inn localElevID før 
	localSyncArray := SyncArray{} //var new her
	localSyncArray.AllElevators = append(localSyncArray.AllElevators, Elevator{})
	//localSyncArray.AllElevators = make(map[string]Elevator)
	//localSyncArray.AckHallStates = make(map[string][N_FLOORS][N_BUTTONS-1]bool)


	//localSyncArray.Owner = localElevatorID


	//localSyncArray.AccessAllElevators(localElevatorID) = Elevator{}

	//set all hall states to unknown in case of reboot
	for floors := 0; floors < N_FLOORS; floors++ {
		for btn := 0; btn < N_BUTTONS-1; btn++ {
			localSyncArray.HallStates[floors][btn] = Hall_unknown
			localSyncArray.AckHallStates[floors][btn] = make(map[string]bool) 
		}
	}

	return localSyncArray //var pointer her 
} 

func updateHallStates(recievedSyncArray SyncArray, localSyncArray SyncArray, localElevatorID string) SyncArray {
	for floors := 0; floors < N_FLOORS; floors++ {
		for btn := 0; btn < N_BUTTONS-1; btn++ {
			switch recievedSyncArray.HallStates[floors][btn] {
			case Hall_unknown: 
				break

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
					if len(localSyncArray.AllElevators) == len(localSyncArray.AckHallStates[floors][btn]) {
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
			}
		}
	}
	return localSyncArray 
}

func addOrders(newBtnEvent ButtonEvent, localSyncArray SyncArray, localElevatorID string, isAlone bool) SyncArray {  //dumt navn 
	if newBtnEvent.Button == B_Cab {
		localSyncArray.AllElevators[0].Requests[newBtnEvent.Floor][B_Cab] = true 
	} else if !isAlone { 
		localSyncArray.HallStates[newBtnEvent.Floor][newBtnEvent.Button] = Hall_unconfirmed //skal være unconfirmed 
		localSyncArray.AckHallStates[newBtnEvent.Floor][newBtnEvent.Button][localElevatorID] = true
	}
	return localSyncArray 
}


func completeOrders(updatedLocalElev Elevator, localSyncArray SyncArray, isAlone bool) SyncArray {  //dumt navn
	for floors := 0; floors < N_FLOORS; floors++ {
		//remove completed cab requests from local Elevator in localSyncArray
		if updatedLocalElev.CompletedReq[floors][B_Cab] == true {
			localSyncArray.AllElevators[0].Requests[floors][B_Cab] = false 
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



func Sync_ElevIndexFinder(input SyncArray, ElevatorID string) int {
	for i:= range input.AllElevators {
		if ElevatorID == input.AllElevators[i].Id {
		return i
		}
	}
	return -1 
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