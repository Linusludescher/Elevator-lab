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
	numFloors := 4 //endre dette??? fjerne??

	//processPairConn := bcast.ProcessPairListner(id)

	// elevio.Init("localhost:"+localhostnr, numFloors)
	elevio.Init("localhost:15657", numFloors) //15657

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	timer_exp_chan := make(chan bool)
	bc_timer_chan := make(chan bool)
	wd_chan := make(chan bool)
	resetTimer_chan := make(chan bool)

	my_elevator, my_wv := elevator.ElevatorInit(id)

	network_channels := network.Init_network(id, &my_elevator, &my_wv)

	go network.PeersOnline(&my_elevator, &my_wv, network_channels)
	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)
	go elevator.BroadcastElevator(bc_timer_chan, 10)
	go timer.OperativeWatchdog(&my_elevator, &my_wv, 10, wd_chan)
	go timer.TimerStart(&my_elevator, &my_wv, 3, timer_exp_chan, drv_obstr, resetTimer_chan)

	for {
		select {
		case <-timer_exp_chan:
			stm.ClosingDoor(&my_elevator, my_wv, wd_chan)

		case buttn := <-drv_buttons:
			stm.ButtonPressed(&my_elevator, &my_wv, buttn, resetTimer_chan, wd_chan)

		case floor_sens := <-drv_floors:
			stm.FloorSensed(&my_elevator, &my_wv, floor_sens, resetTimer_chan, wd_chan)

		case obstr := <-drv_obstr:
			stm.Obstruction(&my_elevator, &my_wv, obstr)

		case <-drv_stop:
			stm.StopButtonPressed(my_elevator)

		case udp_packet := <-network_channels.PacketRx: //legge til
			versioncontrol.Version_update_queue(&my_elevator, &my_wv, udp_packet)

		case <-bc_timer_chan:
			stm.DefaultState(&my_elevator, &my_wv, resetTimer_chan, wd_chan)
			bcast.BcWorldView(my_elevator, my_wv, network_channels.PacketTx)
			// processPairConn.Write([]byte("42"))
			//default:
			my_wv.Display()
		}
	}
}
