package main

import (
	"fmt"
	"project/elevator"
	"project/elevio"
	"project/network"
	"project/stm"
	"project/versioncontrol"
)

//Todo rydding: samle ting i funkdjonrt
//og endre objektnavn på elevator- heter nå e eller my_elevator eller elevator

func main() {

	numFloors := 4

	elevio.Init("localhost:15657", numFloors)
	//elevio.Init("localhost:22222", numFloors)

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	timer_exp_chan := make(chan bool)
	bc_timer_chan := make(chan bool)
	// broadcast_elevator_chan := make(chan elevator.Elevator) //kanskje en buffer her?
	// udp_receive_chan := make(chan network.Packet)           //kanskje en buffer her og?

	my_elevator, my_wv := elevator.ElevatorInit()

	network_channels := network.Init_network(&my_elevator, &my_wv)

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
			stm.FloorSensed(&my_elevator, &my_wv, floor_sens, timer_exp_chan)

		case obstr := <-drv_obstr:
			stm.Obstruction(my_elevator, obstr)

		case <-drv_stop:
			stm.StopButtonPressed(my_elevator)

		case udp_packet := <-network_channels.PacketRx: //legge til
			fmt.Println("Pakke mottatt")
			versioncontrol.Version_update_queue(&my_wv, udp_packet)
		case <-bc_timer_chan:
			network_channels.PacketTx <- my_wv
			stm.DefaultState(&my_elevator, &my_wv, network_channels.PacketTx)
			//default:
			//stm.DefaultState(&my_elevator, network_channels.PacketTx) // Dårlig navn? beskriver dårlig
		}
	}
}
