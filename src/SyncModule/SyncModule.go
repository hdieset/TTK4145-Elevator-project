package syncmodule

import (
	"fmt"
	."param"
)

func SyncModule (localElevatorID string, 
		peerUpdateCh <-chan PeerUpdate,
		networkRx <-chan SyncArray,
		netWorkTx chan<- SyncArray,
		sendSyncArrayToCost chan<- SyncArray,
		receivedLocalElevator <-chan Elevator, 
		receivedButtonPress <-chan ButtonEvent) {
	
	var isAlone bool

	localSyncArray := initLocalSyncArray(localElevatorID)
	
	for {
		select { 
		case newPeerList <- peerUpdateCh:
			// Deleting lost Elevators from local sync array
			for lostElevators := 0; lostElevators < len(newPeerList.Lost); lostElevators++ {
				delete(localSyncArray.AllElevators, newPeerList.Lost[lostElevators] )
			}

			//REDUNDANT????????????????????????????????????????
			// Set flag if this elevator is disconnected from all others
			if len(newPeerList.Peers) == 1 {
				tempSyncArray := initLocalSyncArray(localElevatorID)
				localSyncArray.HallStates = tempSyncArray.HallStates 
				localSyncArray.AckHallStates = tempSyncArray.AckHallStates 

				isAlone = true
			} else {
				isAlone = false
			}

		case recievedSyncArray <- networkRx:
			//Update/overwrite the sender's elevator struct in localSyncArray 
			localSyncArray.AllElevators[recievedSyncArray.Owner] = recievedSyncArray.AllElevators[recievedSyncArray.Owner]
			//KJØRE updateHallStates(recievedSyncArray, localSyncArray, localElevatorID) 
			netWorkTx <- localSyncArray
			sendSyncArrayToCost <- localSyncArray 

		case updatedLocalElev <- receivedLocalElevator: 
			localSyncArray.AllElevators[localElevatorID].Floor = updatedLocalElev.Floor
			localSyncArray.AllElevators[localElevatorID].Direction = updatedLocalElev.Direction
			localSyncArray.AllElevators[localElevatorID].Behaviour = updatedLocalElev.Behaviour

			for floors := 0; floors < N_FLOORS; floors++ {
				//remove completed cab requests from local Elevator in localSyncArray
				if updatedLocalElev.CompletedReq[floors][B_Cab] == true {
					localSyncArray.AllElevators[localElevatorID].Requests[floors][B_Cab] = false 
				}

				for hallBtn := 0; hallBtn < N_BUTTONS-1; hallBtn++ { 
				//if isAlone?? 
				// Work some magic with the hall orders. Future us has this headache (sette Hall_confirmed -> Hall_none)
				}
			}

			netWorkTx <- localSyncArray
			sendSyncArrayToCost <- localSyncArray

		case newBtnEvent <- receivedButtonPress: 
			//add new cab request to local Elevator in localSyncArray
			if newBtnEvent.Button == B_Cab {
				localSyncArray.AllElevators[localElevatorID].Requests[newBtnEvent.Floor][B_Cab] = true 
			} else {
				//SUPERFUNKSJONALITET (bump hall tingen)
			} 

			netWorkTx <- localSyncArray
			sendSyncArrayToCost <- localSyncArray
		}
	}
}

func initLocalSyncArray(localElevatorID string) SyncArray {
	localSyncArray := new(SyncArray)
	localSyncArray.AllElevators = make(map[string]Elevator)
	localSyncArray.AckHallStates = make(map[string][N_FLOORS][N_BUTTONS-1]bool)
	localSyncArray.Owner = localElevatorID	

	//set all hall states to unknown in case of reboot
	for floors := 0; floors < N_FLOORS; floors++ {
		for btn := 0; btn < N_BUTTONS-1; btn++ {
			localSyncArray.HallStates[floors][btn] = Hall_unknown
		}
	}

	return *localSyncArray 
} 

func updateHallStates(recievedSyncArray SyncArray, localSyncArray SyncArray, localElevatorID string) SyncArray {
	//hvis lengden av peerlist == 1, ikke ack fordi vi ikke vil at en enslig elevator skal ta hall orders

}

func completeHallOrders(updatedLocalElev Elevator, localSyncArray SyncArray) SyncArray {
  //husk clear ackd 
}

func addHallOrders(newBtnEvent ButtonType, localSyncArray SyncArray, localElevatorID string) SyncArray {
	//husk å legge til ackd
}







// https://blog.golang.org/go-maps-in-action
/*
m["route"] = 66
n := len(m) <- returns the number of items in the map
delete(m, "route") <- deletes index "route" of object m

*/
