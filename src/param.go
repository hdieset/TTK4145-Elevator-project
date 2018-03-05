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



/*
IP-adresses

Plass 1 	: 129.241.187.***
Plass 2 	: 129.241.187.149
Plass 3 	: 129.241.187.150
Plass 4 	: 129.241.187.***
Plass 5 	: 129.241.187.***
Plass 6 	: 129.241.187.146
Plass 7 	: 129.241.187.***
Plass 8 	: 129.241.187.161
Plass 9 	: 129.241.187.156
Plass 10 	: 129.241.187.***
Plass 11 	: 129.241.187.***
Plass 12 	: 129.241.187.***
Plass 13 	: 129.241.187.***
Plass 14 	: 129.241.187.***
Plass 15 	: 129.241.187.***
Plass 16 	: 129.241.187.***
Plass 17 	: 129.241.187.***
Plass 18 	: 129.241.187.***
Plass 19 	: 129.241.187.***
Plass 20 	: 129.241.187.***
*/
