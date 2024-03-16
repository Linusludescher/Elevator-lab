package requests

import (
	"project/elevator"
	"project/elevio"
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

func RequestsHere(elev elevator.Elevator, worldView elevator.Worldview) bool {
	if elev.CabRequests[elev.Last_Floor] {
		return true
	}
	if worldView.HallRequests[elev.Last_Floor][elevio.BT_HALLUP] == uint8(elev.ElevNum) {
		return true
	}
	return worldView.HallRequests[elev.Last_Floor][elevio.BT_HALLDOWN] == uint8(elev.ElevNum)
}

func RequestsHereCabOrUp(elev elevator.Elevator, worldView elevator.Worldview) bool {
	if worldView.HallRequests[elev.Last_Floor][elevio.BT_HALLUP] == uint8(elev.ElevNum) {
		return true
	}
	return elev.CabRequests[elev.Last_Floor]
}

func RequestsHereCabOrDown(elev elevator.Elevator, worldView elevator.Worldview) bool {
	if worldView.HallRequests[elev.Last_Floor][elevio.BT_HALLDOWN] == uint8(elev.ElevNum) {
		return true
	}
	return elev.CabRequests[elev.Last_Floor]
}
