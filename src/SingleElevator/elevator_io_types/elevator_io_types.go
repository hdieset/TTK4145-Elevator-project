package elevator_io_types

const ( 
	N_FLOORS = 4
	N_BUTTONS = 2
)

type Dirn int 
const (
	D_Down = Dirn(-1)
	D_Stop = Dirn(0)
	D_Up   = Dirn(1)
)

type Button int 
const (
	B_HallUp = Button(-1)
	B_HallDown = Button(0)
	B_Cab = Button(1)
)