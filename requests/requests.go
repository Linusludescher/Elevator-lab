package requests

import (
	"fmt"
	"project/elevator"
	"project/elevio"
	"project/timer"
	"time"
)

func RequestsAbove(e elevator.Elevator) bool {
	for f := e.Last_Floor + 1; f < elevator.N_FLOORS; f++ {
		for btn := 0; btn < elevator.N_BUTTONS; btn++ {
			if e.Requests[f][btn] == uint8(e.ElevNum) {
				return true
			}
		}
	}
	return false
}

func RequestsBelow(e elevator.Elevator) bool {
	for f := 0; f < e.Last_Floor; f++ {
		for btn := 0; btn < elevator.N_BUTTONS; btn++ {
			if e.Requests[f][btn] == uint8(e.ElevNum) {
				return true
			}
		}
	}
	return false
}

func RequestsHere(e elevator.Elevator) bool { //kan bytte ut alle e.Floor med elevio.getFloor!!!!
	for btn := 0; btn < elevator.N_BUTTONS; btn++ {
		if e.Requests[e.Last_Floor][btn] == uint8(e.ElevNum) {
			return true
		}
	}
	return false
}

func RequestsHereCabOrUp(e elevator.Elevator) bool { // stygt, kan ores på en linje
	if e.Requests[e.Last_Floor][elevio.BT_HallUp] == uint8(e.ElevNum) {
		return true
	}
	if e.Requests[e.Last_Floor][elevio.BT_Cab] == uint8(e.ElevNum) {
		return true
	}
	return false
}

func RequestsHereCabOrDown(e elevator.Elevator) bool {
	if e.Requests[e.Last_Floor][elevio.BT_HallDown] == uint8(e.ElevNum) {
		return true
	}
	if e.Requests[e.Last_Floor][elevio.BT_Cab] == uint8(e.ElevNum) {
		return true
	}
	return false
}

func DeleteOrdersHere(e *elevator.Elevator) {
	for orderType := 0; orderType < 3; orderType++ {
		e.Requests[e.Last_Floor][orderType] = 0
		elevio.SetButtonLamp(elevio.ButtonType(orderType), e.Last_Floor, false)
	}
	e.Version++
}

func DeleteAllOrdes(e *elevator.Elevator) {
	for floor := 0; floor < 4; floor++ {
		for orderType := 0; orderType < 3; orderType++ {
			e.Requests[floor][orderType] = 0
			elevio.SetButtonLamp(elevio.ButtonType(orderType), floor, false)
		}
	}
	e.Version++
}

func SetOrderHere(e *elevator.Elevator, buttn elevio.ButtonEvent) {
	if buttn.Floor == e.Last_Floor {
		return
	}
	e.Requests[buttn.Floor][buttn.Button] = 1
	elevio.SetButtonLamp(buttn.Button, buttn.Floor, true) // må endres når flere heiser
	e.Version++
}

func PrintRequests(e elevator.Elevator) {
	for f := 0; f < elevator.N_FLOORS; f++ {
		fmt.Printf("\n")
		for btn := 0; btn < elevator.N_BUTTONS; btn++ {
			fmt.Print(e.Requests[f][btn])
		}
	}
}

func ArrivedAtFloor(e *elevator.Elevator, timer_chan chan bool) {
	elevio.SetMotorDirection(elevio.MD_Stop)
	e.Dirn = elevio.MD_Stop
	DeleteOrdersHere(e)
	e.Version++
	e.Behaviour = elevator.EB_DoorOpen
	go timer.TimerStart(3, timer_chan)
}

func DisplayQueueCont(e *elevator.Elevator){
	e.Display()
	time.Sleep(time.Second)
}
