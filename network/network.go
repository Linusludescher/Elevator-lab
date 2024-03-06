package network

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"project/elevator"
	"project/network/bcast"
	"project/network/localip"
	"project/network/peers"
	"time"
)

// Globals
const (
	N_FLOORS  int = 4
	N_BUTTONS int = 3
)

// Information to send over UDP broadcast
type Packet struct {
	Version     uint64
	ElevatorNum int
	Guid        int
	Queue       [][]uint8
}

type UDPPorts struct {
	UDPTx         int    `json:"UDPTx"`
	UDPRx         []int  `json:"UDPRx"`
	Id            string `json:"ElevNum"`
	UDPstatusPort int    `json:"StatusPort"`
}

type NetworkChan struct {
	PeerUpdateCh chan peers.PeerUpdate
	PeerTxEnable chan bool
	PacketTx     chan Packet
	PacketRx     chan Packet
}

func getNetworkConfig() (elevatorUDPPorts UDPPorts) {
	jsonData, err := os.ReadFile("config.json")

	// can't read the config file, try again
	if err != nil {
		fmt.Printf("/network/udp.go: Error reading config file: %s\n", err)
		getNetworkConfig()
	}

	// Parse jsonData into ElevatorPorts struct
	err = json.Unmarshal(jsonData, &elevatorUDPPorts)

	if err != nil {
		fmt.Printf("/network/upd.go: Error unmarshal json data to ElevatorPorts struct: %s\n", err)

		// try again
		getNetworkConfig()
	}

	return

}

func (packet *Packet) Display() {
	fmt.Printf("Elevator number: \t%v\n", packet.ElevatorNum)
	fmt.Printf("Version: \t\t%v\n", packet.Version)
	fmt.Printf("ID: \t\t\t%v\n", packet.Guid)
	fmt.Println("Floor \t Hall Up \t Hall Down \t Cab")
	for i := 0; i < N_FLOORS; i++ {
		fmt.Printf("%v \t %v \t\t %v \t\t %v \t\n", i+1, packet.Queue[i][0], packet.Queue[i][1], packet.Queue[i][2])
	}
}

func Init_network() (networkChan NetworkChan) {
	// Read from config.json port addresses for Rx and Tx
	ports := getNetworkConfig()

	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`
	var id string = ports.Id
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	// ... or alternatively, we can use the local IP address.
	// (But since we can run multiple programs on the same PC, we also append the
	//  process ID)
	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}

	// We make a channel for receiving updates on the id's of the peers that are
	//  alive on the network
	networkChan.PeerUpdateCh = make(chan peers.PeerUpdate)
	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	networkChan.PeerTxEnable = make(chan bool)

	go peers.Transmitter(ports.UDPstatusPort, id, networkChan.PeerTxEnable)
	go peers.Receiver(ports.UDPstatusPort, networkChan.PeerUpdateCh)

	// We make channels for sending and receiving our custom data types
	networkChan.PacketTx = make(chan Packet)
	networkChan.PacketRx = make(chan Packet)
	// ... and start the transmitter/receiver pair on some port
	// These functions can take any number of channels! It is also possible to
	//  start multiple transmitters/receivers on the same port.
	go bcast.Transmitter(ports.UDPTx, networkChan.PacketTx)

	for rxPort := range ports.UDPRx {
		go bcast.Receiver(rxPort, networkChan.PacketRx)
	}

	// The example message. We just send one of these every second.
	go func() {
		var packet Packet
		for {
			packet.Version++
			networkChan.PacketTx <- packet
			time.Sleep(1 * time.Second)
		}
	}()

	// midlertidlig, slik at vi ikke må skrive så mye kode for å teste
	go func() {
		fmt.Println("Started")
		for {
			select {
			case p := <-networkChan.PeerUpdateCh:
				fmt.Printf("Peer update:\n")
				fmt.Printf("  Peers:    %q\n", p.Peers)
				fmt.Printf("  New:      %q\n", p.New)
				fmt.Printf("  Lost:     %q\n", p.Lost)

			case a := <-networkChan.PacketRx:
				fmt.Printf("Received: %#v\n", a)
				a.Display() // feilmelding hvis a ikke er en struct Packet
			}
		}
	}()
	return
}

func Elevator_to_packet(e elevator.Elevator) Packet {
	packet := Packet{
		Version:     e.Version,
		ElevatorNum: e.ElevNum,
		Guid:        0,
		Queue:       e.Requests,
	}

	return packet
}
