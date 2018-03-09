package main

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

 func main(){
	//p() := fmt.Println()
	//var heis Elevator
	heis := Elevator_uninitialized()
	//heis.config.doorOpenDuration_s = DoorOpenDuration_s
	heis.Floor = 3
	heis.Direction = D_Down
	fmt.Println(heis.Floor)
	heis.Requests[2][1] = 1
	fmt.Println(requests_below(heis))
	fmt.Println(requests_chooseDirection(heis))
	//p(heis.floor)
} 