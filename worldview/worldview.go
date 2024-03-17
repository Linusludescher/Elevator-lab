package worldview

import (
	"fmt"
	"project/elevio"
)

const STARTVERSION uint64 = 5000

type Worldview struct {
	ElevList     []Elevator
	Sender       int
	Version      uint64
	HallRequests [][2]uint8
}

const (
	VERSIONLIMIT  = 18446744073709551615 //2e64-1
	VERSIONBUFFER = 100000               //maks antall sykler ny versjon kan være foran for at e.Version settes godtar lavere p.Version (ved Version overflow)
	// versionInitVal = 10000 //initialisere på høyere verdi enn 0 for ikke problemer med nullstilling ved tilbakekobling etter utfall
)

type ReadWorldviewChannels struct {
	Read_request_worldView_chan chan bool
	Read_request_elev_chan      chan bool
	Read_worldView_chan         chan Worldview
	Read_elev_chan              chan Elevator
}

type UpdateWorldviewChannels struct {
	Set_order_chan          chan elevio.ButtonEvent
	Delete_orders_here_chan chan bool
	Arrived_at_floor_chan   chan bool
	Update_direction_chan   chan elevio.MotorDirection
	Version_up_chan         chan bool
	Update_floor_chan       chan int
	Update_to_incoming_chan chan Worldview
	Update_obstr_chan       chan bool
	Cost_func_chan          chan elevio.ButtonEvent
	Peer_lost_chan          chan int
	Peer_new_chan           chan int
}

func InitReadWorldViewChannels() (readChannels ReadWorldviewChannels) {
	readChannels.Read_request_worldView_chan = make(chan bool, 100)
	readChannels.Read_request_elev_chan = make(chan bool, 100)
	readChannels.Read_worldView_chan = make(chan Worldview, 100)
	readChannels.Read_elev_chan = make(chan Elevator, 100)
	return
}

func InitUpdateWorldviewChannels() (updateChannels UpdateWorldviewChannels) {
	updateChannels.Set_order_chan = make(chan elevio.ButtonEvent, 100)
	updateChannels.Delete_orders_here_chan = make(chan bool, 100)
	updateChannels.Arrived_at_floor_chan = make(chan bool, 100)
	updateChannels.Update_direction_chan = make(chan elevio.MotorDirection, 100)
	updateChannels.Version_up_chan = make(chan bool, 100)
	updateChannels.Update_floor_chan = make(chan int, 100)
	updateChannels.Update_to_incoming_chan = make(chan Worldview, 100)
	updateChannels.Update_obstr_chan = make(chan bool, 100)
	updateChannels.Cost_func_chan = make(chan elevio.ButtonEvent, 100)
	updateChannels.Peer_lost_chan = make(chan int, 100)
	updateChannels.Peer_new_chan = make(chan int, 100)
	return
}

func UpdateWorldview(worldView_p *Worldview,
	elev_p *Elevator,
	reset_timer_chan chan<- bool,
	watchdog_chan chan<- bool,
	readChannels ReadWorldviewChannels,
	updateChannels UpdateWorldviewChannels,
	elevioChannels elevio.ElevioChannels) {
	for {
		select {
		case orderButton := <-updateChannels.Set_order_chan:
			setOrder(elev_p, worldView_p, orderButton, reset_timer_chan, watchdog_chan)

		case <-updateChannels.Delete_orders_here_chan:
			deleteOrdersHere(elev_p, worldView_p)

		case <-updateChannels.Arrived_at_floor_chan:
			arrivedAtFloor(elev_p, worldView_p, reset_timer_chan, watchdog_chan)

		case dirUpdate := <-updateChannels.Update_direction_chan:
			elev_p.UpdateDirection(dirUpdate, watchdog_chan)

		case <-updateChannels.Version_up_chan:
			worldView_p.VersionUp()

		case floor := <-updateChannels.Update_floor_chan:
			if floor != -1 {
				elevioChannels.Set_floor_indicator_chan <- floor
				elev_p.Last_Floor = floor
			}
		case incoming := <-updateChannels.Update_to_incoming_chan:
			updateToIncomingVersion(worldView_p, elev_p, incoming)

		case obstr := <-updateChannels.Update_obstr_chan:
			elev_p.Obstruction = obstr

		case buttn := <-updateChannels.Cost_func_chan:
			costFunction(worldView_p, buttn)

		case peer := <-updateChannels.Peer_lost_chan:
			fmt.Println("mottatt fra channel Peer_lost")
			peerLost(peer, updateChannels.Cost_func_chan, updateChannels.Version_up_chan, worldView_p)

		case peer := <-updateChannels.Peer_new_chan:
			peerNew(peer, updateChannels.Version_up_chan, worldView_p)

		case <-readChannels.Read_request_worldView_chan:
			readChannels.Read_worldView_chan <- *worldView_p
		case <-readChannels.Read_request_elev_chan:
			readChannels.Read_elev_chan <- *elev_p
		}
	}
}

