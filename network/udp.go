package network

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

// Globals
const ( 
    N_FLOORS int = 4
    N_BUTTONS int = 3
)

// Information to send over UDP broadcast
type Packet struct {
    Version     int         `json:"Version"`
    ElevatorNum int         `json:"ElevatorNum"`
    Guid        int         `json:"Guid"`
    Queue       [N_FLOORS][N_BUTTONS]int   `json:"Queue"`
}

type UDPPorts struct {
    UDPBrodcast string      `json:"UDPBroadcastPort"`
    UDPReceve   []string    `json:"UDPRecevePorts"` 
}

func GetNetworkConfig() (elevatorUDPPorts UDPPorts) {
    jsonData, err := os.ReadFile("config.json")

    // can't read the config file, try again
    if err != nil {
        fmt.Printf("/network/udp.go: Error reading config file: %s\n", err)
        GetNetworkConfig()
    }
    
    // Parse jsonData into ElevatorPorts struct
    err = json.Unmarshal(jsonData,&elevatorUDPPorts)

    if err != nil {
        fmt.Printf("/network/upd.go: Error unmarshal json data to ElevatorPorts struct: %s\n", err)

        // try again
        GetNetworkConfig()
    }

    return
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

func Listener(port string) {
	pc, err := net.ListenPacket("udp", port)
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	buf := make([]byte, 1024)

    for {
        n, addr, err := pc.ReadFrom(buf)
        if err != nil {
            print(err)
        }
        // print IP addres
        fmt.Println(addr.String())
        fmt.Println(buf[:n])

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


func UdpInitDail(port string) net.Conn {
	conn, err := net.Dial("udp", port)
	if err != nil {
		panic(err)
	}
    return conn
}

func BroadcastPacket(packet Packet, conn net.Conn){
    var jsonPacket []byte = jsonEncodeElevatorData(packet)
    _, err := conn.Write(jsonPacket)
    if err != nil {
        fmt.Println(err)
    }

}

func jsonEncodeElevatorData(packet Packet) []byte { 
    marshaldPacket, err := json.Marshal(packet)              
    if err != nil {
        panic(err)
    }
    return marshaldPacket
} 


