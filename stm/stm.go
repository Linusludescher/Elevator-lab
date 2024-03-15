package stm

import (
	"project/elevator"
	"project/elevio"
	"project/requests"
)

func ClosingDoor(elev_p *elevator.Elevator, worldView elevator.Worldview, wd_chan chan bool) { //kalle denne for door closed
	wd_chan <- false
	elevio.SetDoorOpenLamp(false)
	if elev_p.Last_dir == elevio.MD_Up {
		if requests.RequestsAbove(*elev_p, worldView) {
			elev_p.UpdateDirection(elevio.MD_Up, wd_chan)
		} else if requests.RequestsBelow(*elev_p, worldView) {
			elev_p.UpdateDirection(elevio.MD_Down, wd_chan)
		} else {
			elev_p.UpdateDirection(elevio.MD_Stop, wd_chan)
		}
	} else if elev_p.Last_dir == elevio.MD_Down {
		if requests.RequestsBelow(*elev_p, worldView) {
			elev_p.UpdateDirection(elevio.MD_Down, wd_chan)
		} else if requests.RequestsAbove(*elev_p, worldView) {
			elev_p.UpdateDirection(elevio.MD_Up, wd_chan)
		} else {
			elev_p.UpdateDirection(elevio.MD_Stop, wd_chan)
		}
	} else {
		elev_p.UpdateDirection(elevio.MD_Stop, wd_chan)
	}
}

func ButtonPressed(elev_p *elevator.Elevator, worldView_p *elevator.Worldview, buttn elevio.ButtonEvent, resetTimer_chan chan bool, wd_chan chan bool) {
	requests.SetOrder(elev_p, worldView_p, buttn, resetTimer_chan, wd_chan)
}

func FloorSensed(elev_p *elevator.Elevator, worldView_p *elevator.Worldview, floor_sens int, resetTimer_chan chan bool, wd_chan chan<- bool) {
	wd_chan <- false

	if floor_sens != -1 {
		elev_p.Last_Floor = floor_sens
		elevio.SetFloorIndicator(floor_sens)
	}
	if elev_p.Dirn == elevio.MD_Up && floor_sens != -1 {
		if requests.RequestsHereCabOrUp(*elev_p, *worldView_p) {
			requests.ArrivedAtFloor(elev_p, worldView_p, resetTimer_chan, wd_chan)
		} else if (!requests.RequestsAbove(*elev_p, *worldView_p)) && requests.RequestsHere(*elev_p, *worldView_p) {
			requests.ArrivedAtFloor(elev_p, worldView_p, resetTimer_chan, wd_chan)
		}
	}
	if elev_p.Dirn == elevio.MD_Down && floor_sens != -1 {
		if requests.RequestsHereCabOrDown(*elev_p, *worldView_p) {
			requests.ArrivedAtFloor(elev_p, worldView_p, resetTimer_chan, wd_chan)
		} else if (!requests.RequestsBelow(*elev_p, *worldView_p)) && requests.RequestsHere(*elev_p, *worldView_p) {
			requests.ArrivedAtFloor(elev_p, worldView_p, resetTimer_chan, wd_chan)
		}
	}
	//softstop
	if (floor_sens == -1 && elev_p.Last_dir == elevio.MD_Down && elev_p.Last_Floor == 0) || (floor_sens == -1 && elev_p.Last_dir == elevio.MD_Up && elev_p.Last_Floor == 3) {
		elev_p.UpdateDirection(elevio.MD_Stop, wd_chan)
	}
}

func Obstruction(elev_p *elevator.Elevator, worldView_p *elevator.Worldview, obstr bool) {
	if obstr && elev_p.Behaviour == elevator.EB_DoorOpen {
		elev_p.Obstruction = obstr
		worldView_p.Version_up()
	} else if !obstr {
		elev_p.Obstruction = obstr
		worldView_p.Version_up()
	}
}

func DefaultState(elev_p *elevator.Elevator, worldView_p *elevator.Worldview, resetTimer_chan chan bool, wd_chan chan bool) {
	go aloneUpdateLights(*worldView_p, *elev_p)
	for floor := range worldView_p.HallRequests {
		for _, order := range worldView_p.HallRequests[floor] {
			if order == uint8(elev_p.ElevNum) && floor == elev_p.Last_Floor && elevio.GetFloor() == floor {
				requests.ArrivedAtFloor(elev_p, worldView_p, resetTimer_chan, wd_chan)
			}
		}
	}
	if (elev_p.Dirn == elevio.MD_Stop) && (elev_p.Behaviour != elevator.EB_DoorOpen) {
		if requests.RequestsAbove(*elev_p, *worldView_p) {
			elev_p.UpdateDirection(elevio.MD_Up, wd_chan)
		} else if requests.RequestsBelow(*elev_p, *worldView_p) {
			elev_p.UpdateDirection(elevio.MD_Down, wd_chan)
		}
	}
}

func aloneUpdateLights(worldView elevator.Worldview, elev elevator.Elevator) {
	only_elev := true
	for i := range worldView.ElevList {
		if worldView.ElevList[i].Online && worldView.ElevList[i].ElevNum != elev.ElevNum {
			only_elev = false
		}
	}
	if only_elev {
		elevator.UpdateLights(worldView, elev.ElevNum)
	}
}
