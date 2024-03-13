package versioncontrol

import (
	"fmt"
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

func Version_if_equal_queue(e elevator.Elevator, my_wv elevator.Worldview, incoming_wv elevator.Worldview) bool {
	areEqual := true
	for i := range my_wv.HallRequests {
		for j := range my_wv.HallRequests[i] {
			if my_wv.HallRequests[i][j] != incoming_wv.HallRequests[i][j] {
				areEqual = false
				fmt.Println("HER ER FEILEN!!")
				break
			}
		}
		if !areEqual {
			break
		}
	}
	for elev := range my_wv.ElevList {
		for f := range my_wv.ElevList[elev].CabRequests {
			if elev == e.ElevNum-1 {
				if e.CabRequests[f] != incoming_wv.ElevList[elev].CabRequests[f] {
					areEqual = false
					break
				}
			} else if my_wv.ElevList[elev].CabRequests[f] != incoming_wv.ElevList[elev].CabRequests[f] {
				areEqual = false
				fmt.Println("HER ER FEILEN22222222222222222222!!")
				fmt.Printf("elev %d\t f: %d\t\n", elev, f)
				fmt.Printf("my cab orders: %v\t incoming carb orders: %v\t \n", my_wv.ElevList[elev].CabRequests, incoming_wv.ElevList[elev].CabRequests)
				break
			}
		}
		if !areEqual {
			break
		}
	}
	return areEqual
}

func Version_update_queue(e_p *elevator.Elevator, my_wv_p *elevator.Worldview, incoming_wv elevator.Worldview) {
	if incoming_wv.Version > my_wv_p.Version { //||  (e.Version > versionLimit-versionStabilityCycles && p.Version < versionStabilityCycles)
		my_wv_p.HallRequests = incoming_wv.HallRequests
		my_wv_p.Display()
		my_wv_p.Version = incoming_wv.Version
		my_wv_p.ElevList = incoming_wv.ElevList
		e_p.CabRequests = incoming_wv.ElevList[e_p.ElevNum-1].CabRequests //La til dette
		my_wv_p.Display()

	} else if (incoming_wv.Version == my_wv_p.Version) && !Version_if_equal_queue(*e_p, *my_wv_p, incoming_wv) {
		fmt.Println("YOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO")
		if incoming_wv.Sender > my_wv_p.Sender {
			my_wv_p.HallRequests = incoming_wv.HallRequests
			my_wv_p.Version = incoming_wv.Version
			my_wv_p.ElevList = incoming_wv.ElevList
			e_p.CabRequests = incoming_wv.ElevList[e_p.ElevNum-1].CabRequests //La til dette
		}
	}
} // må ha noe med når version nullstilles
