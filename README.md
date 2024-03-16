# Elevator-lab
helikopterlab

ny lenke:
https://prod.liveshare.vsengsaas.visualstudio.com/join?56F64805BF6DB492DE0C970C23FDDFDC76D4

For processPairs: 
- kommentere inn i main: 
// processPairConn := bcast.ProcessPairListner(id)
- i default case:
// processPairConn.Write([]byte("42"))
- fjerne panic() i ved watchdog time out! og flere steder sikkert




This is a pure peer to peer system to run multiple elevators.

In order to run the elevator, type:

go run main.go -id <elevator_number>

Each elevator needs to have their own number from 1 to n_elevators.
n_elevators is set to 3 but can be changed in config.json.
One can also change the number of floors, and the ports the elevators communicate over.