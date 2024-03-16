package requests

import (
	io "project/elevio"
	w "project/worldview"
)

func RequestsAbove(elev w.Elevator, worldview w.Worldview) bool {
	for f := elev.Last_Floor + 1; f < len(elev.CabRequests); f++ {
		if elev.CabRequests[f] {
			return true
		}
		if worldview.HallRequests[f][io.BT_HALLUP] == uint8(elev.ElevNum) {
			return true
		}
		if worldview.HallRequests[f][io.BT_HALLDOWN] == uint8(elev.ElevNum) {
			return true
		}
	}
	return false
}

func RequestsBelow(elev w.Elevator, worldview w.Worldview) bool {
	for f := 0; f < elev.Last_Floor; f++ {
		if elev.CabRequests[f] {
			return true
		}
		if worldview.HallRequests[f][io.BT_HALLUP] == uint8(elev.ElevNum) {
			return true
		}
		if worldview.HallRequests[f][io.BT_HALLDOWN] == uint8(elev.ElevNum) {
			return true
		}
	}
	return false
}

func RequestsHere(elev w.Elevator, worldView w.Worldview) bool {
	if elev.CabRequests[elev.Last_Floor] {
		return true
	}
	if worldView.HallRequests[elev.Last_Floor][io.BT_HALLUP] == uint8(elev.ElevNum) {
		return true
	}
	return worldView.HallRequests[elev.Last_Floor][io.BT_HALLDOWN] == uint8(elev.ElevNum)
}

func RequestsHereCabOrUp(elev w.Elevator, worldView w.Worldview) bool {
	if worldView.HallRequests[elev.Last_Floor][io.BT_HALLUP] == uint8(elev.ElevNum) {
		return true
	}
	return elev.CabRequests[elev.Last_Floor]
}

func RequestsHereCabOrDown(elev w.Elevator, worldView w.Worldview) bool {
	if worldView.HallRequests[elev.Last_Floor][io.BT_HALLDOWN] == uint8(elev.ElevNum) {
		return true
	}
	return elev.CabRequests[elev.Last_Floor]
}
