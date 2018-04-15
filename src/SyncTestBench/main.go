package main 

import (
	"fmt"
	."types"
	"time"
	"runtime"
	."network/networkMain"
	."SyncModule"
)


func main() {

	peerUpdateCh 		:= make(chan PeerUpdate)
	networkTx 	 		:= make(chan SyncArray)
	networkRx    		:= make(chan SyncArray)
	syncLocalElevator 	:= make(chan Elevator)
	syncButtonPress		:= make(chan ButtonEvent)
	sendSyncArray		:= make(chan SyncArray)

	runtime.GOMAXPROCS(runtime.NumCPU()) 
	localElevatorID := Network_generateID()

	fmt.Println(runtime.NumCPU())

	go SyncModule(localElevatorID, peerUpdateCh, networkRx, networkTx, sendSyncArray, syncLocalElevator, syncButtonPress)
	go tester(localElevatorID, peerUpdateCh) 
	


	for {
		select {
		case <-sendSyncArray:
			fmt.Println("TestBech: received sync array (to cost)")
		case <-networkTx: 
			fmt.Println("TestBech: received sync array (to network)") 
		}
	}
}

func tester(localElevatorID string, peerUpdateCh chan<- PeerUpdate) {
	var p PeerUpdate
	p.Peers = make([]string, 5)
	fmt.Println("Started tester")
	p.Peers[0] = localElevatorID
	peerUpdateCh <-p
	time.Sleep(1000*time.Millisecond)

	var e Elevator 
	e.Floor = 3 
	e.Direction = D_Stop
	e.Behaviour = EB_Idle 

}
