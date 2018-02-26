package param
//constants for the elevator project 

const (
	PEERPORT 	= 20009
	BCASTPORT 	= 30009	
)

type SyncArray struct {
	CurrentFloor 	[]int 
	Melding 		string
	ErDetFredag 	bool
	MyID 			string
	Suicide			bool
	Iter 			int
}