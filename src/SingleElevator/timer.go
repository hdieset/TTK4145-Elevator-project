package main

import (
	"time"
	"fmt"
)


func get_wall_time() float64 {
	now := time.Now()
	var nano float64 = float64(now.UnixNano())/1000000000
	return nano 
} 
