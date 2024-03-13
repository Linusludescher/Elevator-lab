package versioncontrol

import (
	"fmt"
	"project/elevator"
)

func Version_if_equal_queue(e elevator.Elevator, my_wv elevator.Worldview, incoming_wv elevator.Worldview) bool {
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
			if elev == e.ElevNum-1 {
				if e.CabRequests[f] != incoming_wv.ElevList[elev].CabRequests[f] {
					areEqual = false
					break
				}
			} else if my_wv.ElevList[elev].CabRequests[f] != incoming_wv.ElevList[elev].CabRequests[f] {
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

func Version_update_queue(e_p *elevator.Elevator, my_wv_p *elevator.Worldview, incoming_wv elevator.Worldview) {
	if incoming_wv.Version > my_wv_p.Version || ((my_wv_p.Version > elevator.V_l-elevator.V_s_c) && incoming_wv.Version < elevator.V_s_c) { //||  (e.Version > versionLimit-versionStabilityCycles && p.Version < versionStabilityCycles)
		my_wv_p.HallRequests = incoming_wv.HallRequests
		my_wv_p.Version = incoming_wv.Version
		my_wv_p.ElevList = incoming_wv.ElevList
		e_p.CabRequests = incoming_wv.ElevList[e_p.ElevNum-1].CabRequests //La til dette
		// Slå av og på lys
		elevator.SetHallLights(*my_wv_p, e_p.ElevNum)
	} else if incoming_wv.Version == my_wv_p.Version {
		elevator.SetHallLights(*my_wv_p, e_p.ElevNum)
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
