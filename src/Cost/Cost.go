package cost

import (
    "fmt"
    "os/exec"
    "encoding/json"
    ."types"
)

type assignerCompatibleElev struct {
    Behaviour    string         `json:"behaviour"`
    Floor        int            `json:"floor"`
    Direction    string         `json:"direction"`
    CabRequests  [N_FLOORS]bool `json:"cabRequests"`
}

type assignerCompatibleInput struct {
    HallRequests [N_FLOORS][N_BUTTONS-1]bool       `json:"hallRequests"`
    States       map[string]assignerCompatibleElev `json:"states"`
}


func Cost(sendAssignedOrders chan<- AssignedOrders, receiveSyncArray <-chan SyncArray, LocalElevatorID string){
    const dir string = "$GOPATH" + "/src/Cost"
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

        fmt.Println("********************")
        fmt.Println(formattedResult)
        fmt.Println("********************")

        //Sending new orders to local elevator
        sendAssignedOrders <- newOrderList
    }
}

func elevatorToAssignerConverter (inputElevator Elevator) assignerCompatibleElev {

    var convertedElev assignerCompatibleElev

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

func syncArrayToAssignerConverter (inputSyncArray SyncArray) assignerCompatibleInput {
    convertedSyncArray := assignerCompatibleInput{}
    convertedSyncArray.States = make(map[string]assignerCompatibleElev)

    for elevIdIter := range inputSyncArray.AllElevators {
        convertedSyncArray.States[elevIdIter] = elevatorToAssignerConverter(inputSyncArray.AllElevators[elevIdIter])
    }

    for floors := 0; floors < N_FLOORS; floors++ {
        if inputSyncArray.HallStates[floors][B_HallUp] == Hall_confirmed {
            convertedSyncArray.HallRequests[floors][B_HallUp] = true
        }
        
        if inputSyncArray.HallStates[floors][B_HallDown] == Hall_confirmed {
            convertedSyncArray.HallRequests[floors][B_HallDown] = true
        }
    }

    return convertedSyncArray //??????????????????????????????????????
}
