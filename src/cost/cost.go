package cost

import (
    "fmt"
    "os/exec"
    "encoding/json"
    ."types"
)

func Cost(sendAssignedOrders chan<- AssignedOrders, receiveSyncArray <-chan SyncArray, LocalElevatorID string) {
    var newOrderList AssignedOrders 

    for {
        newSyncArray := <- receiveSyncArray 

        convertedSyncArray := syncArrayToAssignerConverter(newSyncArray)

        assignerInput,_ := json.Marshal(convertedSyncArray)

        cmd := exec.Command("sh", "-c", "./hall_request_assigner --input '" + string(assignerInput) + "' --includeCab ")
    
        result ,err := cmd.Output()

        if err != nil {
            fmt.Println(err)
        }

        var formattedResult map[string][N_FLOORS][N_BUTTONS]bool

        json.Unmarshal(result, &formattedResult) 

        // Adding all global hall requests in order to set correct panel lights 
        newOrderList.GlobalHallReq = convertedSyncArray.HallRequests 
        newOrderList.Local = formattedResult[LocalElevatorID]
       
        sendAssignedOrders <- newOrderList
    }
}

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

    for elevId := range inputSyncArray.AllElevators { 
        temp := elevatorToAssignerConverter(inputSyncArray.AllElevators[elevId]) 
        convertedSyncArray.States[elevId] = &temp
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