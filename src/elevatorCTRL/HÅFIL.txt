Init(addr string, numFloors int) 

SetMotorDirection(dir MotorDirection)

SetButtonLamp(button ButtonType, floor int, value bool)

SetFloorIndicator(floor int)

SetDoorOpenLamp(value bool)

SetStopLamp(value bool)

PollButtons(receiver chan<- ButtonEvent)

PollFloorSensor(receiver chan<- int)

PollStopButton(receiver chan<- bool)

PollObstructionSwitch(receiver chan<- bool)

getButton(button ButtonType, floor int) returns bool

getFloor()returns int 

getStop() returns bool

getObstruction() returns bool
