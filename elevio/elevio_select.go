package elevio

type ButtonLampOrder struct {
	Button_type ButtonType
	OrderFloor  int
	Value       bool
}

func elevio_select(set_floor_indicator_chan <-chan int,
	set_motor_direction_chan <-chan MotorDirection,
	set_button_lamp_chan <-chan ButtonLampOrder,
	set_door_open_lamp_chan <-chan bool,
	get_floor_chan <-chan bool,
	send_floor_chan chan<- int) {
	for {
		select {
		case floor := <-set_floor_indicator_chan:
			SetFloorIndicator(floor)
		case motorDirection := <-set_motor_direction_chan:
			SetMotorDirection(motorDirection)
		case buttonLamp := <-set_button_lamp_chan:
			SetButtonLamp(buttonLamp.Button_type, buttonLamp.OrderFloor, buttonLamp.Value)
		case openDoorLamp := <-set_door_open_lamp_chan:
			SetDoorOpenLamp(openDoorLamp)
		case <-get_floor_chan:
			send_floor_chan <- GetFloor()
		}
	}
}
