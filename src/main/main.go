/* This is our main file for the elevator project. This file starts the
   goroutines and is in general a very nice function. Forged by the evil
   Lord Sauron in the great fires of Mount Doom, this *master* file 
   became entitled with powers to control the greedy greedy goroutines,
   blinded by their illusion of control. One main to rule them all. */
package main

import(
   	"network/bcast"
	"network/peers"
	"network/localip"
	. "param"
	//"flag"
	"time" 
	"os"
	"fmt"
	"runtime"
)

//Just testing somthing to send
type PeerUpdate struct {
	Peers []string
	New   string
	Lost  []string
}


func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// Initializing network goroutine
	id := generateId()
	fmt.Println("Elevator id: ",id)
	//make relevant channels 
	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool) //redundant? 
	networkTx := make(chan SyncArray)
	networkRx := make(chan SyncArray)



	//init go routines 
	go peers.Transmitter(PEERPORT, id, peerTxEnable)
	go peers.Receiver(PEERPORT, peerUpdateCh)
	go bcast.Transmitter(BCASTPORT,networkTx) 
	go bcast.Receiver(BCASTPORT,networkRx)

	var testMsg SyncArray
	testMsg.Melding = "hei"
	go func() {
		for {
		testMsg.Iter++ 
		networkTx <- testMsg
		time.Sleep(3 * time.Second)
		}
	}()

	for {
		select {
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

		case a := <-networkRx:
			fmt.Printf("Received: %#v\n", a.Iter)
		}
}

	
}

// syncArrayRx <-chan syncArray => kanalen mottar fra kanalen

//generate elevator-id with local IP + PID in a string
func generateId()(id string) {
	localIP, err := localip.LocalIP() 
	if err != nil {
		fmt.Println(err) 
		localIP = "DISCONNECTED"
	}
	id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid()) //local IP and PID in a id-string
	return 
}

//.bashrc
//export GOPATH="$HOME/gruppeOgPlass9/project-gruppe-9/:"
