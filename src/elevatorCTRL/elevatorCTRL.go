package main

import "./elevio"
import "fmt"
import "os/exec"
import "os"
import "path/filepath"

func main(){
    dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
            fmt.Println(err)
    }
    fmt.Println(dir,"\n")
	cmd := exec.Command("sh", "-c", dir+"/hall_request_assigner", "-i", `'{"hallRequests":[[false,false],[true,false],[false,false],[false,true]],"states":{"one":{"behaviour":"moving","floor":2,"direction":"up","cabRequests":[false,false,true,true]},"two":{"behaviour":"idle","floor":0,"direction":"stop","cabRequests":[false,false,false,false]}}}'`)
	fmt.Println(cmd, "\n")	
	result, err:= cmd.Output()

	fmt.Println("\nhall request assigner result:", result, err)



    numFloors := 4

    elevio.Init("localhost:15657", numFloors)

    var CurrFloor int 
    var HallReq int
    
    //var dir elevio.MotorDirection = elevio.MD_Up
    //elevio.SetMotorDirection(d)
    
    drv_buttons := make(chan elevio.ButtonEvent)
    drv_floors  := make(chan int)
    drv_obstr   := make(chan bool)
    drv_stop    := make(chan bool)    
    
    go elevio.PollButtons(drv_buttons)
    go elevio.PollFloorSensor(drv_floors)
    go elevio.PollObstructionSwitch(drv_obstr)
    go elevio.PollStopButton(drv_stop)



    elevio.SetMotorDirection(elevio.MD_Down)

    for {
        select {
        case a := <- drv_floors:
            CurrFloor = a
            if CurrFloor == HallReq {
                elevio.SetMotorDirection(elevio.MD_Stop)
            }
            fmt.Printf("%+v\n", a)
        case a := <- drv_buttons:
            fmt.Printf("%+v\n", a)
            HallReq = a.Floor
            elevio.SetButtonLamp(a.Button, a.Floor, true)
            if a.Floor > CurrFloor {
                elevio.SetMotorDirection(elevio.MD_Up) 
            } else if a.Floor < CurrFloor {
                elevio.SetMotorDirection(elevio.MD_Down)
            }
        }

    }

    
  /*  
    for {
        select {
        case a := <- drv_buttons:
            fmt.Printf("%+v\n", a)
            elevio.SetButtonLamp(a.Button, a.Floor, true)
            
        case a := <- drv_floors:
            fmt.Printf("%+v\n", a)
            if a == numFloors-1 {
                dir = elevio.MD_Down
            } else if a == 0 {
                dir = elevio.MD_Up
            }
            elevio.SetMotorDirection(dir)
            
            
        case a := <- drv_obstr:
            fmt.Printf("%+v\n", a)
            if a {
                elevio.SetMotorDirection(elevio.MD_Stop)
            } else {
                elevio.SetMotorDirection(dir)
            }
            
        case a := <- drv_stop:
            fmt.Printf("%+v\n", a)
            for f := 0; f < numFloors; f++ {
                for b := elevio.ButtonType(0); b < 3; b++ {
                    elevio.SetButtonLamp(b, f, false)
                }
            }
        }
    }   */
}
