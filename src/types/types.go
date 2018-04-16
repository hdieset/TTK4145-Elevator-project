package types 

var Panelport string = "localhost:15657"

const (
	PEERPORT 		  = 20258
	BCASTPORT 		  = 30258	
	N_FLOORS 		  = 4
	N_BUTTONS		  = 3
	DOOROPENDURATION  = float64(3.0) 
	MAXTRAVELDURATION = float64(3.1)
)

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

type ElevatorBehaviour int
const(
	EB_Idle		= iota 
	EB_DoorOpen
	EB_Moving
)

type HallReqStates int 
const (
	Hall_unknown 	 = HallReqStates(-1)
	Hall_none 		 = HallReqStates(0)
	Hall_unconfirmed = HallReqStates(1) 
	Hall_confirmed 	 = HallReqStates(2) 
)

type AssignedOrders struct {
	GlobalHallReq [N_FLOORS][N_BUTTONS-1]bool
	Local 		  [N_FLOORS][N_BUTTONS]bool
}

type PeerUpdate struct {
	Peers []string
	New   string
	Lost  []string
}

type ButtonEvent struct {
	Floor  int
	Button ButtonType
}

type Elevator struct {
	Floor 		 int
	Direction	 Dirn
	Requests	 [N_FLOORS][N_BUTTONS]bool
	CompletedReq [N_FLOORS][N_BUTTONS]bool
	Behaviour 	 ElevatorBehaviour
	DoorOpenDuration_s float64
}

type SyncArray struct {
    OwnerId         string
    AllElevators    map[string]Elevator
    HallStates      [N_FLOORS][N_BUTTONS-1]HallReqStates
    AckHallStates   [N_FLOORS][N_BUTTONS-1]map[string]bool
}

type AssignerCompatibleElev struct {
    Behaviour    string         `json:"behaviour"`
    Floor        int            `json:"floor"`
    Direction    string         `json:"direction"`
    CabRequests  [N_FLOORS]bool `json:"cabRequests"`
}

type AssignerCompatibleInput struct {
    HallRequests [N_FLOORS][N_BUTTONS-1]bool       	`json:"hallRequests"`
    States       map[string]*AssignerCompatibleElev `json:"states"`
}