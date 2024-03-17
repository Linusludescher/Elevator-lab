# Elevator-lab

This is a pure peer to peer system to run multiple elevators.

Firstly, you have to build the program by:

        go build -o elevator main.go

Then to run it:

        bash ./runelevator.sh <id>

Each elevator needs to have their own id from 1 to n_elevators.
n_elevators is set to 3 but can be changed in config.json.
Also use config.jason to change the number of floors, and the ports the elevators communicate over.  
Before starting the program, you need to download the executable: hall_request_assigner, found [here](https://github.com/TTK4145/Project-resources/releases/tag/v1.1.1)
The hall_request_assigner must be placed in the the parentfolder of Elevator-lab.
You might have to give permission to run the hall_request_assigner. To do this you can run: chmod a+xw hall_request_assigner.