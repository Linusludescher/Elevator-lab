package elevator

import (
	"encoding/json"
	"fmt"
	"os"
	"project/elevio"
	"time"
)

const (
	N_FLOORS     int    = 4
	N_BUTTONS    int    = 3
	startVersion uint64 = 5000
)

type Behaviour int

const (
	EB_Idle     Behaviour = 1
	EB_Moving   Behaviour = -1
	EB_DoorOpen Behaviour = 0
)

type ConfigData struct {
	N_FLOORS    uint8 `json:"Floors"`
	N_elevators uint8 `json:"n_elevators"`
	ElevatorNum int   `json:"ElevNum"`
}

type Elevator struct {
	Behaviour   Behaviour
	ElevNum     int
	Dirn        elevio.MotorDirection
	Last_dir    elevio.MotorDirection
	Last_Floor  int
	CabRequests []uint8
}

type Worldview struct {
	ElevList     []Elevator
	Sender       int
	Version      uint64
	HallRequests [][]uint8 //legge inn N_FLOOR (etter å ha lest) i første klamme
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

func ElevatorInit() (elevator Elevator, world Worldview) {
	elevatorConfig := readElevatorConfig()
	hall := make([][]uint8, 2)
	for i := range hall {
		hall[i] = make([]uint8, elevatorConfig.N_FLOORS)
	}
	cab := make([]uint8, elevatorConfig.N_FLOORS)

	elevator = Elevator{EB_Idle, elevatorConfig.ElevatorNum, elevio.MD_Stop, elevio.MD_Stop, 0, cab}

	world = Worldview{[]Elevator{elevator}, elevator.ElevNum, startVersion, hall}

	for elevio.GetFloor() != 0 {
		elevio.SetMotorDirection(elevio.MD_Down)
	}
	elevio.SetMotorDirection(elevio.MD_Stop)
	return
}

func (elevator *Elevator) Display() { //lage en for worldview også!
	fmt.Printf("Direction: %v\n", elevator.Dirn)
	fmt.Printf("Last Direction: %v\n", elevator.Last_dir)
	fmt.Printf("Last Floor: %v\n", elevator.Last_Floor)
	fmt.Println("Requests")
	fmt.Println("Floor\t Cab")
	for i := len(elevator.CabRequests) - 1; i >= 0; i-- {
		fmt.Printf("%v \t %v \t\n", i+1, elevator.CabRequests[i])
	}
}

func (elevator *Elevator) UpdateDirection(dir elevio.MotorDirection) {
	elevio.SetMotorDirection(dir)
	elevator.Last_dir = dir
	elevator.Dirn = dir
	if elevator.Dirn != elevio.MD_Stop {
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
