/* This is our main file for the elevator project. This file starts the
   goroutines and is in general a very nice function. Forged by the evil
   Lord Sauron in the great fires of Mount Doom, this *master* file 
   became entitled with powers to control the greedy greedy goroutines,
   blinded by their illusion of control. One main to rule them all. */
package main

import(
   	"network"
	"flag"
	"fmt"
	"os"
	"time"
)

//Just testing somthing to send
type PeerUpdate struct {
	Peers []string
	New   string
	Lost  []string
}
type syncArray struct {
	currentFloor 	[]int 
	melding 		string
	erDetFredag 	bool
	myID 			string
	peers 			PeerUpdate
	suicide			bool
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// Initializing network goroutine
	syncArrayTx := make(chan syncArray)
	syncArrayRx := make(chan syncArray)
	go networkModule.networkMain(syncArrayRx,syncArrayTx);

	a := new(PeerUpdate)

}

// syncArrayRx <-chan syncArray => kanalen mottar fra kanalen