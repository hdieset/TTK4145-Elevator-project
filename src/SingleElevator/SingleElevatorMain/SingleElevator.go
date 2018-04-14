package SingleElevator
import (
	."SingleElevator/fsm"
	"fmt"
	."types"
	."SingleElevator/elevio"
	."SingleElevator/timer"
	//."SingleElevator/extPrc"
)


func SingleElevator(syncLocalElevator chan<- Elevator, 
	syncButtonPress chan<- ButtonEvent,
	receiveAssignedOrders <-chan AssignedOrders,
	stopButtonPressed chan<- bool, network_peerTxEnable chan<- bool ) {

	elevatorStuck := false 

	/* if SIMULATOR {
		ExtPrc_changeElevatorSimPort()
	} else {
		ExtPrc_initElevatorServer()
	} */

 	//initializing fmt, driverconnection and panellights
 	Fsm_init() 

 	//setting up driver channels 
 	drv_buttons  	:= make(chan ButtonEvent)
 	drv_floors 	 	:= make(chan int)
 	doorTimedOut 	:= make(chan bool)
 	movingTimedOut  := make(chan bool)
 	drv_stop	 	:= make(chan bool) 

 	//don't need this one yet 
 	//drv_obstr 	 := make(chan bool)

 	//staring driver polling as go routines with adhering channels 
 	go Elevio_pollButtons(drv_buttons)
 	go Elevio_pollFloorSensor(drv_floors)
	go Timer_timedOut(doorTimedOut, movingTimedOut)
	go Elevio_pollStopButton(drv_stop)

 	//go Elevio_pollObstruction(drv_obstr)
 	
	fmt.Println("Started!")

	if Elevio_getInitialFloor() == -1 {
		Fsm_onInitBetweenFloors()
	}

	for {
		select {
		case buttonPress := <- drv_buttons: 
			if !elevatorStuck {
				syncButtonPress <- buttonPress
			}
		case newOrderlist := <- receiveAssignedOrders:
			Fsm_ReceivedNewOrderList(newOrderlist, syncLocalElevator)
		case arrivedAtFloor := <- drv_floors: 
			if elevatorStuck {
				elevatorStuck = false 
				network_peerTxEnable <- true
				fmt.Println("Dobby is a free el(f)evator!")
			}
			Fsm_onFloorArrival(arrivedAtFloor, syncLocalElevator)
		case <- doorTimedOut:
			Fsm_onDoorTimeout(syncLocalElevator) 
		case <-movingTimedOut: 
			network_peerTxEnable <- false 
			elevatorStuck = true 
			fmt.Println("Aaaaaand we're stuck...")
		case <- drv_stop:
			Elevio_setStopLamp(false)
			fmt.Println("Elevator died peacefully")
			stopButtonPressed <- true
		}
	}

	/* if !SIMULATOR {
		ExtPrc_exitElevatorServer()
	}
	stopButtonPressed <- true */ 
}
