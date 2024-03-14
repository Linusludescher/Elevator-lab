package requests

import (
	"project/costFunc"
	"project/elevator"
	"project/elevio"
	"time"
)

func RequestsAbove(e elevator.Elevator, wv elevator.Worldview) bool {
	for f := e.Last_Floor + 1; f < elevator.N_FLOORS; f++ {
		if e.CabRequests[f] {
			return true
		}
		if wv.HallRequests[f][elevio.BT_HallUp] == uint8(e.ElevNum) {
			return true
		}
		if wv.HallRequests[f][elevio.BT_HallDown] == uint8(e.ElevNum) {
			return true
		}
	}
	return false
}

func RequestsBelow(e elevator.Elevator, wv elevator.Worldview) bool {
	for f := 0; f < e.Last_Floor; f++ {
		if e.CabRequests[f] {
			return true
		}
		if wv.HallRequests[f][elevio.BT_HallUp] == uint8(e.ElevNum) {
			return true
		}
		if wv.HallRequests[f][elevio.BT_HallDown] == uint8(e.ElevNum) {
			return true
		}
	}
	return false
}

func RequestsHere(e elevator.Elevator, wv elevator.Worldview) bool { //kan bytte ut alle e.Floor med elevio.getFloor!!!!
	if e.CabRequests[e.Last_Floor] {
		return true
	}
	if wv.HallRequests[e.Last_Floor][elevio.BT_HallUp] == uint8(e.ElevNum) {
		return true
	}
	if wv.HallRequests[e.Last_Floor][elevio.BT_HallDown] == uint8(e.ElevNum) {
		return true
	}
	return false
}

func RequestsHereCabOrUp(e elevator.Elevator, wv elevator.Worldview) bool { // stygt, kan ores på en linje
	if wv.HallRequests[e.Last_Floor][elevio.BT_HallUp] == uint8(e.ElevNum) {
		return true
	}
	if e.CabRequests[e.Last_Floor] {
		return true
	}
	return false
}

func RequestsHereCabOrDown(e elevator.Elevator, wv elevator.Worldview) bool {
	if wv.HallRequests[e.Last_Floor][elevio.BT_HallDown] == uint8(e.ElevNum) {
		return true
	}
	if e.CabRequests[e.Last_Floor] {
		return true
	}
	return false
}

func DeleteOrdersHere(e_p *elevator.Elevator, wv_p *elevator.Worldview) {
	for orderType := 0; orderType < 2; orderType++ {
		if wv_p.HallRequests[e_p.Last_Floor][orderType] == uint8(e_p.ElevNum) {
			wv_p.HallRequests[e_p.Last_Floor][orderType] = 0
		}
	}
	e_p.CabRequests[e_p.Last_Floor] = false
	wv_p.Version_up()
}

func SetOrder(e_p *elevator.Elevator, wv_p *elevator.Worldview, buttn elevio.ButtonEvent) {
	if (buttn.Floor == e_p.Last_Floor) && (elevio.GetFloor() != -1) {
		return
	}
	if buttn.Button == elevio.BT_Cab {
		e_p.CabRequests[buttn.Floor] = true
	} else if wv_p.HallRequests[buttn.Floor][buttn.Button] == 0 {
		costFunc.CostFunction(wv_p, buttn)
	}
	wv_p.Version_up()
}

func ArrivedAtFloor(e_p *elevator.Elevator, wv_p *elevator.Worldview, reset_ch chan bool, wd_chan chan bool) {
	elevio.SetDoorOpenLamp(true)
	elevio.SetMotorDirection(elevio.MD_Stop)
	e_p.Dirn = elevio.MD_Stop
	DeleteOrdersHere(e_p, wv_p)
	wv_p.Version_up()
	e_p.Behaviour = elevator.EB_DoorOpen
	wd_chan <- true
	reset_ch <- true
}

func DisplayQueueCont(e_p *elevator.Elevator) {
	e_p.Display()
	time.Sleep(time.Second)
}
