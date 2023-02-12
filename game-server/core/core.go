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
		packetData := buildPingPacket()

		// Prepare a common packet first
		dm := DispatcherMessage{SourceID: "", Type: int(MType_PING_EVENT), Data: packetData}

		usernames := GetConnectedUsernames()
		for _, username := range usernames {
			dm.SourceID = username
			toClientDispatcher(dm)
		}
	}
}
