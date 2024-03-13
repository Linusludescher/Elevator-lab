package stm

import (
	"fmt"
	"project/elevator"
	"project/elevio"
	"project/requests"
)

func TimerExp(e_p *elevator.Elevator, wv elevator.Worldview) {
	elevio.SetDoorOpenLamp(false)
	if e_p.Last_dir == elevio.MD_Up {
		if requests.RequestsAbove(*e_p, wv) {
			e_p.UpdateDirection(elevio.MD_Up)
		} else if requests.RequestsBelow(*e_p, wv) {
			e_p.UpdateDirection(elevio.MD_Down)
		} else {
			e_p.UpdateDirection(elevio.MD_Stop)
		}
	} else if e_p.Last_dir == elevio.MD_Down {
		if requests.RequestsBelow(*e_p, wv) {
			e_p.UpdateDirection(elevio.MD_Down)
		} else if requests.RequestsAbove(*e_p, wv) {
			e_p.UpdateDirection(elevio.MD_Up)
		} else {
			e_p.UpdateDirection(elevio.MD_Stop)
		}
	} else {
		e_p.UpdateDirection(elevio.MD_Stop)
	}
}

func ButtonPressed(e_p *elevator.Elevator, wv_p *elevator.Worldview, buttn elevio.ButtonEvent) {
	requests.SetOrder(e_p, wv_p, buttn)
}

func FloorSensed(e_p *elevator.Elevator, wv_p *elevator.Worldview, floor_sens int, timer_chan chan bool, obstruction_chan chan bool) {
	if floor_sens != -1 {
		e_p.Last_Floor = floor_sens
		fmt.Println("new floow")
	}
	if e_p.Dirn == elevio.MD_Up && floor_sens != -1 {
		if requests.RequestsHereCabOrUp(*e_p, *wv_p) {
			requests.ArrivedAtFloor(e_p, wv_p, timer_chan, obstruction_chan)
		} else if (!requests.RequestsAbove(*e_p, *wv_p)) && requests.RequestsHere(*e_p, *wv_p) {
			requests.ArrivedAtFloor(e_p, wv_p, timer_chan, obstruction_chan)
		}
	}
	if e_p.Dirn == elevio.MD_Down && floor_sens != -1 {
		if requests.RequestsHereCabOrDown(*e_p, *wv_p) {
			requests.ArrivedAtFloor(e_p, wv_p, timer_chan, obstruction_chan)
		} else if (!requests.RequestsBelow(*e_p, *wv_p)) && requests.RequestsHere(*e_p, *wv_p) {
			requests.ArrivedAtFloor(e_p, wv_p, timer_chan, obstruction_chan)
		}
	}
}

func Obstruction(e_p *elevator.Elevator,wv_p *elevator.Worldview, obstr bool) {
	e_p.Obstruction = obstr
	wv_p.Version++
}

func StopButtonPressed(e elevator.Elevator) {
	// fjerne hele køen?
	// vente ellerno?
}

func DefaultState(e_p *elevator.Elevator, wv_p *elevator.Worldview, broadcast_elevator_chan chan elevator.Worldview) {
	//e.Display()
	if (e_p.Dirn == elevio.MD_Stop) && (e_p.Behaviour != elevator.EB_DoorOpen) {
		if requests.RequestsAbove(*e_p, *wv_p) {
			e_p.UpdateDirection(elevio.MD_Up)
		} else if requests.RequestsBelow(*e_p, *wv_p) {
			e_p.UpdateDirection(elevio.MD_Down)
		}
	}
}
