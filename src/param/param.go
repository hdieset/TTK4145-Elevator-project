package param
//constants for the elevator project 

const (
	PEERPORT 	= 20009
	BCASTPORT 	= 30009	
	N_FLOORS 	= 4
	N_BUTTONS	= 3
	PANELPORT 	= "localhost:15657"
)

type SyncArray struct {
	CurrentFloor 	[]int 
	Melding 		string
	ErDetFredag 	bool
	MyID 			string
	Suicide			bool
	Iter 			int
}

type Dirn int 
const (
	D_Down = Dirn(-1)
	D_Stop = Dirn(0)
	D_Up   = Dirn(1)
)

type ButtonType int 
const (
	B_HallUp 	= ButtonType(0)
	B_HallDown  = ButtonType(1)
	B_Cab 		= ButtonType(2)
)

type ButtonEvent struct {
	Floor  int
	Button ButtonType
}

type ElevatorBehaviour int
const(
	EB_Idle		= iota   // Just set arbitrary values for these
	EB_DoorOpen
	EB_Moving
)

type Elevator struct {
	Floor 		int
	Direction	Dirn
	Requests	[N_FLOORS][N_BUTTONS]int
	Behaviour 	ElevatorBehaviour
	DoorOpenDuration_s float64
}

/*	
-IP-adresses	
-	
-Plass 1 	: 129.241.187.***	
-Plass 2 	: 129.241.187.149	
-Plass 3 	: 129.241.187.150	
-Plass 4 	: 129.241.187.***	
-Plass 5 	: 129.241.187.***	
-Plass 6 	: 129.241.187.146	
-Plass 7 	: 129.241.187.***	
-Plass 8 	: 129.241.187.161	
-Plass 9 	: 129.241.187.156	
-Plass 10 	: 129.241.187.***	
-Plass 11 	: 129.241.187.***	
-Plass 12 	: 129.241.187.***	
-Plass 13 	: 129.241.187.***	
-Plass 14 	: 129.241.187.***	
-Plass 15 	: 129.241.187.***	
-Plass 16 	: 129.241.187.***	
-Plass 17 	: 129.241.187.***	
-Plass 18 	: 129.241.187.***	
-Plass 19 	: 129.241.187.***	
-Plass 20 	: 129.241.187.***	
-*/