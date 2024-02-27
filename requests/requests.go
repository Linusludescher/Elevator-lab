package requests

import (
	"fmt"
	"project/elevator"
	"project/elevio"
)

func RequestsAbove(e elevator.Elevator) bool {
	for f := e.Last_Floor + 1; f < elevator.N_FLOORS; f++ {
		for btn := 0; btn < elevator.N_BUTTONS; btn++ {
			if e.Requests[f][btn] == 1 {
				return true
			}
		}
	}
	return false
}

func RequestsBelow(e elevator.Elevator) bool {
	for f := 0; f < e.Last_Floor; f++ {
		for btn := 0; btn < elevator.N_BUTTONS; btn++ {
			if e.Requests[f][btn] == 1 {
				return true
			}
		}
	}
	return false
}

func RequestsHere(e elevator.Elevator) bool { //kan bytte ut alle e.Floor med elevio.getFloor!!!!
	for btn := 0; btn < elevator.N_BUTTONS; btn++ {
		if e.Requests[e.Last_Floor][btn] == 1 {
			return true
		}
	}
	return false
}

func RequestsHereCabOrUp(e elevator.Elevator) bool { // stygt, kan ores på en linje
	if e.Requests[e.Last_Floor][elevio.BT_HallUp] == 1 {
		return true
	}
	if e.Requests[e.Last_Floor][elevio.BT_Cab] == 1 {
		return true
	}
	return false
}

func RequestsHereCabOrDown(e elevator.Elevator) bool {
	if e.Requests[e.Last_Floor][elevio.BT_HallDown] == 1 {
		return true
	}
	if e.Requests[e.Last_Floor][elevio.BT_Cab] == 1 {
		return true
	}
	return false
}

func DeleteOrdersHere(e *elevator.Elevator) {
	for orderType := 0; orderType < 3; orderType++ {
		e.Requests[e.Last_Floor][orderType] = 0
		elevio.SetButtonLamp(elevio.ButtonType(orderType), e.Last_Floor, false)
	}
}

func DeleteAllOrdes(e *elevator.Elevator) {
	for floor := 0; floor < 4; floor++ {
		for orderType := 0; orderType < 3; orderType++ {
			e.Requests[floor][orderType] = 0
			elevio.SetButtonLamp(elevio.ButtonType(orderType), floor, false)
		}
	}
}

func SetOrderHere(e *elevator.Elevator, buttn elevio.ButtonEvent) {
	if buttn.Floor == e.Last_Floor {
		return
	}
	e.Requests[buttn.Floor][buttn.Button] = 1
	elevio.SetButtonLamp(buttn.Button, buttn.Floor, true) // må endres når flere heiser
}

func PrintRequests(e elevator.Elevator) {
	for f := 0; f < elevator.N_FLOORS; f++ {
		fmt.Printf("\n")
		for btn := 0; btn < elevator.N_BUTTONS; btn++ {
			fmt.Print(e.Requests[f][btn])
		}
	}
}
