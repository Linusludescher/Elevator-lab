package main

import (
	"flag"
	"fmt"
	"project/elevator"
	"project/elevio"
	"project/network"
	"project/network/bcast"
	"project/stm"
	"project/timer"
	// "strconv"
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

	elevatorConf := elevator.ReadElevatorConfig() //Dette burde bli en initfunksjon, og legges tilbake i elevator package!
	numFloors := int(elevatorConf.N_FLOORS)

	processPairConn := bcast.ProcessPairListner(id)

	//elevio.Init("localhost:"+localhostnr, numFloors)
	elevio.Init("localhost:15657", numFloors) //15657

	//clean_main greier:
	elevioChannels := elevio.InitElevioChannels()
	updateWorldviewChannels := elevator.InitUpdateWorldviewChannels()
	readChannels := elevator.InitReadWorldViewChannels()

	drv_buttons_chan := make(chan elevio.ButtonEvent)
	drv_floors_chan := make(chan int)
	drv_obstr_chan := make(chan bool)
	timer_exp_chan := make(chan bool)
	bc_timer_chan := make(chan bool)
	watchdog_chan := make(chan bool)
	reset_timer_chan := make(chan bool)
	update_lights_chan := make(chan int)

	my_elevator, my_wv := elevator.ElevatorInit(id)
	network_channels := network.InitNetwork(id)

	go elevator.UpdateWorldview(&my_wv, &my_elevator, reset_timer_chan, watchdog_chan, readChannels, updateWorldviewChannels, elevioChannels)
	go elevator.BroadcastElevator(bc_timer_chan, 10)

	go elevio.Elevio_select(elevioChannels)
	go elevio.PollButtons(drv_buttons_chan)
	go elevio.PollFloorSensor(drv_floors_chan)
	go elevio.PollObstructionSwitch(drv_obstr_chan)

	go timer.OperativeWatchdog(10, watchdog_chan)
	go timer.TimerStart(3, timer_exp_chan, reset_timer_chan)

	go stm.Obstruction(updateWorldviewChannels, readChannels, reset_timer_chan, drv_obstr_chan)

	go network.PeersOnline(readChannels, network_channels, updateWorldviewChannels)

	stm.MainFSM(timer_exp_chan, watchdog_chan, processPairConn, drv_buttons_chan, reset_timer_chan,
		drv_floors_chan, network_channels, bc_timer_chan, update_lights_chan, readChannels, updateWorldviewChannels, elevioChannels)
}
