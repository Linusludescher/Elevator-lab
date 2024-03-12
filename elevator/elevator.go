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
	//ElevatorNum int   `json:"ElevNum"`
}

type Elevator struct {
	Online      bool
	Behaviour   Behaviour
	Blocking    bool
	ElevNum     int
	Dirn        elevio.MotorDirection
	Last_dir    elevio.MotorDirection
	Last_Floor  int
	CabRequests []bool
}

type Worldview struct {
	ElevList     []Elevator
	Sender       int //hvorfooor?
	Version      uint64
	HallRequests [][2]uint8
}

func (w Worldview) Display() {
	fmt.Printf("Sender: %d\n", w.Sender)
	fmt.Printf("Version: %d\n", w.Version)

	fmt.Println("Elevator List:")
	for i, elev := range w.ElevList {
		fmt.Printf("  Elevator %d:\n", i+1)
		fmt.Printf("    Online: %v\n", elev.Online)
		fmt.Printf("    Behaviour: %v\n", elev.Behaviour)
		fmt.Printf("    Blocking: %t\n", elev.Blocking)
		fmt.Printf("    ElevNum: %d\n", elev.ElevNum)
		fmt.Printf("    Dirn: %v\n", elev.Dirn)
		fmt.Printf("    Last_dir: %v\n", elev.Last_dir)
		fmt.Printf("    Last_Floor: %d\n", elev.Last_Floor)
		fmt.Printf("    CabRequests: %v\n", elev.CabRequests)
	}

	fmt.Println("Hall Requests:")
	for _, request := range w.HallRequests {
		fmt.Printf("  Floor: %d, Direction: %d\n", request[0], request[1])
	}
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

func ElevatorInit(id int) (e Elevator, wv Worldview) {
	elevatorConfig := readElevatorConfig()
	hall := make([][2]uint8, elevatorConfig.N_FLOORS)
	for i := range hall {
		hall[i] = [2]uint8{0, 0}
	}
	cab := make([]bool, elevatorConfig.N_FLOORS)

	wv = Worldview{[]Elevator{}, e.ElevNum, startVersion, hall}

	for i := 1; i <= int(elevatorConfig.N_elevators); i++ {
		n := Elevator{false, EB_Idle, false, i, elevio.MD_Stop, elevio.MD_Stop, 0, cab}
		wv.ElevList = append(wv.ElevList, n)
	}
	e = wv.ElevList[id-1]

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
