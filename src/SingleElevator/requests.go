package main

import "fmt"

const ( // Some constans that should not be decleared here.
	N_FLOORS = 3
	N_BUTTONS = 4
	DoorOpenDuration_s = 5
)

type Dirn int 
const (
	D_Down = Dirn(-1)
	D_Stop = Dirn(0)
	D_Up   = Dirn(1)
)

type ElevatorBehavior int
const(
	EB_Idle		= iota   // Just set arbitrary values for these
	EB_DoorOpen
	EB_Moving
)

type Elevator struct {
	floor 		int
	dirn		Dirn
	requests	[N_FLOORS][N_FLOORS]int
	behaviour 	ElevatorBehavior
	
	config struct {
		doorOpenDuration_s int
	}
}

/*func requests_above(x int) int {
	
}*/

func main(){
	var heis Elevator
	heis.config.doorOpenDuration_s = DoorOpenDuration_s
	heis.floor = 2
	fmt.Println(heis.floor)
	fmt.Println(heis.config.doorOpenDuration_s)
}