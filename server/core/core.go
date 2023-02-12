package core

import (
	"time"
)

func Start() {
	go clientPinger()
}

func clientPinger() {

	for {
		time.Sleep(1000 * time.Millisecond)
		wb := buildPingWireBytes()

		// Prepare a common packet first
		dm := DispatcherMessage{SourceID: "CORE", Data: wb}

		usernames := GetConnectedUsernames()
		for _, username := range usernames {
			dm.SourceID = username
			toClientDispatcher(dm)
		}
	}
}
