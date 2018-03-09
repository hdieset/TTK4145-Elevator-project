package main 

import (
	//"SingleElevator/elevator_io_types"
	//"SingleElevator/elevator"
	//"SingleElevator/timer"
	//"SingleElevator/elevator_io.go"
	//"param"
	//"SingleElevator/elevator_io"
	)

var elevator Elevator 
//vil at jeg skal lag et ElevatorOutputDevice her....

//vet ikke hva jeg skal gjøre med denne enda, 
//er en __attribute__((constructor)) som bruker elevator.con osv...
func fsm_init() {
	elevator = Elevator_uninitialized()
	elevio.Init(PANELPORT,N_FLOORS)
}

func setAllLights(es Elevator) {
	for floor = 0
}

