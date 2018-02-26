//package main
package network

import(
	"./network/bcast"
	"./network/localip"
	"./network/peers"
	"flag"
	"fmt"
	"os"
	"time"
)


func generateId( ID string) {
	id = ID
	var id string 
	localIP, err := localip.localIP() 
	if err != nil {
		fmt.Println(err) 
		localIP = "DISCONNECTED"
	}
	id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid()) //local IP and PID in a id-string

}

func Networkloop(syncArrayRx <-chan syncArray,
	syncArrayTx chan<- syncArray,
	id){

	//setting up channels for receving updates on the id's opf the peers that are aline on the network 
	peerUpdateCh := make(chan peers.PeerUpdate)

	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTxEnable := make(chan bool)

	go peers.Transmitter(PEERPORT, id, peerTxEnable)
	go peers.Receiver(PEERPORT, peerUpdateCh)

	mapTx := make(chan syncArray)
	mapRx := make(chan syncArray)

	go bcast.Trasmitter(BCASTPORT,mapTx) 
	go bcast.Receiver(BCASTPORT,mapRx)

	/* Now, we start an indefinite for loop, transmitting if 
	   something on Tx channel, and vice versa */

	localCopy := new(syncArray)
	var dead bool = false // Usikker på om dette har blitt gjort riktig
	for {
		select {
		case newPeerList := <- peerUpdateCh: // p.Peers/New/Lost
			// Har mistet eller fått peer
			localCopy.peers = newPeerList
			syncArrayRx 	<- localCopy
		case newMap := <- mapRx: 		// q type syncArray
			// Vi får sync_array fra andre på nettet
			if !localCopy.suicide {
				newMap.myID = id 		// Set myID to id
				localCopy 	= newMap	// Updating local copy
				syncArrayRx <- newMap	// Sending to Sync Module
			}
		case newMap := <- syncArrayRx:
			// Recieved syncArray from Sync module, send.
			// Kanskje periodisk sending skal her?
			localCopy = newMap
			if newMap.suicide {
				peerTxEnable <- false
			}
			mapTx 		<- localCopy
		}
	}
}
