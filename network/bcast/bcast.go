package bcast

import (
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"path/filepath"
	"project/elevator"
	"project/network/conn"
	"reflect"
	"runtime"
	"strconv"
	"time"
)

const BUFSIZE = 1024

// Encodes received values from `chans` into type-tagged JSON, then broadcasts
// it on `port`
func Transmitter(port int, chans ...interface{}) {
	checkArgs(chans...)
	typeNames := make([]string, len(chans))
	selectCases := make([]reflect.SelectCase, len(typeNames))
	for i, ch := range chans {
		selectCases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		}
		typeNames[i] = reflect.TypeOf(ch).Elem().String()
	}

	conn := conn.DialBroadcastUDP(port)
	addr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("255.255.255.255:%d", port))
	for {
		chosen, value, _ := reflect.Select(selectCases)
		jsonstr, _ := json.Marshal(value.Interface())
		ttj, _ := json.Marshal(typeTaggedJSON{
			TypeId: typeNames[chosen],
			JSON:   jsonstr,
		})
		if len(ttj) > BUFSIZE {
			panic(fmt.Sprintf(
				"Tried to send a message longer than the buffer size (length: %d, buffer size: %d)\n\t'%s'\n"+
					"Either send smaller packets, or go to network/bcast/bcast.go and increase the buffer size",
				len(ttj), BUFSIZE, string(ttj)))
		}
		conn.WriteTo(ttj, addr)
	}
}

// Matches type-tagged JSON received on `port` to element types of `chans`, then
// sends the decoded value on the corresponding channel
func Receiver(port int, chans ...interface{}) {
	checkArgs(chans...)
	chansMap := make(map[string]interface{})
	for _, ch := range chans {
		chansMap[reflect.TypeOf(ch).Elem().String()] = ch
	}
	var buf [BUFSIZE]byte
	conn := conn.DialBroadcastUDP(port)
	for {
		n, _, e := conn.ReadFrom(buf[0:])
		if e != nil {
			fmt.Printf("bcast.Receiver(%d, ...):ReadFrom() failed: \"%+v\"\n", port, e)
		}
		var ttj typeTaggedJSON
		json.Unmarshal(buf[0:n], &ttj)
		ch, ok := chansMap[ttj.TypeId]
		if !ok {
			continue
		}
		v := reflect.New(reflect.TypeOf(ch).Elem())
		json.Unmarshal(ttj.JSON, v.Interface())
		reflect.Select([]reflect.SelectCase{{
			Dir:  reflect.SelectSend,
			Chan: reflect.ValueOf(ch),
			Send: reflect.Indirect(v),
		}})
	}
}

type typeTaggedJSON struct {
	TypeId string
	JSON   []byte
}

// Checks that args to Tx'er/Rx'er are valid:
//
//	All args must be channels
//	Element types of channels must be encodable with JSON
//	No element types are repeated
//
// Implementation note:
//   - Why there is no `isMarshalable()` function in encoding/json is a mystery,
//     so the tests on element type are hand-copied from `encoding/json/encode.go`
func checkArgs(chans ...interface{}) {
	n := 0
	for range chans {
		n++
	}
	elemTypes := make([]reflect.Type, n)

	for i, ch := range chans {
		// Must be a channel
		if reflect.ValueOf(ch).Kind() != reflect.Chan {
			panic(fmt.Sprintf(
				"Argument must be a channel, got '%s' instead (arg# %d)",
				reflect.TypeOf(ch).String(), i+1))
		}

		elemType := reflect.TypeOf(ch).Elem()

		// Element type must not be repeated
		for j, e := range elemTypes {
			if e == elemType {
				panic(fmt.Sprintf(
					"All channels must have mutually different element types, arg# %d and arg# %d both have element type '%s'",
					j+1, i+1, e.String()))
			}
		}
		elemTypes[i] = elemType

		// Element type must be encodable with JSON
		checkTypeRecursive(elemType, []int{i + 1})

	}
}

func checkTypeRecursive(val reflect.Type, offsets []int) {
	switch val.Kind() {
	case reflect.Complex64, reflect.Complex128, reflect.Chan, reflect.Func, reflect.UnsafePointer:
		panic(fmt.Sprintf(
			"Channel element type must be supported by JSON, got '%s' instead (nested arg# %v)",
			val.String(), offsets))
	case reflect.Map:
		if val.Key().Kind() != reflect.String {
			panic(fmt.Sprintf(
				"Channel element type must be supported by JSON, got '%s' instead (map keys must be 'string') (nested arg# %v)",
				val.String(), offsets))
		}
		checkTypeRecursive(val.Elem(), offsets)
	case reflect.Array, reflect.Ptr, reflect.Slice:
		checkTypeRecursive(val.Elem(), offsets)
	case reflect.Struct:
		for idx := 0; idx < val.NumField(); idx++ {
			checkTypeRecursive(val.Field(idx).Type, append(offsets, idx+1))
		}
	}
}

func BcWorldView(elev elevator.Elevator, worldView elevator.Worldview, bc_chan chan<- elevator.Worldview) {
	worldView.ElevList[elev.ElevNum-1] = elev
	bc_chan <- worldView
}

func ProcessPairListner(id int) (udpConn *net.UDPConn) {
	broadcastAddr, err := net.ResolveUDPAddr("udp", "localhost:8091")
	if err != nil {
		panic(err)
	}

	// backup, lytter på UDP
	listen_conn, err := net.ListenUDP("udp", broadcastAddr)
	if err != nil {
		panic(err)
	}
	fmt.Println("Listening..")

	buffer := make([]byte, 1024)

	for {
		timeout := time.Now().Add(10 * time.Second)
		listen_conn.SetReadDeadline(timeout)
		_, _, err := listen_conn.ReadFromUDP(buffer)
		if err != nil {
			// Check if the error is a timeout
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Println("Read timeout occurred. Breaking...")
				break
			}
			panic(err)
		}
	}

	listen_conn.Close()

	// starte nytt vindu
	time.Sleep(1000 * time.Millisecond)
	flag := "-id"
	value := strconv.Itoa(id)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic(ok)
	}
	path_to_main := filepath.Join(filepath.Dir(filename), "..", "..", "main.go")
	cmd := exec.Command("gnome-terminal", "--", "go", "run", path_to_main, flag, value)
	fmt.Println(cmd.Args)

	// Run the command
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println("Im the primary")
	udpConn, err = net.DialUDP("udp", nil, broadcastAddr)
	if err != nil {
		panic(err)
	}
	return
}
