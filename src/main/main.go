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
	."SingleElevator/extPrc"
	."Cost"
	."types"
	."SyncModule"
)


func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) 

	if SIMULATOR {
		ExtPrc_changeElevatorSimPort()
	} else {
		ExtPrc_initElevatorServer()
	}

	localElevatorID := Network_generateID()
	fmt.Println("Elevator ID: ", localElevatorID)

	const buffers int = 100

	peerUpdateCh 		:= make(chan PeerUpdate, buffers)
	peerTxEnable 		:= make(chan bool, buffers)  
	networkTx 	 		:= make(chan SyncArray, buffers)
	networkRx    		:= make(chan SyncArray, buffers)
	syncLocalElevator 	:= make(chan Elevator, buffers)
	syncButtonPress		:= make(chan ButtonEvent, buffers)
	sendAssignedOrders 	:= make(chan AssignedOrders, buffers)
	stopButtonPressed 	:= make(chan bool, buffers)
	sendSyncArray		:= make(chan SyncArray, buffers)

	go SingleElevator(syncLocalElevator, syncButtonPress, sendAssignedOrders, stopButtonPressed, peerTxEnable) 
	go Cost(sendAssignedOrders, sendSyncArray, localElevatorID)
	go SyncModule(localElevatorID, peerUpdateCh, networkRx, networkTx, sendSyncArray, syncLocalElevator, syncButtonPress) 
	go Network(localElevatorID, peerTxEnable, peerUpdateCh, networkTx, networkRx)
	
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



