package cost

import (
    "fmt"
    "os/exec"
    "encoding/json"
    ."types"
    "time"
)

func Cost(sendAssignedOrders chan<- AssignedOrders, receiveSyncArray <-chan SyncArray, LocalElevatorID string){
    const dir string = "$GOPATH" + "/src/Cost" //dette vil ikke funke med executable, må endre gopath
    var newOrderList AssignedOrders 

    for {
        //waiting for sync module to send new sync array
        newSyncArray := <- receiveSyncArray 

        t := time.Now()
		fmt.Println(t.String()," - Sync array for konvertering:")
        fmt.Println(newSyncArray)
        fmt.Println("****************************************************************")
        convertedSyncArray := syncArrayToAssignerConverter(newSyncArray)
        fmt.Println("Sync Array som sendes til Marshal:")
        fmt.Println(convertedSyncArray)

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

        fmt.Println("newOrderList som sendes til SingleElevator")
        fmt.Println(newOrderList)

        //newOrderList.Local[0][B_Cab] = true 


        //Sending new orders to local elevator
        sendAssignedOrders <- newOrderList
    }
}

// Her var det mye fucky pointer drit
func elevatorToAssignerConverter (inputElevator *Elevator) AssignerCompatibleElev {

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
    //convertedSyncArray.States[LocalElevatorID]

    for elevIdIter := range inputSyncArray.AllElevators {
        /*var cheesyTemp Elevator 
        cheesyTemp.Behaviour = inputSyncArray.AccessAllElevators(elevIdIter).Behaviour 
        cheesyTemp.Floor = inputSyncArray.AccessAllElevators(elevIdIter).Floor 
        cheesyTemp.Direction = inputSyncArray.AccessAllElevators(elevIdIter).Direction 
        cheesyTemp.Requests = inputSyncArray.AccessAllElevators(elevIdIter).Requests 
        temp := elevatorToAssignerConverter(cheesyTemp)
        convertedSyncArray.AccessStates(elevIdIter).Behaviour = temp.Behaviour
        convertedSyncArray.AccessStates(elevIdIter).Floor = temp.Floor
        convertedSyncArray.AccessStates(elevIdIter).Direction = temp.Direction  
        convertedSyncArray.AccessStates(elevIdIter).CabRequests = temp.CabRequests    
        //convertedSyncArray.AccessStates(elevIdIter).Floor = -1 
        //convertedSyncArray.States[elevIdIter] = elevatorToAssignerConverter(inputSyncArray.AllElevators[elevIdIter])*/

        temp := elevatorToAssignerConverter(inputSyncArray.AllElevators[elevIdIter])
        convertedSyncArray.States[elevIdIter] = &temp

    } //over: bruke accessStates??? 

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