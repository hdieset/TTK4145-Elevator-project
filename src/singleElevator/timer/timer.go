package timer

import "time"

var doorTimerEndTime float64 
var doorTimerActive  bool
var stuckTimerEndTime float64
var stuckTimerActive bool

func get_wall_time() float64 {
	now := time.Now()
	var nano float64 = float64(now.UnixNano())/1000000000
	return nano 
} 

func Timer_doorStart(duration_s float64) {
	doorTimerEndTime = get_wall_time() + duration_s
	doorTimerActive = true;
} 

func Timer_stuckStart(duration_s float64) {
	stuckTimerEndTime = get_wall_time() + duration_s
	stuckTimerActive = true;
}

func Timer_doorStop() {
	doorTimerActive = false
}

func Timer_stuckStop() {
	stuckTimerActive = false
}

func Timer_timedOut(doorTimedOut chan<- bool, stuckTimedOut chan<- bool) {
	sleeptime := 20 * time.Millisecond
	for {
		time.Sleep(sleeptime)
		if doorTimerActive && (get_wall_time() > doorTimerEndTime) {
			Timer_doorStop()
			doorTimedOut <- true
		} else if stuckTimerActive && (get_wall_time() > stuckTimerEndTime){
			Timer_stuckStop()
			stuckTimedOut <- true
		}
	}
}