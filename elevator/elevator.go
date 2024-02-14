package elevator

import "fmt"

const (
    N_FLOORS int = 4
    N_BUTTONS int = 3
)

type Elevator_behaviour int
const (
    EB_Idle Elevator_behaviour = iota + 1
    EB_DoorOpen 
    EB_Moving
)

type Elevator_direction int
const (
    D_Down Elevator_direction = iota + 1
    D_Stop 
    D_Up
)

type Elevator struct {
    floor           int
    dirn            Elevator_direction
    behaviour       Elevator_behaviour
    requests        [N_FLOORS][N_BUTTONS]int
}

func Elevator_uninitialized() (elevator Elevator) {
    var matrix [N_FLOORS][N_BUTTONS]int
    elevator = Elevator{-1, D_Stop, EB_Idle, matrix}
    return
}

func (behaviour Elevator_behaviour) String() string {
    return [...]string{"idle", "doorOpen", "moving"}[behaviour - 1]
}

func (direction Elevator_direction) String() string {
    return [...]string{"down", "stop", "up"}[direction]
}

func (elevator *Elevator) Display () {
    fmt.Printf("Floor: %v\n", elevator.floor)
    fmt.Printf("Direction: %s\n", elevator.dirn.String())
    fmt.Printf("Behaviour: %s\n", elevator.behaviour.String())
    fmt.Println("Requests")
    fmt.Println("Floor \t Hall Up \t Hall Down \t Cab")
    for i := 0; i < N_FLOORS; i++ {
        fmt.Printf("%v \t %v \t\t %v \t\t %v \t\n", i + 1, elevator.requests[i][0], elevator.requests[i][1], elevator.requests[i][2])
    }
}

