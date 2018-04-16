package fsm 

import (
	."singleElevator/elevator"
	."singleElevator/timer"
	."types"
	."singleElevator/elevio"
	."singleElevator/requests"
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
}

func Fsm_ReceivedNewOrderList(newOrders AssignedOrders, syncLocalElevator chan<- Elevator) {
	// Declare flag for special cases
	var idleAtFloor bool = false
	var doorOpenAtFloor bool = false 

	// Clear old elevator.Requests
	var emptyOrderList [N_FLOORS][N_BUTTONS] bool
	elevator.Requests = emptyOrderList

	// Add new orders
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

	prevDir := elevator.Direction
	elevator.Direction = Requests_chooseDirection(elevator)

	if doorOpenAtFloor {
		Timer_doorStart(elevator.DoorOpenDuration_s) 
	} else if idleAtFloor {
		Elevio_setDoorOpenLamp(true)
		Timer_doorStart(elevator.DoorOpenDuration_s)
		elevator.Behaviour = EB_DoorOpen

		// We want to prevent the elevator of driving with an open door, 
		// and also not restart the stuckTimer (when stuck) everytime the function is called 
		// because of updates from other elevators 
	} else if (elevator.Direction != D_Stop) && (elevator.Behaviour != EB_DoorOpen) && (elevator.Direction != prevDir) {
		Elevio_setMotorDirection(elevator.Direction)
		elevator.Behaviour = EB_Moving
		Timer_stuckStart(MAXTRAVELDURATION)
	}

	saveCabOrders(elevator.Requests)
	setAllCabLights()
	setAllHallLights(newOrders)
	sendLocalElevator(syncLocalElevator)
}

func Fsm_onFloorArrival(newFloor int, syncLocalElevator chan<- Elevator) {
	fmt.Println("Arrived at floor", newFloor)

	Timer_stuckStart(MAXTRAVELDURATION)

	elevator.Floor = newFloor
	Elevio_setFloorIndicator(elevator.Floor)

	switch elevator.Behaviour {
	case EB_Moving:
		if Requests_shouldStop(elevator) {
			Timer_stuckStop()
			Elevio_setMotorDirection(D_Stop)
			Elevio_setDoorOpenLamp(true)
			elevator = Requests_clearAtCurrentFloor(elevator)
			Timer_doorStart(elevator.DoorOpenDuration_s)
			saveCabOrders(elevator.Requests)
			setAllCabLights()
			elevator.Behaviour = EB_DoorOpen
		}

	//If the elevator is initialized on a floor, then the stuck timer should stop	
	case EB_Idle: 
		Timer_stuckStop()
	default:
	}

	sendLocalElevator(syncLocalElevator)
}

func Fsm_onDoorTimeout(syncLocalElevator chan<- Elevator) {
	fmt.Println("Door closed")

	switch elevator.Behaviour {
	case EB_DoorOpen:
		elevator.Direction = Requests_chooseDirection(elevator)
		Elevio_setDoorOpenLamp(false)
		Elevio_setMotorDirection(elevator.Direction)

		if elevator.Direction == D_Stop {
			elevator.Behaviour = EB_Idle
		} else {
			elevator.Behaviour = EB_Moving
			Timer_stuckStart(MAXTRAVELDURATION)
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
	// Clear old elevator.CompletedReq
	elevator.CompletedReq = [N_FLOORS][N_BUTTONS]bool{}
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