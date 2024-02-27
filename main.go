package main

import (
	"fmt"
	"project/elevator"
	"project/elevio"
	"project/network"
	"project/requests"
	"project/timer"
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
	elevator_chan := make(chan elevator.Elevator) //kanskje en buffer her?
	udp_receive := make(chan network.Packet)

	my_elevator := elevator.Elevator_uninitialized()

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	// Network init
	ports := network.GetNetworkConfig()
	conn := network.UdpInitDail(ports.UDPBrodcast)
	defer conn.Close()
	go network.Broadcasteverysecond(elevator_chan, conn)
	for _, port := range ports.UDPReceve {
		go network.Listener(udp_receive, port)
	}

	for {
		select {
		case <-timer_chan:
			fmt.Println("timer expired")
			if my_elevator.Last_dir == elevio.MD_Up {
				if requests.RequestsAbove(my_elevator) {
					elevio.SetMotorDirection(elevio.MD_Up)
					my_elevator.Dirn = elevio.MD_Up
					my_elevator.Last_dir = elevio.MD_Up
				} else if requests.RequestsBelow(my_elevator) {
					elevio.SetMotorDirection(elevio.MD_Down)
					my_elevator.Dirn = elevio.MD_Down
					my_elevator.Last_dir = elevio.MD_Down
				}
			} else if my_elevator.Last_dir == elevio.MD_Down {
				if requests.RequestsBelow(my_elevator) {
					elevio.SetMotorDirection(elevio.MD_Down)
					my_elevator.Dirn = elevio.MD_Down
					my_elevator.Last_dir = elevio.MD_Down
				} else if requests.RequestsAbove(my_elevator) {
					elevio.SetMotorDirection(elevio.MD_Up)
					my_elevator.Dirn = elevio.MD_Up
					my_elevator.Last_dir = elevio.MD_Up
				}
			} else {
				my_elevator.Dirn = elevio.MD_Stop
			}

		case buttn := <-drv_buttons:

			requests.SetOrderHere(&my_elevator, buttn) // tuple her etterhvert
			my_elevator.Display()

			//sjekke om timer er ferdig

			if my_elevator.Dirn == elevio.MD_Stop {
				fmt.Printf("test1\n")
				if requests.RequestsAbove(my_elevator) {
					fmt.Printf("test2\n")
					elevio.SetMotorDirection(elevio.MD_Up)
					my_elevator.Last_dir = elevio.MD_Up
					my_elevator.Dirn = elevio.MD_Up
				} else if requests.RequestsBelow(my_elevator) {
					fmt.Printf("test3\n")
					elevio.SetMotorDirection(elevio.MD_Down)
					my_elevator.Last_dir = elevio.MD_Down
					my_elevator.Dirn = elevio.MD_Down
				}
			}

		case floor_sens := <-drv_floors:
			if floor_sens != -1 {
				my_elevator.Last_Floor = floor_sens
				fmt.Println("new floow")
			}
			if my_elevator.Dirn == elevio.MD_Up && floor_sens != -1 {
				fmt.Println("Test4")
				if requests.RequestsHereCabOrUp(my_elevator) {
					fmt.Println("stopping")
					elevio.SetMotorDirection(elevio.MD_Stop)
					my_elevator.Dirn = elevio.MD_Stop
					requests.DeleteOrdersHere(&my_elevator)
					go timer.TimerStart(3, timer_chan)
				} else if (!requests.RequestsAbove(my_elevator)) && requests.RequestsHere(my_elevator) { //samme kode som over
					elevio.SetMotorDirection(elevio.MD_Stop)
					my_elevator.Dirn = elevio.MD_Stop
					requests.DeleteOrdersHere(&my_elevator)
					go timer.TimerStart(3, timer_chan)
				}
			}
			if my_elevator.Dirn == elevio.MD_Down && floor_sens != -1 {
				fmt.Println("Test5")
				if requests.RequestsHereCabOrDown(my_elevator) {
					fmt.Println("stopping")
					elevio.SetMotorDirection(elevio.MD_Stop)
					my_elevator.Dirn = elevio.MD_Stop
					requests.DeleteOrdersHere(&my_elevator)
					go timer.TimerStart(3, timer_chan)
				} else if (!requests.RequestsBelow(my_elevator)) && requests.RequestsHere(my_elevator) { //samme kode som over
					elevio.SetMotorDirection(elevio.MD_Stop)
					my_elevator.Dirn = elevio.MD_Stop
					requests.DeleteOrdersHere(&my_elevator)
					go timer.TimerStart(3, timer_chan)
				}
			}

		case obstr := <-drv_obstr:
			if obstr {
				elevio.SetMotorDirection(elevio.MD_Stop)
			} else {
				elevio.SetMotorDirection(my_elevator.Dirn)
			}
		case <-drv_stop:
			// fjerne hele køen?
			requests.DeleteAllOrdes(&my_elevator)
			// vente ellerno?
		case udp_packet := <-udp_receive:
			my_elevator.Requests = udp_packet.Queue
		default:
			elevator_chan <- my_elevator

			if my_elevator.Dirn == elevio.MD_Stop {
				fmt.Printf("test1\n")
				if requests.RequestsAbove(my_elevator) {
					fmt.Printf("test2\n")
					elevio.SetMotorDirection(elevio.MD_Up)
					my_elevator.Last_dir = elevio.MD_Up
					my_elevator.Dirn = elevio.MD_Up
				} else if requests.RequestsBelow(my_elevator) {
					fmt.Printf("test3\n")
					elevio.SetMotorDirection(elevio.MD_Down)
					my_elevator.Last_dir = elevio.MD_Down
					my_elevator.Dirn = elevio.MD_Down
				}
			}
		}
	}
}
