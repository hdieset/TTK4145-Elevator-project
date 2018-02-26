//constants for the elevator project 
package param

const (
	PEERPORT 	= 15647
	BCASTPORT 	= 16569	
)

type syncArray struct {
	currentFloor 	[]int 
	melding 		string
	erDetFredag 	bool
	myID 			string
	peers 			PeerUpdate
	suicide			bool
}