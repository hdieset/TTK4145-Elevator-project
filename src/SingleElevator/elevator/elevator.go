package main

import (
	//."SingleElevator/timer" 
	"fmt"
	."SingleElevator/elevator_io_types" 
	."SingleElevator/elevator_io_device"
)


type ElevatorBehaviour int
const(
	EB_Idle		= iota   // Just set arbitrary values for these
	EB_DoorOpen
	EB_Moving
)

type Elevator struct {
	floor 		int
	dirn		Dirn
	requests	[N_FLOORS][N_FLOORS]int
	behaviour 	ElevatorBehaviour
	//this used to be a struct, changed it because we dont have "clearRequestVariant"
	doorOpenDuration_s float64
}


func eb_toString(eb ElevatorBehaviour) string {
	switch eb {
		case EB_Idle:
			return "EB_Idle"
		case EB_DoorOpen:
			return "EB_DoorOpen" 
		case EB_Moving:
			return "EB_Moving"
		default: 
			return "EB_UNDEFINED"
	}
}

//problemer med at button ikke er en enum, må fikse senere...
func elevator_print(es Elevator) {
	p := fmt.Printf
	p("  +--------------------+\n")
	p("  |floor = %-2d          |\n", es.floor)
	p("  |dirn  = %-12.12s|\n", Elevio_dirn_toString(es.dirn))
	p("  |behav = %-12.12s|\n", eb_toString(es.behaviour))
    p("  +--------------------+\n")
    p("  |  | up  | dn  | cab |\n")
    for f := N_FLOORS-1; f >= 0; f-- {
    	p("  | %d", f)
    	for btn := 0; btn < N_BUTTONS; btn++ {
    		if((f == N_FLOORS-1 && btn == B_HallUp)  || (f == 0 && btn == B_HallDown)){
                p("|     ")
            } else if es.requests[f][btn] == 1 {
            	p("|  #  ")
            } else {
            	p("|  -  ")
            }
    	}
    	p("|\n")
    }
    p("  +--------------------+\n")
} 



func elevator_uninitialized() Elevator {
	e := Elevator {
		floor: -1,
		dirn: D_Stop, 
		behaviour: EB_Idle, 
		doorOpenDuration_s: 3.0, 
	}
	return e
}

func main() {
	test := elevator_uninitialized()
	elevator_print(test)
}

