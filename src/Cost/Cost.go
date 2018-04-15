package cost

import (
    "fmt"
    "os/exec"
    "encoding/json"
    ."types"
)

func Cost(sendAssignedOrders chan<- AssignedOrders, receiveSyncArray <-chan SyncArray, LocalElevatorID string){
    const dir string = "$GOPATH" + "/src/Cost" //dette vil ikke funke med executable(?), må endre gopath
    var newOrderList AssignedOrders 

    for {
        //waiting for sync module to send new sync array
        newSyncArray := <- receiveSyncArray 

        convertedSyncArray := syncArrayToAssignerConverter(newSyncArray)

        assignerInput,_ := json.Marshal(convertedSyncArray)

        cmd := exec.Command("sh", "-c", dir+"/hall_request_assigner --input '" + string(assignerInput) + "' --includeCab ")

        result ,err := cmd.Output()

        if err != nil {
            fmt.Println(err)
        }

        var formattedResult map[string][N_FLOORS][N_BUTTONS]bool

        json.Unmarshal(result, &formattedResult) 

        //Assigning orders to local elevator and global hall requests for lights on panel
        newOrderList.Local = formattedResult[LocalElevatorID]
        newOrderList.GlobalHallReq = convertedSyncArray.HallRequests 

        //Sending new orders to local elevator
        sendAssignedOrders <- newOrderList
    }
}

// Her var det mye fucky pointer drit
func elevatorToAssignerConverter (inputElevator Elevator) AssignerCompatibleElev {

    var convertedElev AssignerCompatibleElev

    switch inputElevator.Behaviour {
    case EB_Moving:
        convertedElev.Behaviour = "moving"
    case EB_Idle: 
        convertedElev.Behaviour = "idle"
    case EB_DoorOpen:
        convertedElev.Behaviour = "doorOpen"    
    }

    if convertedElev.Floor < 0 {
        convertedElev.Floor = 0
    } else {
        convertedElev.Floor = inputElevator.Floor 
    }
    
    switch inputElevator.Direction {
    case D_Up:
        convertedElev.Direction = "up" 
    case D_Stop:
        convertedElev.Direction = "stop"
    case D_Down:
        convertedElev.Direction = "down"
    }

    for floors := 0; floors < N_FLOORS; floors++ {
        convertedElev.CabRequests[floors] = inputElevator.Requests[floors][B_Cab]
    }

    return convertedElev
}

func syncArrayToAssignerConverter (inputSyncArray SyncArray) AssignerCompatibleInput {
    convertedSyncArray := AssignerCompatibleInput{}
    convertedSyncArray.States = make(map[string]*AssignerCompatibleElev)

    for elevIter := range inputSyncArray.AllElevators { 
        temp := elevatorToAssignerConverter(inputSyncArray.AllElevators[elevIter]) 
        convertedSyncArray.States[inputSyncArray.AllElevators[elevIter].Id] = &temp
    } 

    for floors := 0; floors < N_FLOORS; floors++ {
    	for btn := 0; btn < N_BUTTONS-1; btn++ {
	        if inputSyncArray.HallStates[floors][btn] == Hall_confirmed {
	            convertedSyncArray.HallRequests[floors][btn] = true
	        }
      	}
    }

    return convertedSyncArray 
}