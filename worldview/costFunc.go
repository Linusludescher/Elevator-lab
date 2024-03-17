package worldview

import (
	"encoding/json"
	"os/exec"
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

func costFunction(worldView_p *Worldview, buttn elevio.ButtonEvent) {
	hraExecutable := "hall_request_assigner"
	input := worldViewToCfInput(*worldView_p, buttn)
	jsonBytes, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}

	ret, err := exec.Command("../"+hraExecutable, "-i", string(jsonBytes)).CombinedOutput()
	if err != nil {
		panic(err)
	}

	output := new(map[string][][2]bool)
	err = json.Unmarshal(ret, &output)
	if err != nil {
		panic(err)
	}

	for i, e := range *output {
		for j, r := range e {
			for k, b := range r {
				if b {
					num, err := strconv.Atoi(i)
					if err != nil {
						panic(err)
					}
					worldView_p.HallRequests[j][k] = uint8(num)
				}
			}
		}
	}
}

func worldViewToCfInput(worldView Worldview, buttn elevio.ButtonEvent) (input HRAInput) {
	input.States = make(map[string]HRAElevState)
	for _, elev := range worldView.ElevList {
		if elev.Online && elev.Operative {
			elevstate := HRAElevState{string(elev.Behaviour), elev.Last_Floor, elev.Dirn.String(), elev.CabRequests}
			input.States[strconv.Itoa(elev.ElevNum)] = elevstate
		}
	}
	input.HallRequests = make([][2]bool, len(worldView.HallRequests))
	for i := range input.HallRequests {
		input.HallRequests[i] = [2]bool{false, false}
	}
	input.HallRequests[buttn.Floor][buttn.Button] = true
	return
}
