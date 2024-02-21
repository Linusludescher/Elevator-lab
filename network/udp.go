package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// Globals
const ( 
    UDPSendPort string = ":20024"
    N_FLOORS int = 4
    N_BUTTONS int = 3
)


// NB. hardcoded values
type Packet struct {
    Version     int         `json:"Version"`
    ElevatorNum int         `json:"ElevatorNum"`
    Guid        int         `json:"Guid"`
    Queue       [N_FLOORS][N_BUTTONS]int   `json:"Queue"`
}

func (packet *Packet) display() {
    fmt.Printf("Elevator number: \t%v\n", packet.ElevatorNum)
    fmt.Printf("Version: \t\t%v\n", packet.Version)
    fmt.Printf("ID: \t\t\t%v\n", packet.Guid)
    fmt.Println("Floor \t Hall Up \t Hall Down \t Cab")
    for i := 0; i < N_FLOORS; i++ {
        fmt.Printf("%v \t %v \t\t %v \t\t %v \t\n", i + 1, packet.Queue[i][0], packet.Queue[i][1], packet.Queue[i][2])
    }
}

func listener(port string) {
	pc, err := net.ListenPacket("udp", port)
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	buf := make([]byte, 1024)

    for {
        n, addr, err := pc.ReadFrom(buf)
        if err != nil {
            panic(err)
        }
        // print IP addres
        fmt.Println(addr.String())

        packet := jsonDecodeElevatorData(buf[:n])
        packet.display()
	}

}

func jsonDecodeElevatorData(jsonPacket []byte) (packet Packet) {
    err := json.Unmarshal(jsonPacket, &packet)
    if err != nil {
        fmt.Println(err)
    }
    return packet
}


func udpInitDail(port string) net.Conn {
	conn, err := net.Dial("udp", port)
	if err != nil {
		panic(err)
	}
    return conn
}

func jsonEncodeElevatorData(version int, elevatorNum int, guid int, queue [N_FLOORS][N_BUTTONS]int) []byte { 
    packet := &Packet{Version: version,
                      ElevatorNum: elevatorNum,
                      Guid:guid,
                      Queue: queue}
    marshaldPacket, err := json.Marshal(packet)              
    if err != nil {
        panic(err)
    }
    return marshaldPacket
} 

func main() {
    go listener(UDPSendPort)

    var conn net.Conn = udpInitDail(UDPSendPort)
    defer conn.Close()

    var testQueue = [4][3]int{
        {1, 2, 3},
        {3, 4, 5},
        {1, 4, 5},
        {1, 2, 5},
    }

    var packet []byte = jsonEncodeElevatorData(1, 1, 1, testQueue)

   for { 
    _, err := conn.Write(packet)
    if err != nil {
        fmt.Println(err)
    }
    time.Sleep(1000 * time.Millisecond)
    }

}
