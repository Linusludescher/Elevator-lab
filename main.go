package main

import (
	"fmt"
	"project/elevator"
	"project/elevio"
	"project/network"
	"project/stm"
	versioncontrol "project/version_control"
)

//Todo rydding: samle ting i funkdjonrt
//og endre objektnavn på elevator- heter nå e eller my_elevator eller elevator

func main() {

	numFloors := 4

	elevio.Init("localhost:15657", numFloors)

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	timer_chan := make(chan bool)
	broadcast_elevator_chan := make(chan elevator.Elevator) //kanskje en buffer her?
	udp_receive_chan := make(chan network.Packet)           //kanskje en buffer her og?

	my_elevator := elevator.Elevator_uninitialized()

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	conn := network.NetworkInit(broadcast_elevator_chan, udp_receive_chan)
	defer conn.Close()

	for {
		select {
		case <-timer_chan:
			fmt.Println("timer expired")
			stm.TimerState(&my_elevator)

		case buttn := <-drv_buttons:
			stm.ButtonPressed(&my_elevator, buttn)

		case floor_sens := <-drv_floors:
			stm.FloorSensed(&my_elevator, floor_sens, timer_chan)

		case obstr := <-drv_obstr:
			stm.Obstuction(my_elevator, obstr)

		case <-drv_stop:
			stm.StopButtonPressed(my_elevator)

		case udp_packet := <-udp_receive_chan:
			versioncontrol.Version_update_queue(&my_elevator, udp_packet)

		default:
			stm.DefaultState(&my_elevator, broadcast_elevator_chan) // Dårlig navn? beskriver dårlig
		}
	}
}
