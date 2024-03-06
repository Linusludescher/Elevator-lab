package elevator

import (
	"encoding/json"
	"fmt"
	"os"
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

type ConfigData struct {
	N_FLOORS    uint8 `json:"Floors"`
	N_BUTTONS   uint8 `json:"Buttons"`
	ElevatorNum int   `json:"ElevNum2"`
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

type Elevator struct {
	ElevNum    int
	Version    uint64
	Dirn       elevio.MotorDirection
	Last_dir   elevio.MotorDirection
	Last_Floor int
	Requests   [][]uint8
}

func Elevator_uninitialized() (elevator Elevator) {
	elevatorConfig := readElevatorConfig()
	matrix := make([][]uint8, elevatorConfig.N_FLOORS)
	for i := range matrix {
		matrix[i] = make([]uint8, elevatorConfig.N_BUTTONS)
	}

	elevator = Elevator{elevatorConfig.ElevatorNum, 0, elevio.MD_Stop, elevio.MD_Stop, 0, matrix}
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

func (elevator *Elevator) UpdateDirection(dir elevio.MotorDirection) {
	elevio.SetMotorDirection(dir)
	elevator.Last_dir = dir
	elevator.Dirn = dir
}
