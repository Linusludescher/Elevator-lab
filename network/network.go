package network

import (
	"encoding/json"
	"fmt"
	"os"
	"project/costFunc"
	"project/elevator"
	"project/elevio"
	"project/network/bcast"
	"project/network/peers"
	"strconv"
)

type ConfigUDPPorts struct {
	UDPBase       int `json:"BasePort"`
	N_elevators   int `json:"n_elevators"`
	UDPTx         int
	UDPRx         []int
	Id            int
	UDPstatusPort int `json:"StatusPort"`
}

type NetworkChan struct {
	PeerUpdateCh chan peers.PeerUpdate
	PeerTxEnable chan bool
	PacketTx     chan elevator.Worldview
	PacketRx     chan elevator.Worldview
}

func getNetworkConfig(id int) (configPorts ConfigUDPPorts) { //mer beskrivende navn til configPorts!
	jsonData, err := os.ReadFile("config.json")
	configPorts.Id = id
	// can't read the config file, try again
	if err != nil {
		fmt.Printf("/network/udp.go: Error reading config file: %s\n", err)
		getNetworkConfig(id)
	}

	// Parse jsonData into ElevatorPorts struct
	err = json.Unmarshal(jsonData, &configPorts)
	if err != nil {
		fmt.Printf("/network/upd.go: Error unmarshal json data to ElevatorPorts struct: %s\n", err)

		// try again
		getNetworkConfig(id)
	}
	for i := 1; i < configPorts.N_elevators+1; i++ {
		if i == configPorts.Id {
			configPorts.UDPTx = configPorts.UDPBase + configPorts.Id
		} else {
			configPorts.UDPRx = append(configPorts.UDPRx, configPorts.UDPBase+i)
		}
	}
	return
}

func InitNetwork(id int) (networkChan NetworkChan) {
	// Read from config.json port addresses for Rx and Tx
	ports := getNetworkConfig(id)

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
	return
}

func PeersOnline(worldView_p *elevator.Worldview, network_chan NetworkChan) {
	fmt.Println("Started")
	for {
		p := <-network_chan.PeerUpdateCh
		fmt.Printf("Peer update:\n")
		fmt.Printf("  Peers:    %q\n", p.Peers)
		fmt.Printf("  New:      %q\n", p.New)
		fmt.Printf("  Lost:     %q\n", p.Lost)
		fmt.Printf("  UdpTx: 	%d\n", network_chan.PacketTx)
		fmt.Printf("  UdpRx: 	%d\n", network_chan.PacketRx)

		for _, k := range p.Lost {
			k_int, err := strconv.Atoi(k)
			if err != nil {
				panic(err)
			}
			worldView_p.ElevList[k_int-1].Online = false

			//Assign hall orders to others:
			for floor, f := range worldView_p.HallRequests {
				for buttonType, o := range f {
					if o == uint8(k_int) {
						buttn := elevio.ButtonEvent{Floor: floor, Button: elevio.ButtonType(buttonType)}
						costFunc.CostFunction(worldView_p, buttn)
					}
				}
			}
			worldView_p.Version_up()
		}
		if p.New != "" {
			i, err := strconv.Atoi(p.New)
			if err != nil {
				panic(err)
			}
			worldView_p.ElevList[i-1].Online = true

			worldView_p.Version_up()
		}
	}
}
