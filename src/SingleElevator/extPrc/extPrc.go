package extPrc 

import (
	"os/exec"
	"os"
	"fmt"
	."types"
)

var process *exec.Cmd

func ExtPrc_initElevatorServer() {
	process = exec.Command("ElevatorServer")
	err := process.Start()
 		if err != nil {
 			panic("Failed to start ElevatorServer.")
 		}
 }

func ExtPrc_exitElevatorServer() {
	err := process.Process.Signal(os.Kill)
		if err != nil {
			panic("Failed to terminate ElevatorServer")
		}
}  

func ExtPrc_changeElevatorSimPort() {
 	fmt.Print("Enter simulator port (enter to use default): ")
 	var input string
 	fmt.Scanln(&input)

 	if input != "" {
 		Panelport = "localhost:" + input 
 	}
}