package network 

import(
	"network/bcast"
	"network/peers"
	"network/localip"
	."types"
)

func Network(localElevatorID string,
		peerTxEnable <-chan bool,
		peerUpdateCh <-chan PeerUpdate,
		networkTx chan<- SyncArray,
		networkRx <-chan SyncArray) {

	go peers.Transmitter(PEERPORT, localElevatorID, peerTxEnable)
	go peers.Receiver(PEERPORT, peerUpdateCh)
	go bcast.Transmitter(BCASTPORT,networkTx) 
	go bcast.Receiver(BCASTPORT,networkRx)
}