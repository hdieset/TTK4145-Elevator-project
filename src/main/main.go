/* This is our main file for the elevator project. This file starts the
   goroutines and is in general a very nice function. Forged by the evil
   Lord Sauron in the great fires of Mount Doom, this *master* file 
   became entitled with powers to control the greedy greedy goroutines,
   blinded by their illusion of control. One main to rule them all. */
package main

import(
	"fmt"
	"runtime"
	."network/networkMain"
	."SingleElevator/SingleElevatorMain"
	//."SingleElevator/elevator"
	."SingleElevator/extPrc"
	."Cost"
	."types"
	."SyncModule"
	//"time"
)


func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) //prøver å ha denne først 

	if SIMULATOR {
		ExtPrc_changeElevatorSimPort()
	} else {
		ExtPrc_initElevatorServer()
	}

	localElevatorID := Network_generateID()
	fmt.Println("Elevator ID: ", localElevatorID)

	peerUpdateCh 		:= make(chan PeerUpdate)
	peerTxEnable 		:= make(chan bool)  
	networkTx 	 		:= make(chan SyncArray)
	networkRx    		:= make(chan SyncArray)
	syncLocalElevator 	:= make(chan Elevator)
	syncButtonPress		:= make(chan ButtonEvent)
	sendAssignedOrders 	:= make(chan AssignedOrders)
	stopButtonPressed 	:= make(chan bool)
	sendSyncArray		:= make(chan SyncArray)

	go Network(localElevatorID, peerTxEnable, peerUpdateCh, networkTx, networkRx)
	go SingleElevator(syncLocalElevator, syncButtonPress, sendAssignedOrders, stopButtonPressed, peerTxEnable) 
	go Cost(sendAssignedOrders, sendSyncArray, localElevatorID)
	go SyncModule(localElevatorID, peerUpdateCh, networkRx, networkTx, sendSyncArray, syncLocalElevator, syncButtonPress) 


	/*ticker := time.NewTicker(100*time.Millisecond)
	for {
		select {
		case <- ticker.C :
			//<- stopButtonPressed
		}
	}*/
	
	for {
		select {
		case <- stopButtonPressed:
			return 
		}
	}

	if !SIMULATOR {
		ExtPrc_exitElevatorServer()
	}
}



