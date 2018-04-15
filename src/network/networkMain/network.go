package network 

import(
	"network/bcast"
	"network/peers"
	"network/localip"
	."types"
	"fmt"
	"time"
	"os"
)

func Network(localElevatorID string,
		peerTxEnable <-chan bool,
		peerUpdateCh chan<- PeerUpdate,
		networkTx <-chan SyncArray,
		networkRx chan<- SyncArray) {

	go peers.Transmitter(PEERPORT, localElevatorID, peerTxEnable)
	go peers.Receiver(PEERPORT, peerUpdateCh)
	go bcast.Transmitter(BCASTPORT,localElevatorID, networkTx)  
	go bcast.Receiver(BCASTPORT, localElevatorID, networkRx)
}

func Network_generateID()(id string) {
	localIP, err := localip.LocalIP()
	for err != nil {
		localIP, err = localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			fmt.Println("Failed to connect. Retrying in 3 seconds...")
			time.Sleep(3 * time.Second)
			//localIP = "DISCONNECTED"
		}
	}
	
	id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid()) 
	return 
}