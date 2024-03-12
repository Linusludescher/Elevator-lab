package timer

import (
	"fmt"
	"time"
)

var timerEndTime float64
var timerActive bool

func GetWallTime() float64 {
	now := time.Now()
	return float64(now.Unix()) + float64(now.Nanosecond())*1e-9
}

func TimerStart(duration time.Duration, timer chan bool) {
	fmt.Println("timer started")
	sec_timer := time.NewTimer(duration * time.Second)
	<-sec_timer.C
	timer <- true
}

func TimerStop() {
	timerActive = false
}

func TimerTimedOut() bool {
	return timerActive && GetWallTime() > timerEndTime
}
