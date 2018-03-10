package fsm 

import (
	//"SingleElevator/elevator_io_types"
	//"SingleElevator/elevator"
	."SingleElevator/timer"
	//"SingleElevator/elevator_io.go"
	."param"
	."SingleElevator/elevator_io"
	."SingleElevator/requests"
	"fmt"
)

var elevator Elevator 
//vil at jeg skal lag et ElevatorOutputDevice her....

//vet ikke hva jeg skal gjøre med denne enda, 
//er en __attribute__((constructor)) som bruker elevator.con osv...
func Fsm_init() {
	elevator = Elevator_uninitialized()
	elevio.Init(PANELPORT,N_FLOORS)
}

func setAllLights(es Elevator) {

	for floor := 0; floor < N_FLOORS; floor++ {
		for btn := 0; btn < N_FLOORS; btn++ {
			SetButtonLamp(btn,floor,es.Requests[floor][btn])
		}
	}
}

func Fsm_onRequestButtonPress(btn_floor int, btn ButtonType) {

	switch elevator.Behaviour {
		case EB_DoorOpen: 
			if elevator.Floor == btn_floor {
				Timer_start(elevator.DoorOpenDuration_s) 
			} else {
				elevator.Requests[btn_floor][btn] = true
			}
		case EB_Moving:
			elevator.Requests[btn_floor][btn] = true 
		case EB_Idle:
			if elevator.Floor == btn_floor {
				SetDoorOpenLamp(true)
				Timer_start(elevator.DoorOpenDuration_s)
				elevator.Behaviour = EB_DoorOpen
			} else {
				elevator.Requests[btn_floor][btn] = true 
				elevator.Dirn = Requests_chooseDirection(elevator)
				SetMororDirection(elevator.Dirn)
				elevator.Behaviour = EB_Moving
			} 
	}

	setAllLights(elevator)
	fmt.Printf("\nNew state:")
	Elevator_print(elevator)

}


