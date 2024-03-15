package requests

import (
	"project/elevator"
	"project/elevio"
)

func RequestsAbove(elev_chan <-chan elevator.Elevator, worldView_chan <-chan elevator.Worldview, read_elev_chan chan<- bool, read_worldView_chan chan<- bool) bool {
	read_elev_chan <- true
	read_worldView_chan <- true
	elev := <-elev_chan
	worldView := <-worldView_chan

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

func RequestsBelow(elev_chan <-chan elevator.Elevator, worldView_chan <-chan elevator.Worldview, read_elev_chan chan<- bool, read_worldView_chan chan<- bool) bool {
	read_elev_chan <- true
	read_worldView_chan <- true
	elev := <-elev_chan
	worldView := <-worldView_chan
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

func RequestsHere(elev_chan <-chan elevator.Elevator, worldView_chan <-chan elevator.Worldview, read_elev_chan chan<- bool, read_worldView_chan chan<- bool) bool {
	read_elev_chan <- true
	read_worldView_chan <- true
	elev := <-elev_chan
	worldView := <-worldView_chan
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

func RequestsHereCabOrUp(elev_chan <-chan elevator.Elevator, worldView_chan <-chan elevator.Worldview, read_elev_chan chan<- bool, read_worldView_chan chan<- bool) bool { // TODO: stygt, kan ores på en linje
	read_elev_chan <- true
	read_worldView_chan <- true
	elev := <-elev_chan
	worldView := <-worldView_chan
	if worldView.HallRequests[elev.Last_Floor][elevio.BT_HALLUP] == uint8(elev.ElevNum) {
		return true
	}
	if elev.CabRequests[elev.Last_Floor] {
		return true
	}
	return false
}

func RequestsHereCabOrDown(elev_chan <-chan elevator.Elevator, worldView_chan <-chan elevator.Worldview, read_elev_chan chan<- bool, read_worldView_chan chan<- bool) bool {
	read_elev_chan <- true
	read_worldView_chan <- true
	elev := <-elev_chan
	worldView := <-worldView_chan
	if worldView.HallRequests[elev.Last_Floor][elevio.BT_HALLDOWN] == uint8(elev.ElevNum) {
		return true
	}
	if elev.CabRequests[elev.Last_Floor] {
		return true
	}
	return false
}
