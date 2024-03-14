package stm

import (
	"project/elevator"
	"project/elevio"
	"project/requests"
)

func ClosingDoor(e_p *elevator.Elevator, wv elevator.Worldview, wd_chan chan bool) { //kalle denne for door closed
	wd_chan <- false
	elevio.SetDoorOpenLamp(false)
	if e_p.Last_dir == elevio.MD_Up {
		if requests.RequestsAbove(*e_p, wv) {
			e_p.UpdateDirection(elevio.MD_Up, wd_chan)
		} else if requests.RequestsBelow(*e_p, wv) {
			e_p.UpdateDirection(elevio.MD_Down, wd_chan)
		} else {
			e_p.UpdateDirection(elevio.MD_Stop, wd_chan)
		}
	} else if e_p.Last_dir == elevio.MD_Down {
		if requests.RequestsBelow(*e_p, wv) {
			e_p.UpdateDirection(elevio.MD_Down, wd_chan)
		} else if requests.RequestsAbove(*e_p, wv) {
			e_p.UpdateDirection(elevio.MD_Up, wd_chan)
		} else {
			e_p.UpdateDirection(elevio.MD_Stop, wd_chan)
		}
	} else {
<<<<<<< HEAD
		e.Behaviour = elevator.EB_Idle
		e.Dirn = elevio.MD_Stop
=======
		e_p.UpdateDirection(elevio.MD_Stop, wd_chan)
>>>>>>> main
	}
}

func ButtonPressed(e_p *elevator.Elevator, wv_p *elevator.Worldview, buttn elevio.ButtonEvent) {
	requests.SetOrder(e_p, wv_p, buttn)
}

func FloorSensed(e_p *elevator.Elevator, wv_p *elevator.Worldview, floor_sens int, resetTimer_chan chan bool, wd_chan chan bool) {
	wd_chan <- false

	if floor_sens != -1 {
		e_p.Last_Floor = floor_sens
		elevio.SetFloorIndicator(floor_sens)
	}
<<<<<<< HEAD
	if e.Dirn == elevio.MD_Up && floor_sens != -1 {
		if requests.RequestsHereCabOrUp(*e) {
			fmt.Println("stopping")
			requests.ArrivedAtFloor(e, timer_chan)
		} else if (!requests.RequestsAbove(*e)) && requests.RequestsHere(*e) {
			requests.ArrivedAtFloor(e, timer_chan)
		}
	}
	if e.Dirn == elevio.MD_Down && floor_sens != -1 {
		if requests.RequestsHereCabOrDown(*e) {
			fmt.Println("stopping")
			requests.ArrivedAtFloor(e, timer_chan)
		} else if (!requests.RequestsBelow(*e)) && requests.RequestsHere(*e) {
			requests.ArrivedAtFloor(e, timer_chan)
=======
	if e_p.Dirn == elevio.MD_Up && floor_sens != -1 {
		if requests.RequestsHereCabOrUp(*e_p, *wv_p) {
			requests.ArrivedAtFloor(e_p, wv_p, resetTimer_chan, wd_chan)
		} else if (!requests.RequestsAbove(*e_p, *wv_p)) && requests.RequestsHere(*e_p, *wv_p) {
			requests.ArrivedAtFloor(e_p, wv_p, resetTimer_chan, wd_chan)
		}
	}
	if e_p.Dirn == elevio.MD_Down && floor_sens != -1 {
		if requests.RequestsHereCabOrDown(*e_p, *wv_p) {
			requests.ArrivedAtFloor(e_p, wv_p, resetTimer_chan, wd_chan)
		} else if (!requests.RequestsBelow(*e_p, *wv_p)) && requests.RequestsHere(*e_p, *wv_p) {
			requests.ArrivedAtFloor(e_p, wv_p, resetTimer_chan, wd_chan)
>>>>>>> main
		}
	}
	// if (floor_sens == -1 && e_p.Last_dir == elevio.MD_Down && e_p.Last_Floor == 0) || (floor_sens == -1 && e_p.Last_dir == elevio.MD_Up && e_p.Last_Floor == 3) {
	// 	e_p.UpdateDirection(elevio.MD_Stop, wd_chan)
	// }
}

<<<<<<< HEAD
func Obstruction(e elevator.Elevator, obstr bool) {
	if obstr {
		elevio.SetMotorDirection(elevio.MD_Stop)
	} else {
		elevio.SetMotorDirection(e.Dirn)
	}
=======
func Obstruction(e_p *elevator.Elevator, wv_p *elevator.Worldview, obstr bool) {
	e_p.Obstruction = obstr
	wv_p.Version_up()
>>>>>>> main
}

func StopButtonPressed(e elevator.Elevator) {
	// fjerne hele køen?
	// vente ellerno?
}

func DefaultState(e_p *elevator.Elevator, wv_p *elevator.Worldview, resetTimer_chan chan bool, wd_chan chan bool) {
	go aloneUpdateLights(*wv_p, *e_p)
	for floor := range wv_p.HallRequests {
		for _, l := range wv_p.HallRequests[floor] {
			if l == uint8(e_p.ElevNum) && floor == e_p.Last_Floor && uint8(elevio.GetFloor()) == l {
				requests.ArrivedAtFloor(e_p, wv_p, resetTimer_chan, wd_chan)
			}
		}
	}
	if (e_p.Dirn == elevio.MD_Stop) && (e_p.Behaviour != elevator.EB_DoorOpen) {
		if requests.RequestsAbove(*e_p, *wv_p) {
			e_p.UpdateDirection(elevio.MD_Up, wd_chan)
		} else if requests.RequestsBelow(*e_p, *wv_p) {
			e_p.UpdateDirection(elevio.MD_Down, wd_chan)
		}
	}
}

func aloneUpdateLights(wv elevator.Worldview, e elevator.Elevator) {
	// Lys med kun en heis
	only_elev := true
	for i := range wv.ElevList {
		if wv.ElevList[i].Online && wv.ElevList[i].ElevNum != e.ElevNum {
			only_elev = false
		}
	}
	if only_elev {
		elevator.UpdateLights(wv, e.ElevNum)
	}
	// Lys med en heis ferdig
}
