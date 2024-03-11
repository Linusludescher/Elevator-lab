package stm

import (
	"fmt"
	"project/elevator"
	"project/elevio"
	"project/requests"
)

func TimerState(e *elevator.Elevator, wv elevator.Worldview) {
	if e.Last_dir == elevio.MD_Up {
		if requests.RequestsAbove(*e, wv) {
			e.UpdateDirection(elevio.MD_Up)
		} else if requests.RequestsBelow(*e, wv) {
			e.UpdateDirection(elevio.MD_Down)
		}
	} else if e.Last_dir == elevio.MD_Down {
		if requests.RequestsBelow(*e, wv) {
			e.UpdateDirection(elevio.MD_Down)
		} else if requests.RequestsAbove(*e, wv) {
			e.UpdateDirection(elevio.MD_Up)
		}
	} else {
		e.Dirn = elevio.MD_Stop
	}
}

func ButtonPressed(e *elevator.Elevator, wv *elevator.Worldview, buttn elevio.ButtonEvent) {
	requests.SetOrderHere(e, wv, buttn)
}

func FloorSensed(e *elevator.Elevator, wv *elevator.Worldview, floor_sens int, timer_chan chan bool) {
	if floor_sens != -1 {
		e.Last_Floor = floor_sens
		fmt.Println("new floow")
	}
	if e.Dirn == elevio.MD_Up && floor_sens != -1 {
		fmt.Println("Test4")
		if requests.RequestsHereCabOrUp(*e, *wv) {
			fmt.Println("stopping")
			requests.ArrivedAtFloor(e, wv, timer_chan)
		} else if (!requests.RequestsAbove(*e, *wv)) && requests.RequestsHere(*e, *wv) {
			requests.ArrivedAtFloor(e, wv, timer_chan)
		}
	}
	if e.Dirn == elevio.MD_Down && floor_sens != -1 {
		fmt.Println("Test5")
		if requests.RequestsHereCabOrDown(*e, *wv) {
			fmt.Println("stopping")
			requests.ArrivedAtFloor(e, wv, timer_chan)
		} else if (!requests.RequestsBelow(*e, *wv)) && requests.RequestsHere(*e, *wv) {
			requests.ArrivedAtFloor(e, wv, timer_chan)
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
	// vente ellerno?
}

func DefaultState(e *elevator.Elevator, wv *elevator.Worldview, broadcast_elevator_chan chan elevator.Worldview) {
	//e.Display()
	if e.Dirn == elevio.MD_Stop {
		if requests.RequestsAbove(*e, *wv) {
			fmt.Printf("test2\n")
			e.UpdateDirection(elevio.MD_Up)
		} else if requests.RequestsBelow(*e, *wv) {
			fmt.Printf("test3\n")
			e.UpdateDirection(elevio.MD_Down)
		}
	}
}
