package stm

import (
	"net"
	"project/elevator"
	"project/elevio"
	"project/network"
	"project/network/bcast"
	"project/requests"
	"project/versioncontrol"
	"time"
)

func MainFSM(
	timer_exp_chan <-chan bool,
	watchdog_chan chan bool,
	processPairConn *net.UDPConn,
	drv_buttons_chan chan elevio.ButtonEvent,
	reset_timer_chan chan bool,
	drv_floors_chan chan int,
	network_channels network.NetworkChan,
	bc_timer_chan chan bool,
	update_lights_chan chan int,
	readChannels elevator.ReadWorldviewChannels,
	updateChannels elevator.UpdateWorldviewChannels,
	ioChannels elevio.ElevioChannels) {
	for {
		select {
		case <-timer_exp_chan:
			ClosingDoor(readChannels, watchdog_chan, updateChannels.Update_direction_chan, ioChannels.Set_door_open_lamp_chan)

		case buttn := <-drv_buttons_chan:
			ButtonPressed(updateChannels.Set_order_chan, buttn)

		case floor_sens := <-drv_floors_chan:
			updateChannels.Update_floor_chan <- floor_sens
			FloorSensed(updateChannels.Update_direction_chan, updateChannels.Arrived_at_floor_chan, readChannels, floor_sens, watchdog_chan)

		case incomingWorldview := <-network_channels.PacketRx_chan: //legge til
			versioncontrol.CheckIncomingWorldView(readChannels, updateChannels.Version_up_chan, incomingWorldview, updateChannels.Update_to_incoming_chan, update_lights_chan)

		case <-bc_timer_chan:
			DefaultState(update_lights_chan, updateChannels.Update_direction_chan, updateChannels.Arrived_at_floor_chan, readChannels)
			bcast.BcWorldView(readChannels, network_channels.PacketTx_chan)
			processPairConn.Write([]byte("42"))
			my_worldView := elevator.ReadWorldView(readChannels)
			my_worldView.Display()

		case elevnum := <-update_lights_chan:
			UpdateLights(ioChannels.Set_button_lamp_chan, readChannels, elevnum)
		}
	}
}

func ClosingDoor(readChannels elevator.ReadWorldviewChannels,
	watchdog_chan chan<- bool,
	update_direction_chan chan<- elevio.MotorDirection,
	set_door_open_lamp_chan chan<- bool) {

	watchdog_chan <- false
	set_door_open_lamp_chan <- false
	worldView := elevator.ReadWorldView(readChannels)
	elev := elevator.ReadElevator(readChannels)
	if elev.Last_dir == elevio.MD_UP {
		if requests.RequestsAbove(elev, worldView) {
			update_direction_chan <- elevio.MD_UP
		} else if requests.RequestsBelow(elev, worldView) {
			update_direction_chan <- elevio.MD_DOWN
		} else {
			update_direction_chan <- elevio.MD_STOP
		}
	} else if elev.Last_dir == elevio.MD_DOWN {
		if requests.RequestsBelow(elev, worldView) {
			update_direction_chan <- elevio.MD_DOWN
		} else if requests.RequestsAbove(elev, worldView) {
			update_direction_chan <- elevio.MD_UP
		} else {
			update_direction_chan <- elevio.MD_STOP
		}
	} else {
		update_direction_chan <- elevio.MD_STOP
	}
}

func ButtonPressed(set_order_chan chan<- elevio.ButtonEvent, buttn elevio.ButtonEvent) {
	set_order_chan <- buttn
}

