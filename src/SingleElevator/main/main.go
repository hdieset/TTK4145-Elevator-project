package main
import (
		."/SingleElevator/fsm"
		"fmt"
		."param"
		."/SingleElevator/elevio"
)



 func main() {
 	//initializing fmt and driverconnection 
 	Fsm_init() 

 	//setting up driver channels 
 	drv_buttons  := make(chan ButtonEvent)
 	drv_floors 	 := make(chan int)
 	doorTimedOut := make(chan bool)

 	//don't need theese yet 
 	/*drv_obstr 	 := make(chan bool)
 	drv_stop	 := make(chan bool)*/

 	//staring driver polling as go routines with adhering channels 
 	go Elevio_pollButtons(drv_buttons)
 	go Elevio_pollFloorSensor(drv_floors)
 	go Elevio_pollObstruction(drv_obstr)
 	go Elevio_pollStopButton(drv_stop)
 	go Timer_timedOut(doorTimedOut)
 
	fmt.Println("Started!")

	//checking if elevator is between floors and dealing with it 
	if initFloor := <- drv_floors; initFloor == -1 {
		Fsm_onInitBetweenFloors()
	}



	for {
		select {
		case buttonPress := <- drv_buttons: 
			Fsm_onRequestButtonPress(buttonPress.Floor, buttonPress.Button)
		case arrivedAtFloor := <- drv_floors: 
			Fsm_onFloorArrival(arrivedAtFloor)
		case <- doorTimedOut:
			Fsm_onDoorTimeout() 
			Timer_stop() 
		}
	}
}