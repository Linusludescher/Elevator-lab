package stm

import (
	"project/elevator"
	"project/elevio"
	"project/requests"
)

func ClosingDoor(elev_p *elevator.Elevator, worldView elevator.Worldview, wd_chan chan bool) { //kalle denne for door closed
	wd_chan <- false
	elevio.SetDoorOpenLamp(false)
	if elev_p.Last_dir == elevio.MD_UP {
		if requests.RequestsAbove(*elev_p, worldView) {
			elev_p.UpdateDirection(elevio.MD_UP, wd_chan)
		} else if requests.RequestsBelow(*elev_p, worldView) {
			elev_p.UpdateDirection(elevio.MD_DOWN, wd_chan)
		} else {
			elev_p.UpdateDirection(elevio.MD_STOP, wd_chan)
		}
	} else if elev_p.Last_dir == elevio.MD_DOWN {
		if requests.RequestsBelow(*elev_p, worldView) {
			elev_p.UpdateDirection(elevio.MD_DOWN, wd_chan)
		} else if requests.RequestsAbove(*elev_p, worldView) {
			elev_p.UpdateDirection(elevio.MD_UP, wd_chan)
		} else {
			elev_p.UpdateDirection(elevio.MD_STOP, wd_chan)
		}
	} else {
		elev_p.UpdateDirection(elevio.MD_STOP, wd_chan)
	}
}

func ButtonPressed(elev_p *elevator.Elevator, worldView_p *elevator.Worldview, buttn elevio.ButtonEvent, reset_timer_chan chan bool, wd_chan chan bool) {
	requests.SetOrder(elev_p, worldView_p, buttn, reset_timer_chan, wd_chan)
}

func FloorSensed(elev_p *elevator.Elevator, worldView_p *elevator.Worldview, floor_sens int, reset_timer_chan chan bool, wd_chan chan<- bool) {
	wd_chan <- false

	if floor_sens != -1 {
		elev_p.Last_Floor = floor_sens
		elevio.SetFloorIndicator(floor_sens)
	}
	if elev_p.Dirn == elevio.MD_UP && floor_sens != -1 {
		if requests.RequestsHereCabOrUp(*elev_p, *worldView_p) {
			requests.ArrivedAtFloor(elev_p, worldView_p, reset_timer_chan, wd_chan)
		} else if (!requests.RequestsAbove(*elev_p, *worldView_p)) && requests.RequestsHere(*elev_p, *worldView_p) {
			requests.ArrivedAtFloor(elev_p, worldView_p, reset_timer_chan, wd_chan)
		}
	}
	if elev_p.Dirn == elevio.MD_DOWN && floor_sens != -1 {
		if requests.RequestsHereCabOrDown(*elev_p, *worldView_p) {
			requests.ArrivedAtFloor(elev_p, worldView_p, reset_timer_chan, wd_chan)
		} else if (!requests.RequestsBelow(*elev_p, *worldView_p)) && requests.RequestsHere(*elev_p, *worldView_p) {
			requests.ArrivedAtFloor(elev_p, worldView_p, reset_timer_chan, wd_chan)
		}
	}
	//softstop
	if (floor_sens == -1 && elev_p.Last_dir == elevio.MD_DOWN && elev_p.Last_Floor == 0) || (floor_sens == -1 && elev_p.Last_dir == elevio.MD_UP && elev_p.Last_Floor == 3) {
		elev_p.UpdateDirection(elevio.MD_STOP, wd_chan)
	}
}

func Obstruction(elev_p *elevator.Elevator, worldView_p *elevator.Worldview, obstr bool) {
	if obstr && elev_p.Behaviour == elevator.EB_DOOR_OPEN {
		elev_p.Obstruction = obstr
		worldView_p.VersionUp()
	} else if !obstr {
		elev_p.Obstruction = obstr
		worldView_p.VersionUp()
	}
}

func DefaultState(elev_p *elevator.Elevator, worldView_p *elevator.Worldview, reset_timer_chan chan bool, wd_chan chan bool) {
	go aloneUpdateLights(*worldView_p, *elev_p)
	for floor := range worldView_p.HallRequests {
		for _, order := range worldView_p.HallRequests[floor] {
			if order == uint8(elev_p.ElevNum) && floor == elev_p.Last_Floor && elevio.GetFloor() == floor {
				requests.ArrivedAtFloor(elev_p, worldView_p, reset_timer_chan, wd_chan)
			}
		}
	}
	if (elev_p.Dirn == elevio.MD_STOP) && (elev_p.Behaviour != elevator.EB_DOOR_OPEN) {
		if requests.RequestsAbove(*elev_p, *worldView_p) {
			elev_p.UpdateDirection(elevio.MD_UP, wd_chan)
		} else if requests.RequestsBelow(*elev_p, *worldView_p) {
			elev_p.UpdateDirection(elevio.MD_DOWN, wd_chan)
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
