package costFunc

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"project/elevator"
	"project/elevio"
	"strconv"
)

type HRAElevState struct {
	Behavior    string `json:"behaviour"`
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	CabRequests []bool `json:"cabRequests"`
}

type HRAInput struct {
	HallRequests [][2]bool               `json:"hallRequests"`
	States       map[string]HRAElevState `json:"states"`
}

func CostFunction(wv *elevator.Worldview, buttn elevio.ButtonEvent) {
	hraExecutable := "hall_request_assigner"
	input := wvToCfInput(*wv, buttn)
	jsonBytes, err := json.Marshal(input)
	if err != nil {
		fmt.Println("json.Marshal error: ", err)
		return
	}

	ret, err := exec.Command("../hall_request_assigner/"+hraExecutable, "-i", string(jsonBytes)).CombinedOutput()
	if err != nil {
		fmt.Println("exec.Command error: ", err)
		fmt.Println(string(ret))
		return
	}

	output := new(map[string][][2]bool)
	err = json.Unmarshal(ret, &output)
	if err != nil {
		fmt.Println("json.Unmarshal error: ", err)
		return
	}

	fmt.Printf("output: \n")
	for k, v := range *output {
		fmt.Printf("%6v :  %+v\n", k, v)
	}

	//gjøre noe med output
}

func wvToCfInput(wv elevator.Worldview, buttn elevio.ButtonEvent) (input HRAInput) {
	input.States = make(map[string]HRAElevState)
	for _, elev := range wv.ElevList {
		elevstate := HRAElevState{string(elev.Behaviour), elev.Last_Floor, strconv.Itoa(int(elev.Dirn)), elev.CabRequests}
		input.States[strconv.Itoa(elev.ElevNum)] = elevstate
	}
	input.HallRequests[buttn.Floor][buttn.Button] = true
	return
}
