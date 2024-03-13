package timer

import (
	"fmt"
	"project/elevator"
	"time"
)

func TimerStart(e_p *elevator.Elevator, wv_p *elevator.Worldview, duration time.Duration, timer chan bool, obstruction chan bool) {
	fmt.Println("timer started")
	obstructed := e_p.Obstruction
	sec_timer := time.NewTimer(duration * time.Second)
	defer sec_timer.Stop()
	for {
		select {
		case <-sec_timer.C:
			if !obstructed {
				timer <- true
				return
			} else {
				// Restart the timer if obstructed is true
				sec_timer.Reset(duration * time.Second)
			}
		case obstr := <-obstruction:
			obstructed = obstr
			e_p.Obstruction = obstr
			wv_p.Version_up()
			if obstructed {
				// If obstructed becomes true, restart the timer
				if !sec_timer.Stop() {
					<-sec_timer.C // Drain the timer channel if the timer has already expired
				}
				sec_timer.Reset(duration * time.Second)
			}
		}
	}
}
