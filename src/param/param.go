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
	B_HallDown  = ButtonType(0)
	B_HallUp 	= ButtonType(1)
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