func deleteOrdersHere(elev_p *Elevator, worldView_p *Worldview) {
	for orderType := 0; orderType < 2; orderType++ {
		if worldView_p.HallRequests[elev_p.Last_Floor][orderType] == uint8(elev_p.ElevNum) {
			worldView_p.HallRequests[elev_p.Last_Floor][orderType] = 0
		}
	}
	elev_p.CabRequests[elev_p.Last_Floor] = false
	worldView_p.VersionUp()
}

func setOrder(elev_p *Elevator, worldView_p *Worldview, buttn elevio.ButtonEvent, reset_timer_chan chan<- bool, watchdog_chan chan<- bool) {
	if (buttn.Floor == elev_p.Last_Floor) && (elevio.GetFloor() != -1) {
		arrivedAtFloor(elev_p, worldView_p, reset_timer_chan, watchdog_chan)
		return
	}
	if buttn.Button == elevio.BT_CAB {
		elev_p.CabRequests[buttn.Floor] = true
	} else if worldView_p.HallRequests[buttn.Floor][buttn.Button] == 0 {
		costFunction(worldView_p, buttn)
	}
	worldView_p.VersionUp()
}

func arrivedAtFloor(elev_p *Elevator, worldView_p *Worldview, reset_timer_chan chan<- bool, watchdog_chan chan<- bool) {
	reset_timer_chan <- true
	elevio.SetDoorOpenLamp(true)
	elevio.SetMotorDirection(elevio.MD_STOP)
	elev_p.Dirn = elevio.MD_STOP
	deleteOrdersHere(elev_p, worldView_p)
	worldView_p.VersionUp()
	elev_p.Behaviour = EB_DOOR_OPEN
	watchdog_chan <- true
}

func updateToIncomingVersion(myWorldView_p *Worldview, myElev_p *Elevator, incomingWorldView Worldview) {
	myWorldView_p.HallRequests = incomingWorldView.HallRequests
	myWorldView_p.Version = incomingWorldView.Version
	myWorldView_p.ElevList = incomingWorldView.ElevList
	myElev_p.CabRequests = incomingWorldView.ElevList[myElev_p.ElevNum-1].CabRequests
}

func ReadWorldView(readChannels ReadWorldviewChannels) Worldview {
	readChannels.Read_request_worldView_chan <- true
	return <-readChannels.Read_worldView_chan
}

func ReadElevator(readChannels ReadWorldviewChannels) Elevator {
	readChannels.Read_request_elev_chan <- true
	return <-readChannels.Read_elev_chan
}

func peerLost(peer int, cost_func_chan chan elevio.ButtonEvent, version_up_chan chan bool, worldView_p *Worldview) {
	fmt.Println("peerLost starter")
	//worldView := ReadWorldView(readChannels)
	worldView_p.ElevList[peer-1].Online = false
	fmt.Println("peer lost")
	//Assign hall orders to others:
	for floor, f := range worldView_p.HallRequests {
		for buttonType, o := range f {
			if o == uint8(peer) {
				cost_func_chan <- elevio.ButtonEvent{Floor: floor, Button: elevio.ButtonType(buttonType)}
			}
		}
	}
	version_up_chan <- true
}

func peerNew(peer int, version_up_chan chan bool, worldView_p *Worldview) {
	worldView_p.ElevList[peer-1].Online = true
	version_up_chan <- true
}

func (worldView_p *Worldview) VersionUp() {
	if worldView_p.Version < VERSIONLIMIT {
		worldView_p.Version++
	} else {
		worldView_p.Version = 0
	}
}

func (worldview Worldview) Display() {
	fmt.Printf("Sender: %d\n", worldview.Sender)
	fmt.Printf("Version: %d\n", worldview.Version)

	fmt.Println("Elevator List:")
	for i, elev := range worldview.ElevList {
		fmt.Printf("  Elevator %d:\n", i+1)
		fmt.Printf("    Online: %v\n", elev.Online)
		fmt.Printf("    Operative: %v\n", elev.Operative)
		fmt.Printf("    Behaviour: %v\n", elev.Behaviour)
		fmt.Printf("    Obstruction: %t\n", elev.Obstruction)
		fmt.Printf("    ElevNum: %d\n", elev.ElevNum)
		fmt.Printf("    Dirn: %v\n", elev.Dirn)
		fmt.Printf("    Last_dir: %v\n", elev.Last_dir)
		fmt.Printf("    Last_Floor: %d\n", elev.Last_Floor)
		fmt.Printf("    CabRequests: %v\n", elev.CabRequests)
	}

	fmt.Println("Hall Requests:")
	for i := len(worldview.HallRequests); i > 0; i-- {
		fmt.Printf("floor: %d \thall up: %d\t, halldown: %d\n", i-1, worldview.HallRequests[i-1][0], worldview.HallRequests[i-1][1])
	}
}
