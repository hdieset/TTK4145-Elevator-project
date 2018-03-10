package fsm 

import (
	."SingleElevator/elevator"
	."SingleElevator/timer"
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
	Elevio_init(PANELPORT,N_FLOORS)
	Elevio_setStopLamp(false)
	Elevio_setDoorOpenLamp(false)
	//unsure if we should set floor indicator before we know what floor we're on
	//SetFloorIndicator(0)
}

func setAllLights(es Elevator) {

	for floor := 0; floor < N_FLOORS; floor++ {
		for btn := 0; btn < N_FLOORS; btn++ {
			Elevio_setButtonLamp(btn,floor,es.Requests[floor][btn])
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
				Elevio_setDoorOpenLamp(true)
				Timer_start(elevator.DoorOpenDuration_s)
				elevator.Behaviour = EB_DoorOpen
			} else {
				elevator.Requests[btn_floor][btn] = true 
				elevator.Dirn = Requests_chooseDirection(elevator)
				Elevio_setMororDirection(elevator.Dirn)
				elevator.Behaviour = EB_Moving
			} 
	}

	setAllLights(elevator)
	fmt.Printf("\nNew state:")
	Elevator_print(elevator)

}

func Fsm_onFloorArrival(newFloor int) {
	fmt.Println("Arrived at floor ", newFloor)
	Elevator_print(elevator)

	elevator.Floor = newFloor

	Elevio_setFloorIndicator(elevator.Floor)

	switch elevator.Behaviour {
	case EB_Moving:
		if Requests_shouldStop(elevator) {
			Elevio_setMotorDirection(D_Stop)
			Elevio_setDoorOpenLamp(true)
			elevator = Reqests_clearAtCurrentFloor(elevator)
			Timer_start(elevator.DoorOpenDuration_s)
			setAllLights(elevator)
			elevator.Behaviour = EB_DoorOpen
		}
	default:
		fallthrough
	}

	fmt.Println("New state:")
	Elevator_print(elevator)
}

func Fsm_onDoorTimeout() {
	fmt.Println("fsm_onDoorTimeout")
	Elevator_print(elevator)

	switch elevator.Behaviour {
	case EB_DoorOpen:
		elevator.Direction = Requests_chooseDirection(elevator)
		Elevio_setDoorOpenLamp(false)
		Elevio_SetMotorDirection(elevator.Direction)
		if elevator.Direction == D_Stop {
			elevator.Behaviour = EB_Idle
		} else {
			elevator.Behaviour = EB_Moving
		}
	default:
		fallthrough
	}

	fmt.Println("New state:")
	Elevator_print(elevator)
}