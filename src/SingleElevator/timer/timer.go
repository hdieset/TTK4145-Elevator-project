package timer

import (
	"time"
)

var doorTimerEndTime float64 
var doorTimerActive  bool
var movingTimerEndTime float64
var movingTimerActive bool

func get_wall_time() float64 {
	now := time.Now()
	var nano float64 = float64(now.UnixNano())/1000000000
	return nano 
} 

func Timer_doorStart(duration_s float64) {
	doorTimerEndTime = get_wall_time() + duration_s
	doorTimerActive = true;
} 

func Timer_movingStart(duration_s float64) { //Timer_stuck ?? 
	movingTimerEndTime = get_wall_time() + duration_s
	movingTimerActive = true;
}

func Timer_doorStop() {
	doorTimerActive = false
}

func Timer_movingStop() {
	movingTimerActive = false
}

func Timer_timedOut(doorTimedOut chan<- bool, movingTimedOut chan<- bool) {
	sleeptime := 20*time.Millisecond
	for {
		time.Sleep(sleeptime)
		if doorTimerActive && (get_wall_time() > doorTimerEndTime) {
			Timer_doorStop()
			doorTimedOut <- true
		} else if movingTimerActive && (get_wall_time() > movingTimerEndTime){
			Timer_movingStop()
			movingTimedOut <- true
		}
	}
}

