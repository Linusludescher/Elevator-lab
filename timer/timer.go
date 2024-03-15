package timer

import (
	"os"
	"time"
)

func TimerStart(duration time.Duration, timer_exp_chan chan<- bool, reset_timer_chan <-chan bool) {
	sec_timer := time.NewTimer(duration * time.Second)
	defer sec_timer.Stop()
	sec_timer.Stop()
	for {
		select {
		case <-sec_timer.C:
			timer_exp_chan <- true

		case <-reset_timer_chan:
			if !sec_timer.Stop() {
				<-sec_timer.C // Drain the timer channel if the timer has already expired
			}
			sec_timer.Reset(duration * time.Second)
		}
	}
}

func OperativeWatchdog(d time.Duration, watchdog_chan <-chan bool) {
	watchdog_over := time.NewTimer(0)
	defer watchdog_over.Stop()
	watchdog_over.Stop()
	for {
		select {
		case msg := <-watchdog_chan:
			if msg {
				watchdog_over.Reset(d * time.Second)
			} else {
				watchdog_over.Stop()
			}
		case <-watchdog_over.C:
			os.Exit(1)
		}
	}
}
