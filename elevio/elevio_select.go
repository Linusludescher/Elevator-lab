package elevio

type ButtonLampOrder struct {
	Button_type ButtonType
	OrderFloor  int
	Value       bool
}

type ElevioChannels struct {
	Set_floor_indicator_chan chan int
	Set_motor_direction_chan chan MotorDirection
	Set_button_lamp_chan     chan ButtonLampOrder
	Set_door_open_lamp_chan  chan bool
	Get_floor_chan           chan bool
	Send_floor_chan          chan int
}

func InitElevioChannels() (ioChannels ElevioChannels) {
	ioChannels.Set_floor_indicator_chan = make(chan int)
	ioChannels.Set_motor_direction_chan = make(chan MotorDirection)
	ioChannels.Set_button_lamp_chan = make(chan ButtonLampOrder)
	ioChannels.Set_door_open_lamp_chan = make(chan bool)
	ioChannels.Get_floor_chan = make(chan bool)
	ioChannels.Send_floor_chan = make(chan int)
	return
}

func Elevio_select(ioChannels ElevioChannels) {
	for {
		select {
		case floor := <-ioChannels.Set_floor_indicator_chan:
			SetFloorIndicator(floor)

		case motorDirection := <-ioChannels.Set_motor_direction_chan:
			SetMotorDirection(motorDirection)

		case buttonLamp := <-ioChannels.Set_button_lamp_chan:
			SetButtonLamp(buttonLamp.Button_type, buttonLamp.OrderFloor, buttonLamp.Value)

		case openDoorLamp := <-ioChannels.Set_door_open_lamp_chan:
			SetDoorOpenLamp(openDoorLamp)

		case <-ioChannels.Get_floor_chan:
			ioChannels.Send_floor_chan <- GetFloor()
		}
	}
}
