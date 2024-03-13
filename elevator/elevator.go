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
	Operative   bool
	Behaviour   Behaviour
	Obstruction bool
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
	HallRequests [][2]uint8
}

const (
	V_l   = 18446744073709551615 //2e64-1
	V_s_c = 100000               //maks antall sykler ny versjon kan være foran for at e.Version settes godtar lavere p.Version (ved Version overflow)
	// versionInitVal = 10000 //initialisere på høyere verdi enn 0 for ikke problemer med nullstilling ved tilbakekobling etter utfall
)

func (wv *Worldview) Version_up() {
	if wv.Version < V_l {
		wv.Version++
	} else {
		wv.Version = 0
	}
}

func (w Worldview) Display() {
	fmt.Printf("Sender: %d\n", w.Sender)
	fmt.Printf("Version: %d\n", w.Version)

	fmt.Println("Elevator List:")
	for i, elev := range w.ElevList {
		fmt.Printf("  Elevator %d:\n", i+1)
		fmt.Printf("    Online: %v\n", elev.Online)
		fmt.Printf("    Operative: %v\n", elev.Operative)
		fmt.Printf("    Behaviour: %v\n", elev.Behaviour)
		fmt.Printf("    Obstruction: %t\n", elev.Obstruction)
		fmt.Printf("    ElevNum: %d\n", elev.ElevNum)
		fmt.Printf("    Dirn: %v\n", elev.Dirn)
		fmt.Printf("    Last_dir: %v\n", elev.Last_dir)
		fmt.Printf("    Last_Floor: %d\n", elev.Last_Floor)
		fmt.Printf("    CabRequests: %v\n", elev.CabRequests)
	}

	fmt.Println("Hall Requests:")
	for i := len(w.HallRequests); i > 0; i-- {
		fmt.Printf("floor: %d \thall up: %d\t, halldown: %d\n", i-1, w.HallRequests[i-1][0], w.HallRequests[i-1][1])
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

	wv = Worldview{[]Elevator{}, id, startVersion, hall}

	for i := 1; i <= int(elevatorConfig.N_elevators); i++ {
		cab := make([]bool, elevatorConfig.N_FLOORS)
		n := Elevator{false, true, EB_Idle, false, i, elevio.MD_Stop, elevio.MD_Stop, 0, cab}
		wv.ElevList = append(wv.ElevList, n)
	}
	e = wv.ElevList[id-1]
	e.Online = true

	for elevio.GetFloor() != 0 {
		elevio.SetMotorDirection(elevio.MD_Down)
	}
	elevio.SetMotorDirection(elevio.MD_Stop)
	return
}

func (e Elevator) Display() { //lage en for worldview også!
	fmt.Printf("Direction: %v\n", e.Dirn)
	fmt.Printf("Last Direction: %v\n", e.Last_dir)
	fmt.Printf("Last Floor: %v\n", e.Last_Floor)
	fmt.Println("Requests")
	fmt.Println("Floor\t Cab")
	for i := len(e.CabRequests) - 1; i >= 0; i-- {
		fmt.Printf("%v \t %v \t\n", i+1, e.CabRequests[i])
	}
}

func (e_p *Elevator) UpdateDirection(dir elevio.MotorDirection, wd_chan chan bool) {
	elevio.SetMotorDirection(dir)
	e_p.Last_dir = dir
	e_p.Dirn = dir
	if e_p.Dirn != elevio.MD_Stop {
		e_p.Behaviour = EB_Moving
		fmt.Println("update før sending")
		wd_chan <- true
		fmt.Println("etter sending")
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
