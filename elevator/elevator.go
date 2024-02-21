package elevator

import (
	"project/elevio"
)

const N_FLOORS int = 4
const N_BUTTONS int = 3

type Elevator_behaviour int
const (
    EB_Idle Elevator_behaviour = iota + 1
    EB_DoorOpen 
    EB_Moving
)

type Elevator struct {
    Floor           int
    Dirn            elevio.MotorDirection
    Behaviour       Elevator_behaviour
    Requests        [N_FLOORS][N_BUTTONS]int
}

func Elevator_uninitialized() (elevator Elevator) {
    var matrix [N_FLOORS][N_BUTTONS]int
    elevator = Elevator{-1, elevio.MD_Stop, EB_Idle, matrix}
    return
}
