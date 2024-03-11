package network

import (
	"encoding/json"
	"fmt"
	"os"
	"project/elevator"
	"project/network/bcast"
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
	PacketTx     chan elevator.Worldview
	PacketRx     chan elevator.Worldview
}

func getNetworkConfig() (cp ConfigUDPPorts, Id int) {
	jsonData, err := os.ReadFile("config.json")

	// can't read the config file, try again
	if err != nil {
		fmt.Printf("/network/udp.go: Error reading config file: %s\n", err)
		getNetworkConfig()
	}

	// Parse jsonData into ElevatorPorts struct
	err = json.Unmarshal(jsonData, &cp)
	if err != nil {
		fmt.Printf("/network/upd.go: Error unmarshal json data to ElevatorPorts struct: %s\n", err)

		// try again
		getNetworkConfig()
	}
	for i := 1; i < cp.N_elevators+1; i++ {
		if i == cp.Id {
			cp.UDPTx = cp.UDPBase + cp.Id
		} else {
			cp.UDPRx = append(cp.UDPRx, cp.UDPBase+i)
		}
	}
	Id = cp.Id
	return
}

func Init_network(e *elevator.Elevator, wv *elevator.Worldview) (networkChan NetworkChan) {
	// Read from config.json port addresses for Rx and Tx
	ports, id := getNetworkConfig()

	// We make a channel for receiving updates on the id's of the peers that are
	//  alive on the network
	networkChan.PeerUpdateCh = make(chan peers.PeerUpdate)
	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	networkChan.PeerTxEnable = make(chan bool)

	go peers.Transmitter(ports.UDPstatusPort, strconv.Itoa(id), networkChan.PeerTxEnable)
	go peers.Receiver(ports.UDPstatusPort, networkChan.PeerUpdateCh)

	// We make channels for sending and receiving our custom data types
	networkChan.PacketTx = make(chan elevator.Worldview)
	networkChan.PacketRx = make(chan elevator.Worldview)
	// ... and start the transmitter/receiver pair on some port
	// These functions can take any number of channels! It is also possible to
	//  start multiple transmitters/receivers on the same port.
	go bcast.Transmitter(ports.UDPTx, networkChan.PacketTx)

	for rxPort := range ports.UDPRx {
		fmt.Printf("rxport %d\n", ports.UDPRx[rxPort])
		go bcast.Receiver(ports.UDPRx[rxPort], networkChan.PacketRx)
	}

	// midlertidlig, slik at vi ikke må skrive så mye kode for å teste
	go func(e *elevator.Elevator, wv *elevator.Worldview) {
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
				for _, k := range p.Lost {
					k, err := strconv.Atoi(k)
					if err != nil {
						fmt.Println("Error:", err)
						return
					}
					for i := 0; i < 2; i++ {
						for j := range e.CabRequests[i] {
							if e.CabRequests[j] == 1 {
								//Kost-funksjon
								if wv.HallRequests[i][j] == uint8(k) {
									// Kost fuknsjon
								}
							}
						}
					}
				}
			case <-networkChan.PacketRx:
				fmt.Println("Received:")
				//a.Display() // feilmelding hvis a ikke er en struct Packet
			}
		}
	}(e, wv)
	return
}
