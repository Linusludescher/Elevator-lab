package worldview

import (
	"encoding/json"
	"fmt"
	"os"
	"project/elevio"
	"time"
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
	UDPBase     int   `json:"BasePort"`
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

func ElevatorInit(timer_exp_chan chan bool, id int) (elev Elevator, worldView Worldview) {
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
	timer_exp_chan <- true
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

func (elev_p *Elevator) UpdateDirection(dir elevio.MotorDirection, watchdog_chan chan<- bool) {
	elevio.SetMotorDirection(dir)
	elev_p.Last_dir = dir
	elev_p.Dirn = dir
	if elev_p.Dirn != elevio.MD_STOP {
		elev_p.Behaviour = EB_MOVING
		watchdog_chan <- true
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