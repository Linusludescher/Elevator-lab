package versioncontrol
import (
	"project/elevator"
	"project/network"
)

func Version_up(e *elevator.Elevator) {
	e.Version ++
}

func Version_if_equal_queue(e elevator.Elevator, p network.Packet) (bool){
	areEqual := true
	for i := range e.Requests {
		for j := range e.Requests[i] {
			if e.Requests[i][j] != p.Queue[i][j] {
				areEqual = false
				break
			}
		}
		if !areEqual {
			break
		}
	}
	return areEqual
}

func Version_update_queue(e *elevator.Elevator, p network.Packet) {
	if p.Version > e.Version {
		e.Requests = p.Queue
		e.Version = p.Version
	} else if (p.Version == e.Version) || !Version_if_equal_queue(*e,p) {
		if p.ElevatorNum > e.ElevNum {
			e.Requests = p.Queue
			e.Version = p.Version
		}
	}
}			// må ha noe med når version nullstilles