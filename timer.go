package main

import (
	"time"
)

var timerEndTime float64
var timerActive bool

func getWallTime() float64 {
	now := time.Now()
	return float64(now.Unix()) + float64(now.Nanosecond())*1e-9
}

func timerStart(duration float64) {
	timerEndTime = getWallTime() + duration
	timerActive = true
}

func timerStop() {
	timerActive = false
}

func timerTimedOut() bool {
	return timerActive && getWallTime() > timerEndTime
}
