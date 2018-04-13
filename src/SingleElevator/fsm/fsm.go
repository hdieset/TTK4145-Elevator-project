package fsm 

import (
	."SingleElevator/elevator"
	."SingleElevator/timer"
	."types"
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

func Fsm_onInitBetweenFloors() {
	Elevio_setMotorDirection(D_Down)
	elevator.Direction = D_Down
	elevator.Behaviour = EB_Moving
}

func Fsm_ReceivedNewOrderList(newOrders AssignedOrders, syncLocalElevator chan<- Elevator) {
	//declare flag for special cases
	var idleAtFloor bool = false
	var doorOpenAtFloor bool = false 

	//clear old elevator.Requests
	var emptyOrderList [N_FLOORS][N_BUTTONS] bool
	elevator.Requests = emptyOrderList

	//add new orders to elevator.Requests 
	for floor := 0; floor < N_FLOORS; floor++ {
		for buttons := 0; buttons < N_BUTTONS; buttons++ {
			if newOrders.Local[floor][buttons] {

				switch elevator.Behaviour {
				case EB_DoorOpen: 
					if elevator.Floor == floor {
						doorOpenAtFloor = true 
						elevator.CompletedReq[floor][buttons] = true 
					} else {
						elevator.Requests[floor][buttons] = true
					}
	
				case EB_Moving:
					elevator.Requests[floor][buttons] = true 

				case EB_Idle:
					if elevator.Floor == floor {
						idleAtFloor = true
						elevator.CompletedReq[floor][buttons] = true
					} else {
						elevator.Requests[floor][buttons] = true
					}	
				}
			}
		}
	}

	if doorOpenAtFloor {
		Timer_start(elevator.DoorOpenDuration_s) 
	} 

	if idleAtFloor {
		Elevio_setDoorOpenLamp(true)
		Timer_start(elevator.DoorOpenDuration_s)
		elevator.Behaviour = EB_DoorOpen
	} else {
		elevator.Direction = Requests_chooseDirection(elevator)
		Elevio_setMotorDirection(elevator.Direction)
		elevator.Behaviour = EB_Moving
	}

	setAllHallLights(newOrders)
	setAllCabLights()
	fmt.Println("\nNew state:")
	Elevator_print(elevator)
	sendLocalElevator(syncLocalElevator)
}




func Fsm_onFloorArrival(newFloor int, syncLocalElevator chan<- Elevator) {
	fmt.Println("Arrived at floor", newFloor)
	Elevator_print(elevator)

	elevator.Floor = newFloor

	Elevio_setFloorIndicator(elevator.Floor)

	switch elevator.Behaviour {
	case EB_Moving:
		if Requests_shouldStop(elevator) {
			Elevio_setMotorDirection(D_Stop)
			Elevio_setDoorOpenLamp(true)
			elevator = Requests_clearAtCurrentFloor(elevator)
			Timer_start(elevator.DoorOpenDuration_s)
			//setAllLights(elevator)
			elevator.Behaviour = EB_DoorOpen
		}
	default:
		// NOP
	}

	/*fmt.Println("\nNew state:")
	Elevator_print(elevator)*/ 
	sendLocalElevator(syncLocalElevator)

}

func Fsm_onDoorTimeout(syncLocalElevator chan<- Elevator) {
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

	/*fmt.Println("\nNew state:")
	Elevator_print(elevator)*/ 
	sendLocalElevator(syncLocalElevator)
}




func setAllHallLights(newOrders AssignedOrders) {
	for floor := 0; floor < N_FLOORS; floor++ {
		for btn := 0; btn < N_BUTTONS-1; btn++ {
			Elevio_setButtonLamp(ButtonType(btn),floor,newOrders.GlobalHallReq[floor][btn])
		}
	}
}

func setAllCabLights() {
	for floor := 0; floor < N_FLOORS;floor++ {
		Elevio_setButtonLamp(B_Cab,floor,elevator.Requests[floor][B_Cab])
	}
}

func sendLocalElevator(syncLocalElevator chan<- Elevator) {
	syncLocalElevator <- elevator 
	//clear old elevator.CompletedReq
	var emptyOrderList [N_FLOORS][N_BUTTONS] bool
	elevator.CompletedReq = emptyOrderList
}