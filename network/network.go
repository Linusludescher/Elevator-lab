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
	"strconv"
)

// // Globals
// const (
// 	N_FLOORS  int = 4
// 	N_BUTTONS int = 3
// )

type ConfigUDPPorts struct {
	UDPBase       int `json:"BasePort"`
	N_elevators   int `json:"n_elevators"`
	UDPTx         int
	UDPRx         []int
	Id            int `json:"ElevNum"`
	UDPstatusPort int `json:"StatusPort"`
}

type NetworkChan struct {
	PeerUpdateCh chan peers.PeerUpdate
	PeerTxEnable chan bool
	PacketTx     chan elevator.Elevator
	PacketRx     chan elevator.Elevator
}

func getNetworkConfig() (cp ConfigUDPPorts) {
	jsonData, err := os.ReadFile("config.json")

	// can't read the config file, try again
	if err != nil {
		fmt.Printf("/network/udp.go: Error reading config file: %s\n", err)
		getNetworkConfig()
	}

	// Parse jsonData into ElevatorPorts struct
	err = json.Unmarshal(jsonData, &cp)
	for i := 1; i < cp.N_elevators+1; i++ {
		if i == cp.Id {
			cp.UDPTx = cp.UDPBase + cp.Id
		} else {
			cp.UDPRx = append(cp.UDPRx, cp.UDPBase+i)
		}
	}
	if err != nil {
		fmt.Printf("/network/upd.go: Error unmarshal json data to ElevatorPorts struct: %s\n", err)

		// try again
		getNetworkConfig()
	}

	return

}

func Init_network() (networkChan NetworkChan) {
	// Read from config.json port addresses for Rx and Tx
	ports := getNetworkConfig()
	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`
	var id string = strconv.Itoa(ports.Id)
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
	networkChan.PacketTx = make(chan elevator.Elevator)
	networkChan.PacketRx = make(chan elevator.Elevator)
	// ... and start the transmitter/receiver pair on some port
	// These functions can take any number of channels! It is also possible to
	//  start multiple transmitters/receivers on the same port.
	go bcast.Transmitter(ports.UDPTx, networkChan.PacketTx)

	for rxPort := range ports.UDPRx {
		fmt.Printf("rxport %d\n", ports.UDPRx[rxPort])
		go bcast.Receiver(ports.UDPRx[rxPort], networkChan.PacketRx)
	}

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
				fmt.Printf("  UdpTx: 	%d\n", ports.UDPTx)
				fmt.Printf("  UdpRx: 	%d\n", ports.UDPRx)
			case a := <-networkChan.PacketRx:
				fmt.Println("Received:")
				a.Display() // feilmelding hvis a ikke er en struct Packet
			}
		}
	}()
	return
}
