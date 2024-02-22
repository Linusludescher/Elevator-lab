package main

import (
	"main/network"
	"time"
)


func main() {
    ports := network.GetNetworkConfig()
    
    for _, port := range ports.UDPReceve {
        go network.Listener(port)
    }

    conn := network.UdpInitDail(ports.UDPBrodcast)
    
    defer conn.Close()

    var testQueue = [4][3]int{
        {1, 2, 3},
        {3, 4, 5},
        {1, 4, 5},
        {1, 2, 5},
    }

    packet := network.Packet{
        Version: 1,
        ElevatorNum: 2,
        Guid: 23,
        Queue: testQueue,
    }

    for { 
        network.BroadcastPacket(packet, conn)
        time.Sleep(1000 * time.Millisecond)
    }

}
