package elevator

import (
	//."SingleElevator/timer" 
	"fmt"
	."types" 
)

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

func elevio_dirn_toString(d Dirn) string {
	switch d {
		case D_Up:
			return "D_Up"
		case D_Down:
			return "D_Down" 
		case D_Stop:
			return "D_Stop"
		default: 
			return "D_UNDEFINED"
	}
}

//problemer med at button ikke er en enum, må fikse senere...
func Elevator_print(es Elevator) {
	p := fmt.Printf
	p("  +--------------------+\n")
	p("  |floor = %-2d          |\n", es.Floor)
	p("  |dirn  = %-12.12s|\n", elevio_dirn_toString(es.Direction))
	p("  |behav = %-12.12s|\n", eb_toString(es.Behaviour))
    p("  +--------------------+\n")
    p("  |  | up  | dn  | cab |\n")
    for f := N_FLOORS-1; f >= 0; f-- {
    	p("  | %d", f)
    	for btn := 0; btn < N_BUTTONS; btn++ {
    		if((f == N_FLOORS-1 && ButtonType(btn) == B_HallUp)  || (f == 0 && ButtonType(btn) == B_HallDown)){
                p("|     ")
            } else if es.Requests[f][btn] {
            	p("|  #  ")
            } else {
            	p("|  -  ")
            }
    	}
    	p("|\n")
    }
    p("  +--------------------+\n")
} 

func Elevator_uninitialized() Elevator {
	e := Elevator {
		Floor: -1,
		Direction: D_Stop, 
		Behaviour: EB_Idle, 
		DoorOpenDuration_s: DOOROPENDURATION, 
	}
	return e
}
