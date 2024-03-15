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
	STARTVERSION uint64 = 5000
)

type Behaviour string

const (
	EB_IDLE      Behaviour = "idle"
	EB_MOVING    Behaviour = "moving"
	EB_DOOR_OPEN Behaviour = "doorOpen"
)

type ConfigData struct { //Reader config 2 steder: En gang her, en gang i network!
	N_FLOORS    uint8 `json:"Floors"`
	N_elevators uint8 `json:"n_elevators"`
	//ElevatorNum int   `json:"ElevNum"`
	UDPBase int `json:"BasePort"`
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
	VERSIONLIMIT  = 18446744073709551615 //2e64-1
	VERSIONBUFFER = 100000               //maks antall sykler ny versjon kan være foran for at e.Version settes godtar lavere p.Version (ved Version overflow)
	// versionInitVal = 10000 //initialisere på høyere verdi enn 0 for ikke problemer med nullstilling ved tilbakekobling etter utfall
)

func (worldView_p *Worldview) VersionUp() {
	if worldView_p.Version < VERSIONLIMIT {
		worldView_p.Version++
	} else {
		worldView_p.Version = 0
	}
}

func (worldview Worldview) Display() {
	fmt.Printf("Sender: %d\n", worldview.Sender)
	fmt.Printf("Version: %d\n", worldview.Version)

	fmt.Println("Elevator List:")
	for i, elev := range worldview.ElevList {
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
	for i := len(worldview.HallRequests); i > 0; i-- {
		fmt.Printf("floor: %d \thall up: %d\t, halldown: %d\n", i-1, worldview.HallRequests[i-1][0], worldview.HallRequests[i-1][1])
	}
}

func ReadElevatorConfig() (elevatorData ConfigData) {
	jsonData, err := os.ReadFile("config.json")

	// can't read the config file, try again
	if err != nil {
		fmt.Printf("elevator.go: Error reading config file: %s\n", err)
		// tyr again
		ReadElevatorConfig()
	}

	// Parse jsonData into ElevatorData struct
	err = json.Unmarshal(jsonData, &elevatorData)

	// can't parse the config file, try again
	if err != nil {
		fmt.Printf("elevator.go: Error unmarshal json data to ElevatorData struct: %s\n", err)
		// tyr again
		ReadElevatorConfig()
	}
	return
}

func ElevatorInit(id int) (elev Elevator, worldView Worldview) {
	elevatorConfig := ReadElevatorConfig()
	hallOrders := make([][2]uint8, elevatorConfig.N_FLOORS)
	for i := range hallOrders {
		hallOrders[i] = [2]uint8{0, 0}
	}

	worldView = Worldview{[]Elevator{}, id, STARTVERSION, hallOrders}

	for i := 1; i <= int(elevatorConfig.N_elevators); i++ {
		cabOrders := make([]bool, elevatorConfig.N_FLOORS)
		n := Elevator{false, true, EB_IDLE, false, i, elevio.MD_STOP, elevio.MD_STOP, 0, cabOrders}
		worldView.ElevList = append(worldView.ElevList, n)
	}
	elev = worldView.ElevList[id-1]
	elev.Online = true

	for elevio.GetFloor() != 0 {
		elevio.SetMotorDirection(elevio.MD_DOWN)
	}
	elevio.SetMotorDirection(elevio.MD_STOP)
	return
}

func (elev Elevator) Display() { //lage en for worldview også!
	fmt.Printf("Direction: %v\n", elev.Dirn)
	fmt.Printf("Last Direction: %v\n", elev.Last_dir)
	fmt.Printf("Last Floor: %v\n", elev.Last_Floor)
	fmt.Println("Requests")
	fmt.Println("Floor\t Cab")
	for i := len(elev.CabRequests) - 1; i >= 0; i-- {
		fmt.Printf("%v \t %v \t\n", i+1, elev.CabRequests[i])
	}
}

func (elev_p *Elevator) UpdateDirection(dir elevio.MotorDirection, wd_chan chan<- bool) {
	elevio.SetMotorDirection(dir)
	elev_p.Last_dir = dir
	elev_p.Dirn = dir
	if elev_p.Dirn != elevio.MD_STOP {
		elev_p.Behaviour = EB_MOVING
		wd_chan <- true
	} else {
		elev_p.Behaviour = EB_IDLE
	}
}

func BroadcastElevator(bc_chan chan<- bool, n_ms int) {
	for {
		bc_chan <- true
		time.Sleep(time.Duration(n_ms) * time.Millisecond)
	}
}

func UpdateLights(worldView Worldview, elevnum int) {
	for floor, f := range worldView.HallRequests {
		for buttonType, order := range f {
			elevio.SetButtonLamp(elevio.ButtonType(buttonType), floor, order != 0)
		}
	}
	for i, elev := range worldView.ElevList {
		if i+1 != elevnum {
			continue
		}
		for floor, f := range elev.CabRequests {
			elevio.SetButtonLamp(elevio.BT_CAB, floor, f) //TODO sende noe sånt <- elevio.ButtonLampOrder{1,2,true}

		}
	}
}
