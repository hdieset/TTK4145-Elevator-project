package requests

import ."types"


func requests_above(e Elevator) bool {
	for f := e.Floor + 1; f < N_FLOORS; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if e.Requests[f][btn] {
				return true
			}
		}
	} 
	return false
}

func requests_below(e Elevator) bool {
	for f := 0; f < e.Floor; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if e.Requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func Requests_chooseDirection(e Elevator) Dirn {
	switch e.Direction {
	case D_Up :
		if requests_above(e) {
			return D_Up
		} else if requests_below(e) {
			return D_Down
		} else {
			return D_Stop
		} 
	case D_Down:
		fallthrough
	case D_Stop:
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

func Requests_shouldStop(e Elevator) bool {
	switch e.Direction {
	case D_Down:
		return e.Requests[e.Floor][B_HallDown] || e.Requests[e.Floor][B_Cab] || !requests_below(e)
	case D_Up:
		return e.Requests[e.Floor][B_HallUp] ||  e.Requests[e.Floor][B_Cab] || !requests_above(e)
	case D_Stop:
		fallthrough
	default:
		return true
	}
}

func Requests_clearAtCurrentFloor(e Elevator) Elevator {
	e.CompletedReq[e.Floor][B_Cab] = true 
	e.Requests[e.Floor][B_Cab] = false

	switch e.Direction {
	case D_Up:
		e.Requests[e.Floor][B_HallUp] = false
		e.CompletedReq[e.Floor][B_HallUp] = true 
		if !requests_above(e) {
			e.Requests[e.Floor][B_HallDown] = false
			e.CompletedReq[e.Floor][B_HallDown] = true
		}
	case D_Down:
		e.Requests[e.Floor][B_HallDown] = false
		e.CompletedReq[e.Floor][B_HallDown] = true
		if !requests_below(e) {
			e.Requests[e.Floor][B_HallUp] = false
			e.CompletedReq[e.Floor][B_HallUp] = true
		}
	case D_Stop:
		fallthrough
	default:
		e.Requests[e.Floor][B_HallUp] = false
		e.Requests[e.Floor][B_HallDown] = false
		e.CompletedReq[e.Floor][B_HallUp] = true
		e.CompletedReq[e.Floor][B_HallDown] = true
	}
	return e
}
