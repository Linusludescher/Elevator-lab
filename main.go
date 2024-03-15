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
	"project/versioncontrol"
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
	// localhostnr := strconv.Itoa(19657 + id)
	fmt.Println(id)

	elevatorConf := elevator.ReadElevatorConfig() //Dette burde bli en initfunksjon, og legges tilbake i elevator package!
	numFloors := int(elevatorConf.N_FLOORS)

	//processPairConn := bcast.ProcessPairListner(id)

	// elevio.Init("localhost:"+localhostnr, numFloors)
	elevio.Init("localhost:15657", numFloors) //15657

	drv_buttons_chan := make(chan elevio.ButtonEvent)
	drv_floors_chan := make(chan int)
	drv_obstr_chan := make(chan bool)
	timer_exp_chan := make(chan bool)
	bc_timer_chan := make(chan bool)
	wd_chan := make(chan bool)
	reset_timer_chan := make(chan bool)

	my_elevator, my_wv := elevator.ElevatorInit(id)
	network_channels := network.InitNetwork(id)

	go network.PeersOnline(&my_wv, network_channels)
	go elevio.PollButtons(drv_buttons_chan)
	go elevio.PollFloorSensor(drv_floors_chan)
	go elevio.PollObstructionSwitch(drv_obstr_chan)
	go elevator.BroadcastElevator(bc_timer_chan, 10)
	go timer.OperativeWatchdog(10, wd_chan)
	go timer.TimerStart(3, timer_exp_chan, reset_timer_chan)
	go stm.Obstruction(&my_elevator, &my_wv, drv_obstr_chan, reset_timer_chan)

	for {
		select {
		case <-timer_exp_chan:
			stm.ClosingDoor(&my_elevator, my_wv, wd_chan)

		case buttn := <-drv_buttons_chan:
			stm.ButtonPressed(&my_elevator, &my_wv, buttn, reset_timer_chan, wd_chan)

		case floor_sens := <-drv_floors_chan:
			stm.FloorSensed(&my_elevator, &my_wv, floor_sens, reset_timer_chan, wd_chan)

		case udp_packet := <-network_channels.PacketRx_chan: //legge til
			versioncontrol.VersionUpdateQueue(&my_elevator, &my_wv, udp_packet)

		case <-bc_timer_chan:
			stm.DefaultState(&my_elevator, &my_wv, reset_timer_chan, wd_chan)
			bcast.BcWorldView(my_elevator, my_wv, network_channels.PacketTx_chan)
			// processPairConn.Write([]byte("42"))
			my_wv.Display()
		}
	}
}
