package elevator_io_device 

import ."param"

//there should be more functions above and under theese stringfunctions

func Elevio_dirn_toString(d Dirn) string {
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

func Elevio_button_toString(b ButtonType) string {
	switch b {
		case B_HallUp:
			return "B_HallUp"
		case B_HallDown:
			return "B_HallDown" 
		case B_Cab:
			return "B_Cab"
		default: 
			return "B_UNDEFINED"
	}
}