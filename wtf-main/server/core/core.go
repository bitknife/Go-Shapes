package core

import (
	"bitknife.se/wtf/shared"
	"log"
	"time"
)

const (
	PING_INTERVAL_SEC = 5
)

func Run(pingIntervalMsec int) {
	/**
	Core routine(s) are responsible for all things _NOT_ related to "gamey" stuff.
		- Housekeeping
		- Metrics collection/delivery
		- ...
	*/

	// Just an example for now
	go broadCastPing(pingIntervalMsec)
}

func broadCastPing(pingIntervalMsec int) {

	for {
		time.Sleep(time.Duration(pingIntervalMsec) * time.Millisecond)

		usernames := GetConnectedUsernames()
		go func() {
			for _, username := range usernames {
				toClientDispatcher(username, shared.BuildPingPacket())
			}
		}()

	}
	log.Panicln("Pinger returns!")
}
