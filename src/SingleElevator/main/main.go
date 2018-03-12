package main 

//test main for SingleElevator 
import (
	."SingleElevator/SingleElevatorMain"
	"fmt"
	."param"
	."SingleElevator/elevator"
	."SingleElevator/extPrc"
	."Cost"
)

func main() {

	if SIMULATOR {
		ExtPrc_changeElevatorSimPort()
	} else {
		ExtPrc_initElevatorServer()
	}

	syncLocalElevator := make(chan Elevator)
	syncButtonPress := make(chan ButtonEvent)
	sendAssignedOrders := make(chan AssignedOrders)
	stopButtonPressed := make(chan bool)

	go SingleElevator(syncLocalElevator, syncButtonPress, sendAssignedOrders, stopButtonPressed)
	go Cost(sendAssignedOrders)

	for {
		select {
		case receivedElev := <- syncLocalElevator: 
			fmt.Println("Received elevator object")
			Elevator_print(receivedElev)
			fmt.Println(receivedElev.CompletedReq)

		case newButtonPress:= <- syncButtonPress:
			fmt.Println("New button push at floor: ", newButtonPress.Floor)
			fmt.Println("And button type (Up =0, Dwn = 1, Cab =2): ", newButtonPress.Button)

		case <- stopButtonPressed:
			return
		}
	}
	if !SIMULATOR {
		ExtPrc_exitElevatorServer()
	}
}