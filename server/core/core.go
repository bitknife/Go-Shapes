package core

import (
	"time"
)

const (
	PING_INTERVAL = 1000
)

func Start() {
	go clientPinger()
}

func clientPinger() {

	for {
		time.Sleep(PING_INTERVAL * time.Millisecond)

		dm := DispatcherMessage{SourceID: "CORE", Packet: buildPingPacket()}

		// TODO: Improve? for each connected client
		usernames := GetConnectedUsernames()
		for _, username := range usernames {
			dm.SourceID = username
			toClientDispatcher(dm)
		}
	}
}
