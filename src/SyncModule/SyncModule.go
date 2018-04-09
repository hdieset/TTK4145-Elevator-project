package syncmodule

import (
	"fmt"
	."param"
)

func SyncModule () {
	localSyncArray := new(SyncArray)
	localSyncArray.AllElevators = make(map[string]Elevator)


}



// https://blog.golang.org/go-maps-in-action
/*
m["route"] = 66
n := len(m) <- returns the number of items in the map
delete(m, "route") <- deletes index "route" of object m

*/
