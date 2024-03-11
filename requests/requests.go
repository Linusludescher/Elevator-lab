package requests

import (
	"project/elevator"
	"project/elevio"
	"project/timer"
	"time"
)

func RequestsAbove(e elevator.Elevator, wv elevator.Worldview) bool {
	for f := e.Last_Floor + 1; f < elevator.N_FLOORS; f++ {
		if e.CabRequests[f] == 1 {
			return true
		}
		if wv.HallRequests[elevio.BT_HallUp][f] > 0 {
			return true
		}
		if wv.HallRequests[elevio.BT_HallDown][f] > 0 {
			return true
		}
	}
	return false
}

func RequestsBelow(e elevator.Elevator, wv elevator.Worldview) bool {
	for f := 0; f < e.Last_Floor; f++ {
		if e.CabRequests[f] == 1 {
			return true
		}
		if wv.HallRequests[elevio.BT_HallUp][f] > uint8(e.ElevNum) {
			return true
		}
		if wv.HallRequests[elevio.BT_HallDown][f] > uint8(e.ElevNum) {
			return true
		}
	}
	return false
}

func RequestsHere(e elevator.Elevator, wv elevator.Worldview) bool { //kan bytte ut alle e.Floor med elevio.getFloor!!!!
	if e.CabRequests[e.Last_Floor] == 1 {
		return true
	}
	if wv.HallRequests[elevio.BT_HallUp][e.Last_Floor] == uint8(e.ElevNum) {
		return true
	}
	if wv.HallRequests[elevio.BT_HallDown][e.Last_Floor] == uint8(e.ElevNum) {
		return true
	}
	return false
}

func RequestsHereCabOrUp(e elevator.Elevator, wv elevator.Worldview) bool { // stygt, kan ores på en linje
	if wv.HallRequests[elevio.BT_HallUp][e.Last_Floor] == uint8(e.ElevNum) {
		return true
	}
	if e.CabRequests[e.Last_Floor] == 1 {
		return true
	}
	return false
}

func RequestsHereCabOrDown(e elevator.Elevator, wv elevator.Worldview) bool {
	if wv.HallRequests[elevio.BT_HallDown][e.Last_Floor] == uint8(e.ElevNum) {
		return true
	}
	if e.CabRequests[e.Last_Floor] == 1 {
		return true
	}
	return false
}

func DeleteOrdersHere(e *elevator.Elevator, wv *elevator.Worldview) {
	for orderType := 0; orderType < 2; orderType++ {
		wv.HallRequests[orderType][e.Last_Floor] = 0
		elevio.SetButtonLamp(elevio.ButtonType(orderType), e.Last_Floor, false) //det med lys må fikses
	}
	e.CabRequests[e.Last_Floor] = 0
	wv.Version++
}

func SetOrderHere(e *elevator.Elevator, wv *elevator.Worldview, buttn elevio.ButtonEvent) {
	if buttn.Floor == e.Last_Floor {
		return
	}
	if buttn.Button == elevio.BT_Cab {
		e.CabRequests[e.Last_Floor] = 1
	} else {
		wv.HallRequests[buttn.Button][e.Last_Floor] = uint8(e.ElevNum)
	}
	elevio.SetButtonLamp(buttn.Button, buttn.Floor, true) // må endres når flere heiser
	wv.Version++
}

func ArrivedAtFloor(e *elevator.Elevator, wv *elevator.Worldview, timer_chan chan bool) {
	elevio.SetMotorDirection(elevio.MD_Stop)
	e.Dirn = elevio.MD_Stop
	DeleteOrdersHere(e, wv)
	wv.Version++
	go timer.TimerStart(3, timer_chan)
}

func DisplayQueueCont(e *elevator.Elevator) {
	e.Display()
	time.Sleep(time.Second)
}
