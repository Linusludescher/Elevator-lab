package timer

import (
	"os"
	"project/costFunc"
	"project/elevator"
	"project/elevio"
	"time"
)

func TimerStart(elev_p *elevator.Elevator, worldView_p *elevator.Worldview, duration time.Duration, timer_exp_chan chan<- bool, obstruction <-chan bool, reset_ch <-chan bool) {
	obstructed := elev_p.Obstruction
	sec_timer := time.NewTimer(duration * time.Second)
	defer sec_timer.Stop()
	sec_timer.Stop()
	for {
		select {
		case <-sec_timer.C:
			if !obstructed {
				timer_exp_chan <- true
			} else {
				// Restart the timer if obstructed is true
				sec_timer.Reset(duration * time.Second)
			}
		case obstr := <-obstruction:
			obstructed = obstr
			elev_p.Obstruction = obstr
			worldView_p.Version_up()
			if obstructed {
				// If obstructed becomes true, restart the timer
				if !sec_timer.Stop() {
					<-sec_timer.C // Drain the timer channel if the timer has already expired
				}
				sec_timer.Reset(duration * time.Second)
			}
		case <-reset_ch:
			sec_timer.Reset(duration * time.Second)
		}
	}
}

func OperativeWatchdog(elev_p *elevator.Elevator, worldView_p *elevator.Worldview, d time.Duration, wd_chan <-chan bool) {
	wd_over := time.NewTimer(0)
	defer wd_over.Stop()
	wd_over.Stop()
	for {
		select {
		case msg := <-wd_chan:
			if msg {
				wd_over.Reset(d * time.Second)
			} else {
				wd_over.Stop()
				elev_p.Operative = true
			}
		case <-wd_over.C:
			elev_p.Operative = false
			worldView_p.ElevList[elev_p.ElevNum-1].Operative = false
			for floor, f := range worldView_p.HallRequests {
				for buttonType, o := range f {
					if o == uint8(elev_p.ElevNum) {
						buttn := elevio.ButtonEvent{Floor: floor, Button: elevio.ButtonType(buttonType)}
						costFunc.CostFunction(worldView_p, buttn)
					}
				}
			}
			worldView_p.Version_up()
			os.Exit(1)
		}
	}
}
