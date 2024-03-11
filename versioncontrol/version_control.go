package versioncontrol

import (
	"project/elevator"
)

const (
	versionLimit = 18446744073709551615 //2e64-1
	// versionStabilityCycles = 100 //maks antall sykler ny versjon kan være foran for at e.Version settes godtar lavere p.Version (ved Version overflow)
	// versionInitVal = 10000 //initialisere på høyere verdi enn 0 for ikke problemer med nullstilling ved tilbakekobling etter utfall
)

func Version_up(e *elevator.Elevator) {
	if e.Version < versionLimit {
		e.Version++
	} else {
		e.Version = 0
	}
}

func Version_if_equal_queue(e elevator.Elevator, p elevator.Elevator) bool {
	areEqual := true
	for i := range e.Requests {
		for j := range e.Requests[i] {
			if e.Requests[i][j] != p.Requests[i][j] {
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

func Version_update_queue(e *elevator.Elevator, p elevator.Elevator) {
	if p.Version > e.Version { //||  (e.Version > versionLimit-versionStabilityCycles && p.Version < versionStabilityCycles)
		e.Requests = p.Requests
		e.Version = p.Version
	} else if (p.Version == e.Version) || !Version_if_equal_queue(*e, p) {
		if p.ElevNum > e.ElevNum {
			e.Requests = p.Requests
			e.Version = p.Version
		}
	}
} // må ha noe med når version nullstilles
