package versioncontrol

import (
	"project/elevator"
)

const (
	versionLimit = 18446744073709551615 //2e64-1
	// versionStabilityCycles = 100 //maks antall sykler ny versjon kan være foran for at e.Version settes godtar lavere p.Version (ved Version overflow)
	// versionInitVal = 10000 //initialisere på høyere verdi enn 0 for ikke problemer med nullstilling ved tilbakekobling etter utfall
)

func Version_up(wv *elevator.Worldview) {
	if wv.Version < versionLimit {
		wv.Version++
	} else {
		wv.Version = 0
	}
}

func Version_if_equal_queue(my_wv elevator.Worldview, incoming_wv elevator.Worldview) bool {
	areEqual := true
	for i := range my_wv.HallRequests {
		for j := range my_wv.HallRequests[i] {
			if my_wv.HallRequests[i][j] != incoming_wv.HallRequests[i][j] {
				areEqual = false
				break
			}
		}
		if !areEqual {
			break
		}
	}
	for elev := range my_wv.ElevList {
		for f := range my_wv.ElevList[elev].CabRequests {
			if my_wv.ElevList[elev].CabRequests[f] != incoming_wv.ElevList[elev].CabRequests[f] {
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

func Version_update_queue(my_wv *elevator.Worldview, incoming_wv elevator.Worldview) {
	if incoming_wv.Version > my_wv.Version { //||  (e.Version > versionLimit-versionStabilityCycles && p.Version < versionStabilityCycles)
		my_wv.HallRequests = incoming_wv.HallRequests
		my_wv.Version = incoming_wv.Version
		my_wv.ElevList = incoming_wv.ElevList
	} else if (incoming_wv.Version == my_wv.Version) || !Version_if_equal_queue(*my_wv, incoming_wv) {
		if incoming_wv.Sender > my_wv.Sender {
			my_wv.HallRequests = incoming_wv.HallRequests
			my_wv.Version = incoming_wv.Version
			my_wv.ElevList = incoming_wv.ElevList
		}
	}
} // må ha noe med når version nullstilles
