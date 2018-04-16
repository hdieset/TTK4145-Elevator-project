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

func Fsm_init() {
	elevator = Elevator_uninitialized()
	Elevio_init(Panelport,N_FLOORS)
	Elevio_setStopLamp(true)
	Elevio_setDoorOpenLamp(false)
}

func Fsm_onInitBetweenFloors() {
	Elevio_setMotorDirection(D_Down)
	elevator.Direction = D_Down
	elevator.Behaviour = EB_Moving
	//Timer_movingStart(MAXTRAVELDURATION) // Denne bør vel være her????????????????
}

func Fsm_ReceivedNewOrderList(newOrders AssignedOrders, syncLocalElevator chan<- Elevator) {
	//declare flag for special cases
	var idleAtFloor bool = false
	var doorOpenAtFloor bool = false 

	//clear old elevator.Requests
	var emptyOrderList [N_FLOORS][N_BUTTONS] bool
	elevator.Requests = emptyOrderList

	fmt.Println("Elevator Behaviour:", elevator.Behaviour)

	//add new orders to elevator.Requests 
	for floor := 0; floor < N_FLOORS; floor++ {
		for buttons := 0; buttons < N_BUTTONS; buttons++ {
			if newOrders.Local[floor][buttons] {

				switch elevator.Behaviour {
				case EB_DoorOpen: 
					if elevator.Floor == floor {
						doorOpenAtFloor = true 
						elevator.CompletedReq[floor][buttons] = true 
						//elevator.Requests[floor][buttons] = false -------------- er allerede false 
					} else {
						elevator.Requests[floor][buttons] = true
					}
	
				case EB_Moving:
					elevator.Requests[floor][buttons] = true 

				case EB_Idle:
					fmt.Println("Case EB_Idle")
					if elevator.Floor == floor {
						idleAtFloor = true
						elevator.CompletedReq[floor][buttons] = true
						fmt.Println("Idle at floor = true**************************************")
						//elevator.Requests[floor][buttons] = false
					} else {
						elevator.Requests[floor][buttons] = true
					}	
				}
			}
		}
	}

	if doorOpenAtFloor {
		Timer_doorStart(elevator.DoorOpenDuration_s) 
	} 

	if idleAtFloor {
		fmt.Println("Gikk inn i if setning - idleAtFloor*******************************")
		Elevio_setDoorOpenLamp(true)
		Timer_doorStart(elevator.DoorOpenDuration_s)
		elevator.Behaviour = EB_DoorOpen
	} else if elevator.Behaviour != EB_DoorOpen {
		if elevator.Direction = Requests_chooseDirection(elevator); elevator.Direction != D_Stop {
			Elevio_setMotorDirection(elevator.Direction)
			elevator.Behaviour = EB_Moving
		}
		//Timer_movingStart(MAXTRAVELDURATION) // LAGT TIL HER ??????????????????!!!!!!!!!
	}

	setAllHallLights(newOrders)
	setAllCabLights()
	//fmt.Println("\nNew state:")
	//Elevator_print(elevator)
	sendLocalElevator(syncLocalElevator)
}

func Fsm_onFloorArrival(newFloor int, syncLocalElevator chan<- Elevator) {
	fmt.Println("Arrived at floor", newFloor)
	
	//Timer_movingStart(MAXTRAVELDURATION) 
	//og starte timeren når heisen blir satt til MOVING 

	elevator.Floor = newFloor
	//Elevator_print(elevator)

	Elevio_setFloorIndicator(elevator.Floor)

	switch elevator.Behaviour {
	case EB_Moving:
		if Requests_shouldStop(elevator) {
			//Timer_movingStop()
			Elevio_setMotorDirection(D_Stop)
			Elevio_setDoorOpenLamp(true)
			elevator = Requests_clearAtCurrentFloor(elevator)
			Timer_doorStart(elevator.DoorOpenDuration_s) //endre til consten?? 
			setAllCabLights()
			elevator.Behaviour = EB_DoorOpen
		}
	default:
	}

	sendLocalElevator(syncLocalElevator)
}

func Fsm_onDoorTimeout(syncLocalElevator chan<- Elevator) {
	fmt.Println("Door closed")
	//Elevator_print(elevator)

	switch elevator.Behaviour {
	case EB_DoorOpen:
		elevator.Direction = Requests_chooseDirection(elevator)
		Elevio_setDoorOpenLamp(false)
		Elevio_setMotorDirection(elevator.Direction)

		if elevator.Direction == D_Stop {
			fmt.Println("doorTimed out: Behaviour is set to IDLE")
			elevator.Behaviour = EB_Idle
		} else {
			elevator.Behaviour = EB_Moving
			//Timer_movingStart(MAXTRAVELDURATION)
		}
	default:
	}

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