package main

import (
    "fmt"
    "os/exec"
    "encoding/json"
)



func main(){
    const input string = `{"hallRequests":[[false,false],[true,false],[false,false],[false,true]],"states":{"one":{"behaviour":"moving","floor":2,"direction":"up","cabRequests":[false,false,true,true]},"two":{"behaviour":"idle","floor":0,"direction":"stop","cabRequests":[false,false,false,false]}}}`
    
    //dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

    const dir string = "$GOPATH" + "/src/elevatorCTRL"    
    //result, err := exec.Command("bash", "-c", dir+"/hall_request_assigner --includeCab --input '" + input+ "'").Output
    fmt.Println(dir+"/hall_request_assigner")

    cmd := exec.Command("sh", "-c", dir+"/hall_request_assigner --input '" + input+ "' --includeCab ")

    result ,err := cmd.Output()
    var a map[string][][]bool
    json.Unmarshal(result, &a)
    
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(a)

        

}