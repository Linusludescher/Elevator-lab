package elevator

import "project/elevio"

func UpdateWorldview(worldView_p *Worldview,
					elev_p *Elevator,
					reset_timer_chan chan <- bool,
					watchdog_chan chan <- bool,
					set_order_chan <- chan elevio.ButtonEvent, 
					delete_orders_here_chan <- chan bool, 
					arrived_at_floor_chan <- chan bool,
					update_direction_chan <- chan elevio.MotorDirection,
					version_up_chan <- chan bool, ){
	for{
		select{
		case orderButton := <- set_order_chan:
			SetOrder(elev_p, worldView_p, orderButton, reset_timer_chan, watchdog_chan)
		case <- delete_orders_here_chan:
			DeleteOrdersHere(elev_p, worldView_p)
		case <- arrived_at_floor_chan: 
			ArrivedAtFloor(elev_p, worldView_p, reset_timer_chan, watchdog_chan)
		case dirUpdate := <- update_direction_chan:
			 elev_p.UpdateDirection(dirUpdate, watchdog_chan)
		case <- version_up_chan: 
			worldView_p.VersionUp()
		}
	}
}

func DeleteOrdersHere(elev_p *Elevator, worldView_p *Worldview) {
	for orderType := 0; orderType < 2; orderType++ {
		if worldView_p.HallRequests[elev_p.Last_Floor][orderType] == uint8(elev_p.ElevNum) {
			worldView_p.HallRequests[elev_p.Last_Floor][orderType] = 0
		}
	}
	elev_p.CabRequests[elev_p.Last_Floor] = false
	worldView_p.VersionUp()
}

func SetOrder(elev_p *Elevator, worldView_p *Worldview, buttn elevio.ButtonEvent, reset_timer_chan chan <- bool, wd_chan chan <- bool) {
	if (buttn.Floor == elev_p.Last_Floor) && (elevio.GetFloor() != -1) {
		ArrivedAtFloor(elev_p, worldView_p, reset_timer_chan, wd_chan)
		return
	}
	if buttn.Button == elevio.BT_CAB {
		elev_p.CabRequests[buttn.Floor] = true
	} else if worldView_p.HallRequests[buttn.Floor][buttn.Button] == 0 {
		//costFunc.CostFunction(worldView_p, buttn)
	}
	worldView_p.VersionUp()
}

func ArrivedAtFloor(elev_p *Elevator, worldView_p *Worldview, reset_timer_chan chan<- bool, wd_chan chan<- bool) {
	elevio.SetDoorOpenLamp(true)
	elevio.SetMotorDirection(elevio.MD_STOP)
	elev_p.Dirn = elevio.MD_STOP
	DeleteOrdersHere(elev_p, worldView_p)
	worldView_p.VersionUp()
	elev_p.Behaviour = EB_DOOR_OPEN
	wd_chan <- true
	reset_timer_chan <- true
}