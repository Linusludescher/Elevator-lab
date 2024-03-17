package elevio

import (
	"fmt"
	"net"
	"sync"
	"time"
)

const _POLLRATE = 20 * time.Millisecond

var _initialized bool = false
var _numFloors int = 4
var _mtx sync.Mutex
var _conn net.Conn

type MotorDirection int

const (
	MD_UP   MotorDirection = 1
	MD_DOWN MotorDirection = -1
	MD_STOP MotorDirection = 0
)

type ButtonType int

const (
	BT_HALLUP   ButtonType = 0
	BT_HALLDOWN ButtonType = 1
	BT_CAB      ButtonType = 2
)

type ButtonEvent struct {
	Floor  int
	Button ButtonType
}

func (md MotorDirection) String() string {
	switch md {
	case MD_UP:
		return "up"
	case MD_DOWN:
		return "down"
	case MD_STOP:
		return "stop"
	default:
		return "unknown"
	}
}

func Init(addr string, numFloors int) {
	if _initialized {
		fmt.Println("Driver already initialized!")
		return
	}
	_numFloors = numFloors
	_mtx = sync.Mutex{}
	var err error
	_conn, err = net.Dial("tcp", addr)
	if err != nil {
		panic(err.Error())
	}
	_initialized = true
}

func SetMotorDirection(dir MotorDirection) {
	write([4]byte{1, byte(dir), 0, 0})
}

func SetButtonLamp(button ButtonType, floor int, value bool) {
	write([4]byte{2, byte(button), byte(floor), toByte(value)})
}

func SetFloorIndicator(floor int) {
	write([4]byte{3, byte(floor), 0, 0})
}

func SetDoorOpenLamp(value bool) {
	write([4]byte{4, toByte(value), 0, 0})
}

func SetStopLamp(value bool) {
	write([4]byte{5, toByte(value), 0, 0})
}

func PollButtons(receiver_chan chan<- ButtonEvent) {
	prev := make([][3]bool, _numFloors)
	for {
		time.Sleep(_POLLRATE)
		for f := 0; f < _numFloors; f++ {
			for b := ButtonType(0); b < 3; b++ {
				v := GetButton(b, f)
				if v != prev[f][b] && v {
					receiver_chan <- ButtonEvent{f, ButtonType(b)}
				}
				prev[f][b] = v
			}
		}
	}
}

func PollFloorSensor(receiver_chan chan<- int) {
	prev := -1
	for {
		time.Sleep(_POLLRATE)
		v := GetFloor()
		if v != prev && v != -1 {
			receiver_chan <- v
		}
		prev = v
	}
}

func PollStopButton(receiver_chan chan<- bool) {
	prev := false
	for {
		time.Sleep(_POLLRATE)
		v := GetStop()
		if v != prev {
			receiver_chan <- v
		}
		prev = v
	}
}

func PollObstructionSwitch(receiver_chan chan<- bool) {
	prev := false
	for {
		time.Sleep(_POLLRATE)
		v := GetObstruction()
		if v != prev {
			fmt.Println(v)

			receiver_chan <- v
		}
		prev = v
	}
}

func GetButton(button ButtonType, floor int) bool {
	a := read([4]byte{6, byte(button), byte(floor), 0})
	return toBool(a[1])
}

func GetFloor() int {
	a := read([4]byte{7, 0, 0, 0})
	if a[1] != 0 {
		return int(a[2])
	} else {
		return -1
	}
}

func GetStop() bool {
	a := read([4]byte{8, 0, 0, 0})
	return toBool(a[1])
}

func GetObstruction() bool {
	a := read([4]byte{9, 0, 0, 0})
	return toBool(a[1])
}

func read(in [4]byte) [4]byte {
	_mtx.Lock()
	defer _mtx.Unlock()

	_, err := _conn.Write(in[:])
	if err != nil {
		panic("Lost connection to Elevator Server")
	}

	var out [4]byte
	_, err = _conn.Read(out[:])
	if err != nil {
		panic("Lost connection to Elevator Server")
	}

	return out
}

func write(in [4]byte) {
	_mtx.Lock()
	defer _mtx.Unlock()

	_, err := _conn.Write(in[:])
	if err != nil {
		panic("Lost connection to Elevator Server")
	}
}

func toByte(a bool) byte {
	var b byte = 0
	if a {
		b = 1
	}
	return b
}

func toBool(a byte) bool {
	var b bool = false
	if a != 0 {
		b = true
	}
	return b
}
