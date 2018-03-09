package elevator

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
	Floor 		int
	Direction	Dirn
	Requests	[N_FLOORS][N_FLOORS]int
	Behaviour 	ElevatorBehaviour
	DoorOpenDuration_s float64
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
func Elevator_print(es Elevator) {
	p := fmt.Printf
	p("  +--------------------+\n")
	p("  |floor = %-2d          |\n", es.Floor)
	p("  |dirn  = %-12.12s|\n", Elevio_dirn_toString(es.Direction))
	p("  |behav = %-12.12s|\n", eb_toString(es.Behaviour))
    p("  +--------------------+\n")
    p("  |  | dn  | up  | cab |\n")
    for f := N_FLOORS-1; f >= 0; f-- {
    	p("  | %d", f)
    	for btn := 0; btn < N_BUTTONS; btn++ {
    		if((f == N_FLOORS-1 && Button(btn) == B_HallUp)  || (f == 0 && Button(btn) == B_HallDown)){
                p("|     ")
            } else if es.Requests[f][btn] == 1 {
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
		DoorOpenDuration_s: 3.0, 
	}
	return e
}

/*func main() {
	test := Elevator_uninitialized()
	test.Direction = D_Up
	test.Requests[2][1] = 1
	test.Requests[1][2] = 1
	Elevator_print(test)
}*/