func FloorSensed(update_direction_chan chan<- elevio.MotorDirection, arrived_at_floor_chan chan<- bool, readChannels elevator.ReadWorldviewChannels, floor_sens int, watchdog_chan chan<- bool) {
	watchdog_chan <- false
	worldView := elevator.ReadWorldView(readChannels)
	elev := elevator.ReadElevator(readChannels)

	if elev.Dirn == elevio.MD_UP && floor_sens != -1 {
		if requests.RequestsHereCabOrUp(elev, worldView) {
			arrived_at_floor_chan <- true
		} else if (!requests.RequestsAbove(elev, worldView)) && requests.RequestsHere(elev, worldView) {
			arrived_at_floor_chan <- true
		}
	}
	if elev.Dirn == elevio.MD_DOWN && floor_sens != -1 {
		if requests.RequestsHereCabOrDown(elev, worldView) {
			arrived_at_floor_chan <- true
		} else if (!requests.RequestsBelow(elev, worldView)) && requests.RequestsHere(elev, worldView) {
			arrived_at_floor_chan <- true
		}
	}
	//softstop: TODO kan ikke bruke hardkoda values!
	if (floor_sens == -1 && elev.Last_dir == elevio.MD_DOWN && elev.Last_Floor == 0) ||
		(floor_sens == -1 && elev.Last_dir == elevio.MD_UP && elev.Last_Floor == 3) {
		update_direction_chan <- elevio.MD_STOP
	}
} // TODO: her kan noe gjenbrukes, og kanskje deles opp?

func Obstruction(updateChannels elevator.UpdateWorldviewChannels, readChannels elevator.ReadWorldviewChannels, reset_timer_chan chan<- bool, obstruction_chan chan bool) {
	last_obst := false
	for {
		select {
		case obst := <-obstruction_chan:
			last_obst = obst
		default:
			readChannels.Read_request_elev_chan <- true
			elev := <-readChannels.Read_elev_chan
			if elev.Behaviour != elevator.EB_DOOR_OPEN {
				break
			}
			updateChannels.Update_obstr_chan <- last_obst
			updateChannels.Version_up_chan <- true
			if last_obst {
				reset_timer_chan <- true
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func DefaultState(update_lights_chan chan int, update_direction_chan chan<- elevio.MotorDirection, arrived_at_floor_chan chan<- bool, readChannels elevator.ReadWorldviewChannels) {
	worldView := elevator.ReadWorldView(readChannels)
	elev := elevator.ReadElevator(readChannels)

	go aloneUpdateLights(worldView, elev, update_lights_chan) // TODO: Noe rart med denne eller?

	for floor := range worldView.HallRequests {
		for _, order := range worldView.HallRequests[floor] {
			if order == uint8(elev.ElevNum) && floor == elev.Last_Floor && elevio.GetFloor() == floor {
				arrived_at_floor_chan <- true
			}
		}
	}
	if (elev.Dirn == elevio.MD_STOP) && (elev.Behaviour != elevator.EB_DOOR_OPEN) {
		if requests.RequestsAbove(elev, worldView) {
			update_direction_chan <- elevio.MD_UP
		} else if requests.RequestsBelow(elev, worldView) {
			update_direction_chan <- elevio.MD_DOWN
		}
	}
}

func UpdateLights(set_button_lamp_chan chan<- elevio.ButtonLampOrder, readChannels elevator.ReadWorldviewChannels, elevnum int) {
	worldView := elevator.ReadWorldView(readChannels)
	for floor, f := range worldView.HallRequests {
		for buttonType, order := range f {
			set_button_lamp_chan <- elevio.ButtonLampOrder{Button_type: elevio.ButtonType(buttonType), OrderFloor: floor, Value: order != 0}
		}
	}
	for i, elev := range worldView.ElevList {
		if i+1 != elevnum {
			continue
		}
		for floor, f := range elev.CabRequests {
			set_button_lamp_chan <- elevio.ButtonLampOrder{Button_type: elevio.BT_CAB, OrderFloor: floor, Value: f}
		}
	}
}

func aloneUpdateLights(worldView elevator.Worldview, elev elevator.Elevator, update_lights_chan chan<- int) { //TODO: dårlig funksjon: burde deles opp i checkAlone og noe annet
	only_elev := true
	for i := range worldView.ElevList {
		if worldView.ElevList[i].Online && worldView.ElevList[i].ElevNum != elev.ElevNum {
			only_elev = false
		}
	}
	if only_elev {
		update_lights_chan <- elev.ElevNum
	}
}
