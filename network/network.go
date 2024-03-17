package network

import (
	"encoding/json"
	"fmt"
	"os"
	"project/network/bcast"
	"project/network/peers"
	w "project/worldview"
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
	PeerUpdate_chan   chan peers.PeerUpdate
	PeerTxEnable_chan chan bool
	PacketTx_chan     chan w.Worldview
	PacketRx_chan     chan w.Worldview
}

func getNetworkConfig(id int) (configPorts ConfigUDPPorts) {
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
	ports := getNetworkConfig(id)

	networkChan.PeerUpdate_chan = make(chan peers.PeerUpdate)
	networkChan.PeerTxEnable_chan = make(chan bool)
	networkChan.PacketTx_chan = make(chan w.Worldview)
	networkChan.PacketRx_chan = make(chan w.Worldview)

	go peers.Transmitter(ports.UDPstatusPort, strconv.Itoa(id), networkChan.PeerTxEnable_chan)
	go peers.Receiver(ports.UDPstatusPort, networkChan.PeerUpdate_chan)
	go bcast.Transmitter(ports.UDPTx, networkChan.PacketTx_chan)
	for rxPort := range ports.UDPRx {
		fmt.Printf("rxport %d\n", ports.UDPRx[rxPort])
		go bcast.Receiver(ports.UDPRx[rxPort], networkChan.PacketRx_chan)
	}
	return
}

func PeersOnline(readChannels w.ReadWorldviewChannels, network_chan NetworkChan, updateWorldviewChannels w.UpdateWorldviewChannels) {
	for {
		p := <-network_chan.PeerUpdate_chan
		fmt.Printf("Peer update:\n")
		fmt.Printf("  Peers:    %q\n", p.Peers)
		fmt.Printf("  New:      %q\n", p.New)
		fmt.Printf("  Lost:     %q\n", p.Lost)
		fmt.Printf("  UdpTx: 	%v\n", network_chan.PacketTx_chan)
		fmt.Printf("  UdpRx: 	%v\n", network_chan.PacketRx_chan)

		for _, k := range p.Lost {
			k_int, err := strconv.Atoi(k)
			if err != nil {
				panic(err)
			}
			updateWorldviewChannels.Peer_lost_chan <- k_int
		}
		if p.New != "" {
			i, err := strconv.Atoi(p.New)
			if err != nil {
				panic(err)
			}
			updateWorldviewChannels.Peer_new_chan <- i
		}
	}
}
