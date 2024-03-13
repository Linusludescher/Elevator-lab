package main

import (
	"flag"
	"fmt"
	"project/elevator"
	"project/elevio"
	"project/network"
	"project/network/bcast"
	"project/stm"
	"project/versioncontrol"
)

//Todo rydding: samle ting i funkdjonrt
//og endre objektnavn på elevator- heter nå e eller my_elevator eller elevator

func main() {
	idFlag := flag.Int("id", 1, "Specifies an ID number")

	// Parse the command-line flags
	flag.Parse()

	// Retrieve the value of the idFlag
	id := *idFlag

	numFloors := 4 //endre dette??? fjerne??

	elevio.Init("localhost:15657", numFloors)
	//elevio.Init("localhost:22222", numFloors)

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	timer_exp_chan := make(chan bool)
	bc_timer_chan := make(chan bool)

	my_elevator, my_wv := elevator.ElevatorInit(id)

	network_channels := network.Init_network(id, &my_elevator, &my_wv)

	go network.PeersOnline(&my_elevator, &my_wv, network_channels)
	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)
	go elevator.BroadcastElevator(bc_timer_chan, 10)

	for {
		select {
		case <-timer_exp_chan:
			fmt.Println("timer expired")
			stm.TimerExp(&my_elevator, my_wv)

		case buttn := <-drv_buttons:
			stm.ButtonPressed(&my_elevator, &my_wv, buttn)

		case floor_sens := <-drv_floors:
			stm.FloorSensed(&my_elevator, &my_wv, floor_sens, timer_exp_chan, drv_obstr)

		case obstr := <-drv_obstr:
			stm.Obstruction(&my_elevator, &my_wv, obstr)

		case <-drv_stop:
			stm.StopButtonPressed(my_elevator)

		case udp_packet := <-network_channels.PacketRx: //legge til
			versioncontrol.Version_update_queue(&my_elevator, &my_wv, udp_packet)

		case <-bc_timer_chan:
			bcast.BcWorldView(my_elevator, &my_wv, network_channels.PacketTx)

			stm.DefaultState(&my_elevator, &my_wv, network_channels.PacketTx)
			//default:
			//my_wv.Display()
		}
	}
}
