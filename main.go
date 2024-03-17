package main

import (
	"flag"
	"fmt"
	"project/elevio"
	"project/fsm"
	"project/network"
	"project/network/bcast"
	"project/timer"
	w "project/worldview"
)

//Todo rydding: samle ting i funkdjonrt
//og endre objektnavn på elevator- heter nå e eller my_elevator eller elevator

func main() {

	idFlag := flag.Int("id", 1, "Specifies an ID number")

	// Parse the command-line flags
	flag.Parse()

	// Retrieve the value of the idFlag
	id := *idFlag
	//localhostnr := strconv.Itoa(19657 + id)
	fmt.Println(id)

	elevatorConf := w.ReadElevatorConfig() //Dette burde bli en initfunksjon, og legges tilbake i elevator package!
	numFloors := int(elevatorConf.N_FLOORS)

	processPairConn := bcast.ProcessPairListner(id)

	//elevio.Init("localhost:"+localhostnr, numFloors)
	elevio.Init("localhost:15657", numFloors) //15657

	//clean_main greier:
	elevioChannels := elevio.InitElevioChannels()
	updateWorldviewChannels := w.InitUpdateWorldviewChannels()
	readChannels := w.InitReadWorldViewChannels()

	drv_buttons_chan := make(chan elevio.ButtonEvent, 100)
	drv_floors_chan := make(chan int, 100)
	drv_obstr_chan := make(chan bool, 100)
	timer_exp_chan := make(chan bool, 100)
	bc_timer_chan := make(chan bool, 100)
	watchdog_chan := make(chan bool, 100)
	reset_timer_chan := make(chan bool, 100)
	update_lights_chan := make(chan int, 100)

	my_elevator, my_wv := w.ElevatorInit(timer_exp_chan, id)
	network_channels := network.InitNetwork(id)

	go w.UpdateWorldview(&my_wv, &my_elevator, reset_timer_chan, watchdog_chan, readChannels, updateWorldviewChannels, elevioChannels)
	go w.BroadcastElevator(bc_timer_chan, 10)
	go elevio.ElevioUpdate(elevioChannels)
	go elevio.PollButtons(drv_buttons_chan)
	go elevio.PollFloorSensor(drv_floors_chan)
	go elevio.PollObstructionSwitch(drv_obstr_chan)
	go timer.OperativeWatchdog(10, watchdog_chan)
	go timer.TimerStart(3, timer_exp_chan, reset_timer_chan)
	go fsm.Obstruction(updateWorldviewChannels, readChannels, reset_timer_chan, drv_obstr_chan)
	go network.PeersOnline(readChannels, network_channels, updateWorldviewChannels)

	fsm.MainFSM(timer_exp_chan, watchdog_chan, processPairConn, drv_buttons_chan, reset_timer_chan,
		drv_floors_chan, network_channels, bc_timer_chan, update_lights_chan, readChannels, updateWorldviewChannels, elevioChannels)
}
