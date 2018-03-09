package requests

import (
	."SingleElevator/elevator"
	."SingleElevator/elevator_io_types"
	"fmt"
)

const ( // Some constans that should not be decleared here.
	DoorOpenDuration_s = 5
)

func requests_above(e Elevator) bool {
	for f := e.Floor + 1; f < N_FLOORS; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if e.Requests[f][btn] == 1 {
				return true
			}
		}
	} 
	return false
}

func requests_below(e Elevator) bool {
	for f := 0; f < e.Floor; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if e.Requests[f][btn] == 1 {
				return true
			}
		}
	}
	return false
}

func requests_chooseDirection(e Elevator) Dirn {
	switch e.Direction {
	case D_Up :
		if requests_above(e) {
			return D_Up
		} else if requests_below(e) {
			return D_Down
		} else {
			return D_Stop
		}
	case D_Stop | D_Down:
		if requests_below(e) {
			return D_Down
		} else if requests_above(e) {
			return D_Up
		} else {
			return D_Stop
		}
	default:
		return D_Stop
	}
	return D_Stop
}

func requests_shouldStop(e Elevator) bool {
	switch e.Direction {
	case D_Down:
		return e.Requests[e.Floor][B_HallDown] == 1 || e.Requests[e.Floor][B_Cab] == 1 || !requests_below(e)
	case D_Up:
		return e.Requests[e.Floor][B_HallUp] == 1 ||  e.Requests[e.Floor][B_Cab] == 1 || !requests_above(e)
	case D_Stop:
		fallthrough
	default:
		return true
	}
}

func reqests_clearAtCurrentFloor(e Elevator) Elevator {
	e.Requests[e.Floor][B_Cab] = 0
	switch e.Direction {
	case D_Up:
		e.Requests[e.Floor][B_HallUp] = 0
		if !requests_above(e) {
			e.Requests[e.Floor][B_HallDown] = 0
		}
	case D_Down:
		e.Requests[e.Floor][B_HallDown] = 0
		if !requests_below(e) {
			e.Requests[e.Floor][B_HallUp] = 0
		}
	case D_Stop:
		fallthrough
	default:
		e.Requests[e.Floor][B_HallUp] = 0
		e.Requests[e.Floor][B_HallDown] = 0
	}
	return e
}

/* // For testing purposes
 func main(){
	//p() := fmt.Println()
	//var heis Elevator
	heis := Elevator_uninitialized()
	//heis.config.doorOpenDuration_s = DoorOpenDuration_s
	heis.Floor = 1
	heis.Direction = D_Up
	fmt.Println(heis.Floor)
	heis.Requests[1][B_HallUp] = 1
	heis.Requests[1][B_HallDown] = 1
	heis.Requests[1][B_Cab] = 1
	Elevator_print(heis)
	fmt.Println("**********************************")
	Elevator_print(reqests_clearAtCurrentFloor(heis))
	//p(heis.floor)
} 
*/