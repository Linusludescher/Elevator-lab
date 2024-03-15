package versioncontrol

import (
	"fmt"
	"project/elevator"
)

func Version_if_equal_queue(e elevator.Elevator, my_worldView elevator.Worldview, incoming_worldView elevator.Worldview) bool {
	areEqual := true
	for i := range my_worldView.HallRequests {
		for j := range my_worldView.HallRequests[i] {
			if my_worldView.HallRequests[i][j] != incoming_worldView.HallRequests[i][j] {
				areEqual = false
				break
			}
		}
		if !areEqual {
			break
		}
	}
	for elev := range my_worldView.ElevList {
		for f := range my_worldView.ElevList[elev].CabRequests {
			if elev == e.ElevNum-1 {
				if e.CabRequests[f] != incoming_worldView.ElevList[elev].CabRequests[f] {
					areEqual = false
					break
				}
			} else if my_worldView.ElevList[elev].CabRequests[f] != incoming_worldView.ElevList[elev].CabRequests[f] {
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

func Version_update_queue(e_p *elevator.Elevator, my_worldView_p *elevator.Worldview, incoming_worldView elevator.Worldview) {
	if incoming_worldView.Version > my_worldView_p.Version || ((my_worldView_p.Version > elevator.V_l-elevator.V_s_c) && incoming_worldView.Version < elevator.V_s_c) {
		my_worldView_p.HallRequests = incoming_worldView.HallRequests
		my_worldView_p.Version = incoming_worldView.Version
		my_worldView_p.ElevList = incoming_worldView.ElevList
		e_p.CabRequests = incoming_worldView.ElevList[e_p.ElevNum-1].CabRequests //La til dette
		// Slå av og på lys
		go elevator.UpdateLights(*my_worldView_p, e_p.ElevNum)
	} else if incoming_worldView.Version == my_worldView_p.Version {
		go elevator.UpdateLights(*my_worldView_p, e_p.ElevNum)
	} else if (incoming_worldView.Version == my_worldView_p.Version) && !Version_if_equal_queue(*e_p, *my_worldView_p, incoming_worldView) {
		fmt.Println("YOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO")
		if incoming_worldView.Sender > my_worldView_p.Sender {
			my_worldView_p.HallRequests = incoming_worldView.HallRequests
			my_worldView_p.Version = incoming_worldView.Version
			my_worldView_p.ElevList = incoming_worldView.ElevList
			e_p.CabRequests = incoming_worldView.ElevList[e_p.ElevNum-1].CabRequests //La til dette
			my_worldView_p.Version_up()
		}
	}
} // må ha noe med når version nullstilles
