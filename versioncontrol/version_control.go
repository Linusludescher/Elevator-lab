package versioncontrol

import (
	"fmt"
	w "project/worldview"
)

func versionIfEqualQueue(elev w.Elevator, my_worldView w.Worldview, incoming_worldView w.Worldview) bool {
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
	for elevator := range my_worldView.ElevList {
		for f := range my_worldView.ElevList[elevator].CabRequests {
			if elevator == elev.ElevNum-1 {
				if elev.CabRequests[f] != incoming_worldView.ElevList[elevator].CabRequests[f] {
					areEqual = false
					break
				}
			} else if my_worldView.ElevList[elevator].CabRequests[f] != incoming_worldView.ElevList[elevator].CabRequests[f] {
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

func CheckIncomingWorldView(readChannels w.ReadWorldviewChannels,
	version_up_chan chan<- bool,
	incomingWorldView w.Worldview,
	update_to_incoming_worldview_chan chan<- w.Worldview,
	update_lights_chan chan<- int) {

	myWorldView := w.ReadWorldView(readChannels)
	myElev := w.ReadElevator(readChannels)
	if incomingWorldView.Version > myWorldView.Version || ((myWorldView.Version > w.VERSIONLIMIT-w.VERSIONBUFFER) && incomingWorldView.Version < w.VERSIONBUFFER) {
		update_to_incoming_worldview_chan <- incomingWorldView

		update_lights_chan <- myElev.ElevNum
	} else if incomingWorldView.Version == myWorldView.Version {
		update_lights_chan <- myElev.ElevNum
	} else if (incomingWorldView.Version == myWorldView.Version) && !versionIfEqualQueue(myElev, myWorldView, incomingWorldView) {
		fmt.Println("YOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO")
		if incomingWorldView.Sender > myWorldView.Sender {
			update_to_incoming_worldview_chan <- incomingWorldView
			version_up_chan <- true
		}
	}
}
