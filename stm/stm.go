package stm

import (
	"project/elevio"
	"project/elevator"
	"fmt"
	"project/requests"
)

func Stm(e elevator.Elevator) {
	state_moving := make(chan int)
	state_idle := make(chan int)
	state_doorOpen := make(chan int)
	for {
		select {
		case <-state_moving:
			if e.Dirn == elevio.MD_Up{
				if requests.RequestsAbove(e) {
					
				} else if requestsHere(e) {
					return EB_DoorOpen // Define EB_DoorOpen accordingly
				} else if requestsBelow(e) {
					return EB_Moving // Define EB_Moving accordingly
				}
			}

		case <-state_idle:

		case <-state_doorOpen:
		}
	}
}

// 	for {
// 		select {
// 		case a := <-drv_buttons:
// 			fmt.Printf("%+v\n", a)
// 			elevio.SetButtonLamp(a.Button, a.Floor, true)

// 		case a := <-drv_floors:
// 			fmt.Printf("%+v\n", a)
// 			if a == numFloors-1 {
// 				d = elevio.MD_Down
// 			} else if a == 0 {
// 				d = elevio.MD_Up
// 			}
// 			elevio.SetMotorDirection(d)

// 		case a := <-drv_obstr:
// 			fmt.Printf("%+v\n", a)
// 			if a {
// 				elevio.SetMotorDirection(elevio.MD_Stop)
// 			} else {
// 				elevio.SetMotorDirection(d)
// 			}

// 		case a := <-drv_stop:
// 			fmt.Printf("%+v\n", a)
// 			for f := 0; f < numFloors; f++ {
// 				for b := elevio.ButtonType(0); b < 3; b++ {
// 					elevio.SetButtonLamp(b, f, false)
// 				}
// 			}
// 		}
// 	}
// }
