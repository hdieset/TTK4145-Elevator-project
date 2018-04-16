# TTK4145 Real-time programming - Elevator project
```none
 _____  _                      _                                        _              _   
|  ___|| |                    | |                                      (_)            | |  
| |__  | |  ___ __   __  __ _ | |_   ___   _ __   _ __   _ __   ___     _   ___   ___ | |_ 
|  __| | | / _ \\ \ / / / _` || __| / _ \ | '__| | '_ \ | '__| / _ \   | | / _ \ / __|| __|
| |___ | ||  __/ \ V / | (_| || |_ | (_) || |    | |_) || |   | (_) |  | ||  __/| (__ | |_ 
\____/ |_| \___|  \_/   \__,_| \__| \___/ |_|    | .__/ |_|    \___/   | | \___| \___| \__|
                                                 | |                  _/ |                 
                                                 |_|                 |__/                  
```
## Henrik Nyholm og Herman K. Dieset
![NTNU](https://innsida.ntnu.no/c/wiki/get_page_attachment?p_l_id=22780&nodeId=24647&title=Bruksregler+for+NTNU-logoen&fileName=variant1.jpg)

A truly uplifting experience. An executable can be found in the folder `executable` where the necessary executable `hall_request_assigner` lies. 

The code is separated into four major blocks; a cost module, networking module, singleElevator module and a syncModule.
- The singleElevator controls the movements of the elevator, sets its lights, reads buttons and so on. It sends information to the syncModule regarding its state.
- The cost module gets information about all elevators on the network and assigns a workload to the connected elevator. It also transmitts information about Hall Orders to set the correct panel lights.
- The syncModule is a core part which contains a state machine that ensures redundancy of all orders. It shares state information to other elevators through the networking module. A list of orders which can be assigned to all elevators is sent to the cost module for filtering. 
- The network sends and recieves orders to other elevator nodes on an udp network.
These are depicted below.
![Alt text](elevatorModules.png?raw=true "Title")
