package elevator

import (
	"fmt"
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
	Dirn       elevio.MotorDirection
    Last_dir   elevio.MotorDirection
	Last_Floor int
	Requests   [N_FLOORS][N_BUTTONS]int
}

func Elevator_uninitialized() (elevator Elevator) {
	var matrix [N_FLOORS][N_BUTTONS]int
	elevator = Elevator{elevio.MD_Stop, elevio.MD_Stop, 0, matrix}
	for elevio.GetFloor() != 0 {
		elevio.SetMotorDirection(elevio.MD_Down)
	}
	elevio.SetMotorDirection(elevio.MD_Stop)
	return
}

func (elevator *Elevator) Display() {
	fmt.Printf("Direction: %v\n", elevator.Dirn)
    fmt.Printf("Last Direction: %v\n", elevator.Last_dir)
	fmt.Printf("Last Floor: %v\n", elevator.Last_Floor)
	fmt.Println("Requests")
	fmt.Println("Floor \t Hall Up \t Hall Down \t Cab")
	for i := 0; i < N_FLOORS; i++ {
		fmt.Printf("%v \t %v \t\t %v \t\t %v \t\n", i+1, elevator.Requests[i][0], elevator.Requests[i][1], elevator.Requests[i][2])
	}
}

func (elevator *Elevator) UpdateDirection(dir elevio.MotorDirection){
	elevio.SetMotorDirection(dir)
	elevator.Last_dir = dir
	elevator.Dirn = dir
}