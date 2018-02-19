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

const (
	PEERPORT 	= 15647
	BCASTPORT 	= 16569
)

func networkMain(syncArrayRx <-chan syncArray,
	syncArrayTx chan<- syncArray ) {

	var id string 
	localIP, err := localip.localIP() 
	if err != nil {
		fmt.Println(err) 
		localIP = "DISCONNECTED"
	}
	id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid()) //local IP and PID in a id-string

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
	for {
		select {
		case p := <- peerUpdateCh: 	// p.Peers/New/Lost
			// Har mistet eller fått peer
			q.peers = p 
		case q := <- mapRx: 		// q type syncArray
			// Vi får sync_array fra andre på nettet
			q.myID = id // Set myID to id
			syncArrayTx <- q
		case q := <- syncArrayRx:
			// Recieved syncArray from Sync module, send.
			// Kanskje periodisk sending skal her?
			mapTx <- q
		}
	}
}
