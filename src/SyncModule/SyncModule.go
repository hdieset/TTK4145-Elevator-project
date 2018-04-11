package syncmodule

import (
	"fmt"
	."param"
)

func SyncModule (localElevatorID string) {
	var localSyncArray SyncArray
	//localSyncArray := new(SyncArray)  - gammel versjon førte til at vi måtte sende peker 
	localSyncArray.AllElevators = make(map[string]Elevator)
	localSyncArray.Owner = localElevatorID	

	for {
		select { 
		case bla: // Får syncArray fra nettverk
			// - Sjekke owner, oppdatere owner i lokale syncarray

		case yolo : // Får peerlistupdate fra nettverk

		case bleu : // Får upd8 fra SingleElevator

		case noldus : // får buttonevent fra SingleElevator



		} 
	}



}



// https://blog.golang.org/go-maps-in-action
/*
m["route"] = 66
n := len(m) <- returns the number of items in the map
delete(m, "route") <- deletes index "route" of object m

*/
