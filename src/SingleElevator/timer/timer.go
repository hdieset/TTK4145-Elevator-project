package timer

import (
	"time"
)


func get_wall_time() float64 {
	now := time.Now()
	var nano float64 = float64(now.UnixNano())/1000000000
	return nano 
} 


var timerEndTime float64 
var timerActive  bool

func Timer_start(duration_s float64) {
	timerEndTime = get_wall_time() + duration_s
	timerActive = true;
}


func Timer_stop() {
	timerActive = false
}


func Timer_timedOut() bool {
	return (timerActive && get_wall_time() > timerEndTime)
}

