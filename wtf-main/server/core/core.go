package core

import (
	"bitknife.se/wtf/shared"
	"log"
	"time"
)

const (
	PING_INTERVAL = 1000
)

func Run() {
	/**
	Core routine(s) are respinsible for all things not related to "gamey" stuff.
		- Housekeeping
		- Metrics collection/delivery
		- ...
	*/
	go clientPinger()
}

func clientPinger() {

	for {
		time.Sleep(PING_INTERVAL * time.Millisecond)

		dm := DispatcherMessage{SourceID: "CORE", Packet: shared.BuildPingPacket()}

		// TODO: Improve? for each connected client
		usernames := GetConnectedUsernames()
		for _, username := range usernames {
			dm.SourceID = username
			toClientDispatcher(dm)
		}
	}
	log.Panicln("Pinger returns!")
}
