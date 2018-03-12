package fsm 

import (
	."SingleElevator/elevator"
	."SingleElevator/timer"
	."param"
	."SingleElevator/elevio"
	."SingleElevator/requests"
	"fmt"
)

var elevator Elevator 
//vil at jeg skal lag et ElevatorOutputDevice her....

//vet ikke hva jeg skal gjøre med denne enda, 
//er en __attribute__((constructor)) som bruker elevator.con osv...
func Fsm_init() {
	elevator = Elevator_uninitialized()
	Elevio_init(Panelport,N_FLOORS)
	Elevio_setStopLamp(true)
	Elevio_setDoorOpenLamp(false)
	//unsure if we should set floor indicator before we know what floor we're on
	//SetFloorIndicator(0)
}

func setAllLights(es Elevator) {

	for floor := 0; floor < N_FLOORS; floor++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			Elevio_setButtonLamp(ButtonType(btn),floor,es.Requests[floor][btn])
		}
	}
}

func Fsm_onInitBetweenFloors() {
	Elevio_setMotorDirection(D_Down)
	elevator.Direction = D_Down
	elevator.Behaviour = EB_Moving
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
				elevator.Direction = Requests_chooseDirection(elevator)
				Elevio_setMotorDirection(elevator.Direction)
				elevator.Behaviour = EB_Moving
			} 
	}

	setAllLights(elevator)
	fmt.Println("\nNew state:")
	Elevator_print(elevator)

}

func Fsm_onFloorArrival(newFloor int) {
	fmt.Println("Arrived at floor", newFloor)
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
		// NOP
	}

	fmt.Println("\nNew state:")
	Elevator_print(elevator)
}

func Fsm_onDoorTimeout() {
	fmt.Println("Door closed")
	Elevator_print(elevator)

	switch elevator.Behaviour {
	case EB_DoorOpen:
		elevator.Direction = Requests_chooseDirection(elevator)
		Elevio_setDoorOpenLamp(false)
		Elevio_setMotorDirection(elevator.Direction)

		if elevator.Direction == D_Stop {
			elevator.Behaviour = EB_Idle
		} else {
			elevator.Behaviour = EB_Moving
		}
	default:
		// NOP
	}

	fmt.Println("\nNew state:")
	Elevator_print(elevator)
}