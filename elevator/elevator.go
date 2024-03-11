package elevator

import (
	"encoding/json"
	"fmt"
	"os"
	"project/elevio"
	"time"
)

const N_FLOORS int = 4
const N_BUTTONS int = 3

type ConfigData struct {
	N_FLOORS    uint8 `json:"Floors"`
	N_elevators uint8 `json:"n_elevators"`
	ElevatorNum int   `json:"ElevNum"`
}

type Behaviour int

const (
	EB_Idle     Behaviour = 1
	EB_Moving   Behaviour = -1
	EB_DoorOpen Behaviour = 0
)

type Elevator struct {
    Behaviour  Behaviour
	ElevNum    int
	Version    uint64
	Dirn       elevio.MotorDirection
	Last_dir   elevio.MotorDirection
	Last_Floor int
	Requests   [][]uint8
}

func readElevatorConfig() (elevatorData ConfigData) {
	jsonData, err := os.ReadFile("config.json")

	// can't read the config file, try again
	if err != nil {
		fmt.Printf("elevator.go: Error reading config file: %s\n", err)
		// tyr again
		readElevatorConfig()
	}

	// Parse jsonData into ElevatorData struct
	err = json.Unmarshal(jsonData, &elevatorData)

	// can't parse the config file, try again
	if err != nil {
		fmt.Printf("elevator.go: Error unmarshal json data to ElevatorData struct: %s\n", err)
		// tyr again
		readElevatorConfig()
	}
	return
}

func Elevator_uninitialized() (elevator Elevator) {
	elevatorConfig := readElevatorConfig()
	matrix := make([][]uint8, elevatorConfig.N_FLOORS)
	for i := range matrix {
		matrix[i] = make([]uint8, elevatorConfig.N_elevators+2)
	}

	elevator = Elevator{EB_Idle, elevatorConfig.ElevatorNum, 0, elevio.MD_Stop, elevio.MD_Stop, 0, matrix}
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
	fmt.Printf("Version: %v\n", elevator.Version)
	fmt.Println("Requests")
	fmt.Println("Floor \t Hall Up \t Hall Down \t Cab")
	for i := N_FLOORS - 1; i >= 0; i-- {
		fmt.Printf("%v \t %v \t\t %v \t\t %v \t\n", i+1, elevator.Requests[i][0], elevator.Requests[i][1], elevator.Requests[i][2])
	}
}

func (elevator *Elevator) UpdateDirection(dir elevio.MotorDirection) {
	elevio.SetMotorDirection(dir)
	elevator.Last_dir = dir
	elevator.Dirn = dir
    if elevator.Dirn != elevio.MD_Stop{
        elevator.Behaviour = EB_Moving
    } else {
        elevator.Behaviour = EB_Idle
    }
}

func BroadcastElevator(bc_chan chan bool, n_ms int) {
	for {
		bc_chan <- true
		time.Sleep(time.Duration(n_ms) * time.Millisecond)
	}
}
