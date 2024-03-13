package timer

import (
	"project/costFunc"
	"project/elevator"
	"project/elevio"
	"time"
)

func TimerStart(e_p *elevator.Elevator, wv_p *elevator.Worldview, duration time.Duration, timer chan bool, obstruction chan bool) {
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

func OperativeWatchdog(e_p *elevator.Elevator, wv_p *elevator.Worldview, d time.Duration, wd_chan chan bool) {
	wd_over := time.NewTimer(0)
	defer wd_over.Stop()
	wd_over.Stop()
	var test int = 0
	for {
		select {
		case msg := <-wd_chan:
			test++
			if msg {
				wd_over.Reset(d * time.Second)
			} else {
				wd_over.Stop()
				e_p.Operative = true
			}
		case <-wd_over.C:
			e_p.Operative = false
			wv_p.ElevList[e_p.ElevNum-1].Operative = false
			for floor, f := range wv_p.HallRequests {
				for buttonType, o := range f {
					if o == uint8(e_p.ElevNum) {
						buttn := elevio.ButtonEvent{Floor: floor, Button: elevio.ButtonType(buttonType)}
						costFunc.CostFunction(wv_p, buttn)
					}
				}
			}
			wv_p.Version_up()
		}
	}
}
