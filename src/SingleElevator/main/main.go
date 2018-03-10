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
 	drv_buttons := make(chan ButtonEvent)
 	drv_floors 	:= make(chan int)
 	drv_obstr 	:= make(chan bool)
 	drv_stop	:= make(chan bool)

 	//staring driver polling as go routines with adhering channels 
 	go Elevio_pollButtons(drv_buttons)
 	go Elevio_pollFloorSensor(drv_floors)
 	go Elevio_pollObstruction(drv_obstr)
 	go Elevio_pollStopButton(drv_stop)
 
	fmt.Println("Started!")

	//checking if elevator is between floors and dealing with it 
	if initFloor := <- drv_floors; initFloor == -1 {
		Fsm_onInitBetweenFloors()
	}



	for {
		select {
		case buttonPress := <- drv_buttons: 
			Fsm_onRequestButtonPress(buttonPress.Floor, buttonPress.Button)
		
		case arrive


		}

	}


	


}