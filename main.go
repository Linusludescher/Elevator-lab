package main

import (
	"flag"
	"fmt"
	"project/elevio"
	"project/fsm"
	"project/network"
	"project/timer"
	w "project/worldview"
)

func main() {

	idFlag := flag.Int("id", 1, "Specifies an ID number")
	flag.Parse()
	id := *idFlag
	fmt.Printf("elevator number: %d:\n", id)
	elevatorConf := w.ReadElevatorConfig()
	numFloors := int(elevatorConf.N_FLOORS)

	//id, numFloors := w.GetElevatorCredentials()
	elevio.Init("localhost:15657", numFloors)

	elevioChannels := elevio.InitElevioChannels()
	updateWorldviewChannels := w.InitUpdateWorldviewChannels()
	readChannels := w.InitReadWorldViewChannels()
	network_channels := network.InitNetwork(id)

	drv_buttons_chan := make(chan elevio.ButtonEvent, 100)
	drv_floors_chan := make(chan int, 100)
	drv_obstr_chan := make(chan bool, 100)
	timer_exp_chan := make(chan bool, 100)
	bc_timer_chan := make(chan bool, 100)
	watchdog_chan := make(chan bool, 100)
	reset_timer_chan := make(chan bool, 100)
	update_lights_chan := make(chan int, 100)

	my_elevator, my_wv := w.WorldviewInit(timer_exp_chan, id)

	go w.UpdateWorldview(&my_wv, &my_elevator, reset_timer_chan, watchdog_chan, readChannels, updateWorldviewChannels, elevioChannels)
	go w.StartBroadcastLoop(bc_timer_chan, 10)
	go elevio.ElevioUpdate(elevioChannels)
	go elevio.PollButtons(drv_buttons_chan)
	go elevio.PollFloorSensor(drv_floors_chan)
	go elevio.PollObstructionSwitch(drv_obstr_chan)
	go timer.OperativeWatchdog(10, watchdog_chan)
	go timer.TimerStart(3, timer_exp_chan, reset_timer_chan)
	go fsm.Obstruction(updateWorldviewChannels, readChannels, reset_timer_chan, drv_obstr_chan)
	go network.PeersOnline(readChannels, network_channels, updateWorldviewChannels)

	fsm.MainFSM(timer_exp_chan, watchdog_chan, drv_buttons_chan, reset_timer_chan,
		drv_floors_chan, network_channels, bc_timer_chan, update_lights_chan, readChannels, updateWorldviewChannels, elevioChannels)
}
