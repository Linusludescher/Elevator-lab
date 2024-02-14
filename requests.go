package main

type DirnBehaviourPair struct {
	Dirn             dir
	ElevatorBehavior behaviour
}

func requestsAbove(e Elevator) bool {
	for f := e.floor + 1; f < N_FLOORS; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if e.requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func requestsBelow(e Elevator) bool {
	for f := 0; f < e.floor; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if e.requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func requestsHere(e Elevator) bool {
	for btn := 0; btn < N_BUTTONS; btn++ {
		if e.requests[e.floor][btn] {
			return true
		}
	}
	return false
}

func requestsChooseDirection(e Elevator) DirnBehaviourPair {
	switch e.dirn {
	case D_Up:
		return DirnBehaviourPair{
			Dirn: D_Up,
			ElevatorBehavior: func() int {
				if requestsAbove(e) {
					return EB_Moving // Define EB_Moving accordingly
				} else if requestsHere(e) {
					return EB_DoorOpen // Define EB_DoorOpen accordingly
				} else if requestsBelow(e) {
					return EB_Moving // Define EB_Moving accordingly
				}
				return EB_Idle // Define EB_Idle accordingly
			}(),
		}
	case D_Down:
		return DirnBehaviourPair{
			Dirn: D_Down,
			ElevatorBehavior: func() int {
				if requestsBelow(e) {
					return EB_Moving // Define EB_Moving accordingly
				} else if requestsHere(e) {
					return EB_DoorOpen // Define EB_DoorOpen accordingly
				} else if requestsAbove(e) {
					return EB_Moving // Define EB_Moving accordingly
				}
				return EB_Idle // Define EB_Idle accordingly
			}(),
		}
	case D_Stop:
		return DirnBehaviourPair{
			Dirn: D_Stop,
			ElevatorBehavior: func() int {
				if requestsHere(e) {
					return EB_DoorOpen // Define EB_DoorOpen accordingly
				} else if requestsAbove(e) {
					return EB_Moving // Define EB_Moving accordingly
				} else if requestsBelow(e) {
					return EB_Moving // Define EB_Moving accordingly
				}
				return EB_Idle // Define EB_Idle accordingly
			}(),
		}
	default:
		return DirnBehaviourPair{
			Dirn:             D_Stop,
			ElevatorBehavior: EB_Idle, // Define EB_Idle accordingly
		}
	}
}

func requestsShouldStop(e Elevator) bool {
	switch e.dirn {
	case D_Down:
		return e.requests[e.floor][B_HallDown] || e.requests[e.floor][B_Cab] || !requestsBelow(e)
	case D_Up:
		return e.requests[e.floor][B_HallUp] || e.requests[e.floor][B_Cab] || !requestsAbove(e)
	case D_Stop:
		return true
	default:
		return false
	}
}

func requestsShouldClearImmediately(e Elevator, btnFloor int, btnType int) bool {
	switch e.config.clearRequestVariant {
	case CV_All: // Define CV_All accordingly
		return e.floor == btnFloor
	case CV_InDirn: // Define CV_InDirn accordingly
		return e.floor == btnFloor && ((e.dirn == D_Up && btnType == B_HallUp) || (e.dirn == D_Down && btnType == B_HallDown) || e.dirn == D_Stop || btnType == B_Cab)
	default:
		return false
	}
}

func requestsClearAtCurrentFloor(e Elevator) Elevator {
	switch e.config.clearRequestVariant {
	case CV_All: // Define CV_All accordingly
		for btn := 0; btn < N_BUTTONS; btn++ {
			e.requests[e.floor][btn] = false
		}
	case CV_InDirn: // Define CV_InDirn accordingly
		e.requests[e.floor][B_Cab] = false
		switch e.dirn {
		case D_Up:
			if !requestsAbove(e) && !e.requests[e.floor][B_HallUp] {
				e.requests[e.floor][B_HallDown] = false
			}
			e.requests[e.floor][B_HallUp] = false
		case D_Down:
			if !requestsBelow(e) && !e.requests[e.floor][B_HallDown] {
				e.requests[e.floor][B_HallUp] = false
			}
			e.requests[e.floor][B_HallDown] = false
		case D_Stop:
			e.requests[e.floor][B_HallUp] = false
			e.requests[e.floor][B_HallDown] = false
		default:
			// Do nothing
		}
	default:
		// Do nothing
	}

	return e
}
