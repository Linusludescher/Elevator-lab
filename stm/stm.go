package stm

import (
	"fmt"
	"project/elevator"
	"project/elevio"
	"project/requests"
)

func TimerState(e *elevator.Elevator) {
	if e.Last_dir == elevio.MD_Up {
		if requests.RequestsAbove(*e) {
			e.UpdateDirection(elevio.MD_Up)
		} else if requests.RequestsBelow(*e) {
			e.UpdateDirection(elevio.MD_Down)
		}
	} else if e.Last_dir == elevio.MD_Down {
		if requests.RequestsBelow(*e) {
			e.UpdateDirection(elevio.MD_Down)
		} else if requests.RequestsAbove(*e) {
			e.UpdateDirection(elevio.MD_Up)
		}
	} else {
		e.Dirn = elevio.MD_Stop
	}
}

func ButtonPressed(e *elevator.Elevator, buttn elevio.ButtonEvent) {
	requests.SetOrderHere(e, buttn) // tuple her etterhvert
	e.Display()
}

func FloorSensed(e *elevator.Elevator, floor_sens int, timer_chan chan bool) {
	if floor_sens != -1 {
		e.Last_Floor = floor_sens
		fmt.Println("new floow")
	}
	if e.Dirn == elevio.MD_Up && floor_sens != -1 {
		fmt.Println("Test4")
		if requests.RequestsHereCabOrUp(*e) {
			fmt.Println("stopping")
			requests.ArrivedAtFloor(e, timer_chan)
		} else if (!requests.RequestsAbove(*e)) && requests.RequestsHere(*e) {
			requests.ArrivedAtFloor(e, timer_chan)
		}
	}
	if e.Dirn == elevio.MD_Down && floor_sens != -1 {
		fmt.Println("Test5")
		if requests.RequestsHereCabOrDown(*e) {
			fmt.Println("stopping")
			requests.ArrivedAtFloor(e, timer_chan)
		} else if (!requests.RequestsBelow(*e)) && requests.RequestsHere(*e) {
			requests.ArrivedAtFloor(e, timer_chan)
		}
	}
}

func Obstuction(e elevator.Elevator, obstr bool) {
	if obstr {
		elevio.SetMotorDirection(elevio.MD_Stop)
	} else {
		elevio.SetMotorDirection(e.Dirn)
	}
}

func StopButtonPressed(e elevator.Elevator) {
	// fjerne hele køen?
	requests.DeleteAllOrdes(&e)

	// vente ellerno?
}

func DefaultState(e *elevator.Elevator, broadcast_elevator_chan chan elevator.Elevator) {
	broadcast_elevator_chan <- *e
	//e.Display()
	if e.Dirn == elevio.MD_Stop {
		if requests.RequestsAbove(*e) {
			fmt.Printf("test2\n")
			e.UpdateDirection(elevio.MD_Up)
		} else if requests.RequestsBelow(*e) {
			fmt.Printf("test3\n")
			e.UpdateDirection(elevio.MD_Down)
		}
	}
}
