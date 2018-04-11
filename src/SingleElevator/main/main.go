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

	syncLocalElevator 	:= make(chan Elevator)
	syncButtonPress		:= make(chan ButtonEvent)
	sendAssignedOrders 	:= make(chan AssignedOrders)
	stopButtonPressed 	:= make(chan bool)
	sendSyncArray		:= make(chan SyncArray)

	//ID FØR DETTE, placeholder: 
	var LocalElevatorID string 
	LocalElevatorID = "penis"

	var localSyncArray SyncArray
	//localSyncArray := new(SyncArray)  - gammel versjon førte til at vi måtte sende peker 
	localSyncArray.AllElevators = make(map[string]Elevator)	

	// Creating a test localSyncArray
	localSyncArray.HallStates[1][B_HallDown] = Hall_confirmed
	localSyncArray.HallStates[3][B_HallDown] = Hall_confirmed
	var penisheis Elevator
	penisheis.Behaviour = EB_Moving
	penisheis.Floor = 2
	penisheis.Direction = D_Up
	penisheis.Requests[2][B_Cab] = true
	penisheis.Requests[3][B_Cab] = true
	localSyncArray.AllElevators["penis"] = penisheis
	var fitteheis Elevator
	fitteheis.Behaviour = EB_Idle
	fitteheis.Floor = 0
	fitteheis.Direction = D_Stop
	//fitteheis.Requests[3][B_Cab] = true
	localSyncArray.AllElevators["fitte"] = fitteheis


	go SingleElevator(syncLocalElevator, syncButtonPress, sendAssignedOrders, stopButtonPressed)
	go Cost(sendAssignedOrders, sendSyncArray, LocalElevatorID)

	sendSyncArray <- localSyncArray //her måtte vi sende en peker. 

	for {
		select {
		case receivedElev := <- syncLocalElevator: // Skal til SyncModule
			fmt.Println("Received elevator object")
			Elevator_print(receivedElev)
			fmt.Println(receivedElev.CompletedReq)

		case newButtonPress:= <- syncButtonPress: // Skal til SyncModule
			fmt.Println("New button push at floor: ", newButtonPress.Floor)
			fmt.Println("And button type (Up =0, Dwn = 1, Cab =2): ", newButtonPress.Button)

		case <- stopButtonPressed: // skal til main
			return
		}
	}
	if !SIMULATOR {
		ExtPrc_exitElevatorServer()
	}
}

