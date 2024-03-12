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

type Behaviour string

const (
	EB_Idle     Behaviour = "idle"
	EB_Moving   Behaviour = "moving"
	EB_DoorOpen Behaviour = "doorOpen"
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
	CabRequests []bool
}

type Worldview struct {
	ElevList     []Elevator
	Sender       int
	Version      uint64
	HallRequests [][2]uint8 //legge inn N_FLOOR (etter å ha lest) i første klamme
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

func ElevatorInit() (e Elevator, world Worldview) {
	elevatorConfig := readElevatorConfig()
	hall := make([][2]uint8, elevatorConfig.N_FLOORS)
	for i := range hall {
		hall[i] = [2]uint8{0, 0}
	}
	cab := make([]bool, elevatorConfig.N_FLOORS)

	e = Elevator{EB_Idle, elevatorConfig.ElevatorNum, elevio.MD_Stop, elevio.MD_Stop, 0, cab}

	world = Worldview{[]Elevator{e}, e.ElevNum, startVersion, hall}

	for elevio.GetFloor() != 0 {
		elevio.SetMotorDirection(elevio.MD_Down)
	}
	elevio.SetMotorDirection(elevio.MD_Stop)
	return
}

func (e_p *Elevator) Display() { //lage en for worldview også!
	fmt.Printf("Direction: %v\n", e_p.Dirn)
	fmt.Printf("Last Direction: %v\n", e_p.Last_dir)
	fmt.Printf("Last Floor: %v\n", e_p.Last_Floor)
	fmt.Println("Requests")
	fmt.Println("Floor\t Cab")
	for i := len(e_p.CabRequests) - 1; i >= 0; i-- {
		fmt.Printf("%v \t %v \t\n", i+1, e_p.CabRequests[i])
	}
}

func (e_p *Elevator) UpdateDirection(dir elevio.MotorDirection) {
	elevio.SetMotorDirection(dir)
	e_p.Last_dir = dir
	e_p.Dirn = dir
	if e_p.Dirn != elevio.MD_Stop {
		e_p.Behaviour = EB_Moving
	} else {
		e_p.Behaviour = EB_Idle
	}
}

func BroadcastElevator(bc_chan chan bool, n_ms int) {
	for {
		bc_chan <- true
		time.Sleep(time.Duration(n_ms) * time.Millisecond)
	}
}
