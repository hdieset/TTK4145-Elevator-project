package syncmodule

import (
	"fmt"
	."types"
	"time"
)

func SyncModule (localElevatorID string, 
		peerUpdateCh <-chan PeerUpdate,
		networkRx <-chan SyncArray,
		networkTx chan<- SyncArray,
		sendSyncArrayToCost chan<- SyncArray,
		receivedLocalElevator <-chan Elevator, 
		receivedButtonPress <-chan ButtonEvent) {
	
	var isAlone bool
	localSyncArray := initLocalSyncArray(localElevatorID)
	ticker := time.NewTicker(100*time.Millisecond)

	for {
		select { 
		case newPeerList <- peerUpdateCh:
			// Deleting lost Elevators from local sync array
			//fjerne alle ackd fra LocalSyncArray hvis peer er lost

			for lostElevators := 0; lostElevators < len(newPeerList.Lost); lostElevators++ {
				delete(localSyncArray.AllElevators, newPeerList.Lost[lostElevators] )
				//delete acks from disconnected elevators
				for floors := 0; floors < N_FLOORS; floors++ {
					for btn := 0; btn < N_BUTTONS-1; btn++ {
						delete(localSyncArray.AckHallStates[floors][btn], newPeerList.Lost[lostElevators])
					}
				}
			}

			if len(newPeerList.Peers) == 1 {
				// Set non-confirmed Hallorders to unknown state.
				for floor := 0; floor < N_FLOORS; floor++ {
					for btn := 0; btn < N_BUTTONS-1; btn++ {
						if localSyncArray.HallStates[floor][btn] == Hall_none || localSyncArray.HallStates[floor][btn] == Hall_unconfirmed {
							localSyncArray.HallStates[floor][btn] = Hall_unknown
						}
					}
				}

				tempSyncArray := initLocalSyncArray(localElevatorID)
				localSyncArray.AckHallStates = tempSyncArray.AckHallStates 

				// Set flag if this elevator is disconnected from all others
				isAlone = true
			} else {
				isAlone = false
			}

		case recievedSyncArray <- networkRx:
			//Update/overwrite the sender's elevator struct in localSyncArray 
			localSyncArray.AllElevators[recievedSyncArray.Owner] = recievedSyncArray.AllElevators[recievedSyncArray.Owner]
			localSyncArray = updateHallStates(recievedSyncArray, localSyncArray, localElevatorID) 
			//networkTx <- localSyncArray
			sendSyncArrayToCost <- localSyncArray 

		case recievedLocalElev <- localElevatorCh: 
			localSyncArray.AllElevators[localElevatorID].Floor = recievedLocalElev.Floor
			localSyncArray.AllElevators[localElevatorID].Direction = recievedLocalElev.Direction
			localSyncArray.AllElevators[localElevatorID].Behaviour = recievedLocalElev.Behaviour

			localSyncArray = completeHallOrders(updatedLocalElev, localSyncArray, isAlone)

			//networkTx <- localSyncArray
			sendSyncArrayToCost <- localSyncArray

		case newBtnEvent <- receivedButtonPress: 
			localSyncArray = addHallOrders(newBtnEvent, localSyncArray, localElevatorID, isAlone)
			//add new cab request to local Elevator in localSyncArray

			//networkTx <- localSyncArray
			sendSyncArrayToCost <- localSyncArray

		case <-ticker.C:
			networkTx <- localSyncArray
		}
	}
}

func initLocalSyncArray(localElevatorID string) SyncArray {
	localSyncArray := SyncArray{} //var new her 
	localSyncArray.AllElevators = make(map[string]Elevator)
	//localSyncArray.AckHallStates = make(map[string][N_FLOORS][N_BUTTONS-1]bool)
	localSyncArray.Owner = localElevatorID	

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
	return localSyncArray // vAR POINTERH ER 
}


func completeHallOrders(updatedLocalElev Elevator, localSyncArray SyncArray, localElevatorID string, isAlone bool) SyncArray {
	for floors := 0; floors < N_FLOORS; floors++ {
		//remove completed cab requests from local Elevator in localSyncArray
		if recievedLocalElev.CompletedReq[floors][B_Cab] == true {
			localSyncArray.AllElevators[localElevatorID].Requests[floors][B_Cab] = false 
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

func addHallOrders(newBtnEvent ButtonEvent, localSyncArray SyncArray, localElevatorID string, isAlone bool) SyncArray {
	if newBtnEvent.Button == B_Cab {
		localSyncArray.AllElevators[localElevatorID].Requests[newBtnEvent.Floor][B_Cab] = true 
	} else if !isAlone {
		localSyncArray.HallStates[newBtnEvent.Floor][newBtnEvent.Button] = Hall_unconfirmed
		localSyncArray.AckHallStates[newBtnEvent.Floor][newBtnEvent.Button][localElevatorID] = true
	}
	return localSyncArray //mest sann ikke peker her 
}

