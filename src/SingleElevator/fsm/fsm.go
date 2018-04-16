package fsm 

import (
	."SingleElevator/elevator"
	."SingleElevator/timer"
	."types"
	."SingleElevator/elevio"
	."SingleElevator/requests"
	"fmt"
	"os"
	"io/ioutil"
)

var elevator Elevator 

func Fsm_init() {
	elevator = Elevator_uninitialized()
	elevator.Requests = readBackupCabOrders()
	Elevio_init(Panelport,N_FLOORS)
	Elevio_setStopLamp(true)
	Elevio_setDoorOpenLamp(false)
	setAllCabLights()
	setAllHallLights(AssignedOrders{})
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
					if elevator.Floor == floor {
						idleAtFloor = true
						elevator.CompletedReq[floor][buttons] = true
						//elevator.Requests[floor][buttons] = false
					} else {
						elevator.Requests[floor][buttons] = true
					}	
				}
			}
		}
	}

	prevDir := elevator.Direction
	if doorOpenAtFloor {
		Timer_doorStart(elevator.DoorOpenDuration_s) 
	} else if idleAtFloor {
		Elevio_setDoorOpenLamp(true)
		Timer_doorStart(elevator.DoorOpenDuration_s)
		elevator.Behaviour = EB_DoorOpen
	} else if elevator.Direction = Requests_chooseDirection(elevator); (elevator.Direction != D_Stop) && (elevator.Direction != prevDir) && (elevator.Behaviour != EB_DoorOpen) {
		Elevio_setMotorDirection(elevator.Direction)
		elevator.Behaviour = EB_Moving
		Timer_movingStart(MAXTRAVELDURATION)
	}

	saveCabOrders(elevator.Requests)
	setAllCabLights()
	setAllHallLights(newOrders)
	//fmt.Println("\nNew state:")
	//Elevator_print(elevator)
	sendLocalElevator(syncLocalElevator)
}

func Fsm_onFloorArrival(newFloor int, syncLocalElevator chan<- Elevator) {
	fmt.Println("Arrived at floor", newFloor)
	
	Timer_movingStart(MAXTRAVELDURATION) 

	elevator.Floor = newFloor

	Elevio_setFloorIndicator(elevator.Floor)

	switch elevator.Behaviour {
	case EB_Moving:
		if Requests_shouldStop(elevator) {
			Timer_movingStop()
			Elevio_setMotorDirection(D_Stop)
			Elevio_setDoorOpenLamp(true)
			elevator = Requests_clearAtCurrentFloor(elevator)
			Timer_doorStart(elevator.DoorOpenDuration_s)
			saveCabOrders(elevator.Requests)
			setAllCabLights()
			elevator.Behaviour = EB_DoorOpen
		}
	//If the elevator is initialized on a floor, then stuck timer stops 	
	case EB_Idle: 
		Timer_movingStop()
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
			elevator.Behaviour = EB_Idle
		} else {
			elevator.Behaviour = EB_Moving
			Timer_movingStart(MAXTRAVELDURATION)
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

func readBackupCabOrders() [N_FLOORS][N_BUTTONS]bool {
    var output [N_FLOORS][N_BUTTONS]bool

    b, err := ioutil.ReadFile("cabBackup.txt")
    if err != nil {
        saveCabOrders([N_FLOORS][N_BUTTONS]bool{})
    } else {
    	for i := 0; i < N_FLOORS; i++ {
	        if b[i] == 49 {
	            output[i][B_Cab] = true
	        }
    	}
	}
    return output
}

func saveCabOrders(input [N_FLOORS][N_BUTTONS]bool) {
    file, err := os.Create("cabBackup.txt")
    if err != nil {
        fmt.Println("Cannot create file", err)
    } else {
        defer file.Close()

        for i := 0; i < N_FLOORS; i++ {
            if input[i][B_Cab] == true {
                fmt.Fprintf(file, "1")
            } else {
                fmt.Fprintf(file, "0")
            }
        }
    }
}