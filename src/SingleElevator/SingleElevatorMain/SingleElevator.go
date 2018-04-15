package SingleElevator
import (
	."SingleElevator/fsm"
	"fmt"
	."types"
	."SingleElevator/elevio"
	."SingleElevator/timer"
	//."SingleElevator/extPrc"
	//"time"
)


func SingleElevator(syncLocalElevator chan<- Elevator, 
	syncButtonPress chan<- ButtonEvent,
	receiveAssignedOrders <-chan AssignedOrders,
	stopButtonPressed chan<- bool,
	network_peerTxEnable chan<- bool) {

	elevatorStuck := false 

	/* if SIMULATOR {
		ExtPrc_changeElevatorSimPort()
	} else {
		ExtPrc_initElevatorServer()
	} */

 	//initializing fmt, driverconnection and panellights
 	Fsm_init() 

 	const buffers int = 100
 	//setting up driver channels 
 	drv_buttons  	:= make(chan ButtonEvent, buffers)
 	drv_floors 	 	:= make(chan int, buffers)
 	doorTimedOut 	:= make(chan bool, buffers)
 	movingTimedOut  := make(chan bool, buffers)
 	drv_stop	 	:= make(chan bool, buffers) 

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
			fmt.Println("Knapp for faen")
			if !elevatorStuck {
				syncButtonPress <- buttonPress
			}
		case newOrderlist := <- receiveAssignedOrders:
			Fsm_ReceivedNewOrderList(newOrderlist, syncLocalElevator)
		case arrivedAtFloor := <- drv_floors:
			fmt.Println("Single Elevator says: arrivedAtFloor") 
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
