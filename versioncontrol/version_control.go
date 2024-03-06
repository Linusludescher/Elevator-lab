package versioncontrol

import (
	"project/elevator"
	"project/network"
	"project/network/peers"
	"strconv"
)

// Alle funksjoner er med stor bokstav. Vi må sjekke hva som trenger å være eksternt

const (
	versionLimit = 18446744073709551615 //2e64-1
	// versionStabilityCycles = 100 //maks antall sykler ny versjon kan være foran for at e.Version settes godtar lavere p.Version (ved Version overflow)
	// versionInitVal = 10000 //initialisere på høyere verdi enn 0 for ikke problemer med nullstilling ved tilbakekobling etter utfall

    // Midlertildig deklarering av maks antal heiser på netwerk
    MAX_ELEVATORS_ON_NETWORK = 3
)

func Version_up(e *elevator.Elevator) {
	if e.Version < versionLimit {
		e.Version++
	} else {
		e.Version = 0
	}
}

// Hvorfor sjekker vi egentlig køen?
func Version_if_equal_queue(e elevator.Elevator, p network.Packet) bool {
	for i := range e.Requests {
		for j := range e.Requests[i] {
			if e.Requests[i][j] != p.Queue[i][j] {
                return true
			}
		}
	}
	return false
}


func Version_update_queue(e *elevator.Elevator, p network.Packet) {
	if p.Version > e.Version { //||  (e.Version > versionLimit-versionStabilityCycles && p.Version < versionStabilityCycles)
		e.Requests = p.Queue
		e.Version = p.Version
	} else if (p.Version == e.Version) || !Version_if_equal_queue(*e, p) {
		if p.ElevatorNum > e.ElevNum {
			e.Requests = p.Queue
			e.Version = p.Version
		}
	}
} // må ha noe med når version nullstilles



// Kode fra johannes
// Er drit kode, men prøver å skjønne strukturen på stm for versionskontrolg

// Er fortsatt usikker på hovrdan vi skal løse det at en nye heis blir med netverket. 
// Men tror kanskje det ikke er så vaneksleig, Vi kan ha logik på at nye heiser må gota aktuell version før de får gjøre endiringer.

// Usikker på hvordan vi skal forholde oss til verifisering at et ordre er blit registrert hos de andre heisene.
// Lagre siste ordre fra en heis som FIFO. Også sjekke om ny version er versionen main heis sendte. 
// Hvis ja, slett et elemet fra siste ordre køen.
// Trenger en funksjon som sjekker om vår ordre gikk igjenom. (Anser jeg som kompleks)


// Vi må også lagre tilstanden til heisene som ikke er på nettet, slik at de kan gjennopprete seg til øsnket tilstand.

// Trenger en funskjon som paser på at ordre av heiser som ikke lengre er på nettet også blir tatt. (Anser jeg som kompleks)

type current_elevator struct {
    version uint64 
    sender  int
    queue   [][]uint8
}


func check_if_versions_are_equal(elevators []current_elevator, main_elevator elevator.Elevator) bool {
    for _, elevator := range elevators {
        if elevator.version != main_elevator.Version{
            return false
        }
    }

    return true
}

func check_version_limit(cur_version uint64) uint64 {
    if cur_version == versionLimit {
        return 0
    } else {
        return cur_version
    }

}
// Tar in arry med nyeste version fra alle heiser og heiser på netwerket og oppdaterer egen Version_control
// Vi trenger å holde styr på den nyeste versjonen fra hvær heis
func Version_control(current_elevators_on_network peers.PeerUpdate, current_elevators_version [MAX_ELEVATORS_ON_NETWORK - 1]elevator.Elevator, main_elevator elevator.Elevator) (new_version elevator.Elevator) {
    // Trivial case 
    // Main elevator is the only elevator on network 
    if len(current_elevators_on_network.Peers) == 0{
        new_version = main_elevator
        return
    }

    // Lag en liste med alle heiser på netwerket 
    // Burde være en egen funksjon
    var current_elevators = []current_elevator {}
    for _, onNetwork := range current_elevators_on_network.Peers{
        for _, elevator_version := range current_elevators_version{
            if onNetwork == strconv.Itoa(elevator_version.ElevNum){
                current_elevators = append(current_elevators, current_elevator{elevator_version.Version,
                elevator_version.ElevNum,
                elevator_version.Requests})
            }
        }

    }

    // Are version equal? No - change version to a new version
    if check_if_versions_are_equal(current_elevators, main_elevator) {
        for _, elev := range current_elevators {
            if elev.version > main_elevator.Version || elev.version == check_version_limit(main_elevator.Version) {
                main_elevator.Sender = elev.sender
                main_elevator.Version = elev.version
                main_elevator.Requests = elev.queue
                new_version = main_elevator
                return
            }
        }
    } 


    // Shit, same version

    // Check if all active versions have the same the same sender
    // If yes, then do nothing.

    // Different senders, chose first in current_elevators (can't be empty)
    return
}
