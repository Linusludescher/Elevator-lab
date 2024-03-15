package requests

import (
	"project/costFunc"
	"project/elevator"
	"project/elevio"
	"time"
)

func RequestsAbove(elev elevator.Elevator, worldView elevator.Worldview) bool {
	for f := elev.Last_Floor + 1; f < elevator.N_FLOORS; f++ {
		if elev.CabRequests[f] {
			return true
		}
		if worldView.HallRequests[f][elevio.BT_HALLUP] == uint8(elev.ElevNum) {
			return true
		}
		if worldView.HallRequests[f][elevio.BT_HALLDOWN] == uint8(elev.ElevNum) {
			return true
		}
	}
	return false
}

func RequestsBelow(elev elevator.Elevator, worldView elevator.Worldview) bool {
	for f := 0; f < elev.Last_Floor; f++ {
		if elev.CabRequests[f] {
			return true
		}
		if worldView.HallRequests[f][elevio.BT_HALLUP] == uint8(elev.ElevNum) {
			return true
		}
		if worldView.HallRequests[f][elevio.BT_HALLDOWN] == uint8(elev.ElevNum) {
			return true
		}
	}
	return false
}

func RequestsHere(elev elevator.Elevator, worldView elevator.Worldview) bool { //kan bytte ut alle e.Floor med elevio.getFloor!!!!
	if elev.CabRequests[elev.Last_Floor] {
		return true
	}
	if worldView.HallRequests[elev.Last_Floor][elevio.BT_HALLUP] == uint8(elev.ElevNum) {
		return true
	}
	if worldView.HallRequests[elev.Last_Floor][elevio.BT_HALLDOWN] == uint8(elev.ElevNum) {
		return true
	}
	return false
}

func RequestsHereCabOrUp(elev elevator.Elevator, worldView elevator.Worldview) bool { // stygt, kan ores på en linje
	if worldView.HallRequests[elev.Last_Floor][elevio.BT_HALLUP] == uint8(elev.ElevNum) {
		return true
	}
	if elev.CabRequests[elev.Last_Floor] {
		return true
	}
	return false
}

func RequestsHereCabOrDown(elev elevator.Elevator, worldView elevator.Worldview) bool {
	if worldView.HallRequests[elev.Last_Floor][elevio.BT_HALLDOWN] == uint8(elev.ElevNum) {
		return true
	}
	if elev.CabRequests[elev.Last_Floor] {
		return true
	}
	return false
}

func DeleteOrdersHere(elev_p *elevator.Elevator, worldView_p *elevator.Worldview) {
	for orderType := 0; orderType < 2; orderType++ {
		if worldView_p.HallRequests[elev_p.Last_Floor][orderType] == uint8(elev_p.ElevNum) {
			worldView_p.HallRequests[elev_p.Last_Floor][orderType] = 0
		}
	}
	elev_p.CabRequests[elev_p.Last_Floor] = false
	worldView_p.VersionUp()
}

func SetOrder(elev_p *elevator.Elevator, worldView_p *elevator.Worldview, buttn elevio.ButtonEvent, reset_timer_chan chan bool, wd_chan chan bool) {
	if (buttn.Floor == elev_p.Last_Floor) && (elevio.GetFloor() != -1) {
		ArrivedAtFloor(elev_p, worldView_p, reset_timer_chan, wd_chan)
		return
	}
	if buttn.Button == elevio.BT_CAB {
		elev_p.CabRequests[buttn.Floor] = true
	} else if worldView_p.HallRequests[buttn.Floor][buttn.Button] == 0 {
		costFunc.CostFunction(worldView_p, buttn)
	}
	worldView_p.VersionUp()
}

func ArrivedAtFloor(elev_p *elevator.Elevator, worldView_p *elevator.Worldview, reset_timer_chan chan<- bool, wd_chan chan<- bool) {
	elevio.SetDoorOpenLamp(true)
	elevio.SetMotorDirection(elevio.MD_STOP)
	elev_p.Dirn = elevio.MD_STOP
	DeleteOrdersHere(elev_p, worldView_p)
	worldView_p.VersionUp()
	elev_p.Behaviour = elevator.EB_DOOR_OPEN
	wd_chan <- true
	reset_timer_chan <- true
}

func DisplayQueueCont(elev_p *elevator.Elevator) {
	elev_p.Display()
	time.Sleep(time.Second)
}
