package core

import (
	"bitknife.se/wtf/shared"
	"log"
	"time"
)

const (
	PING_INTERVAL_SEC = 5
)

func Run() {
	/**
	Core routine(s) are responsible for all things _NOT_ related to "gamey" stuff.
		- Housekeeping
		- Metrics collection/delivery
		- ...
	*/

	// Just an example for now
	go broadCastPing()
}

func broadCastPing() {

	for {
		time.Sleep(PING_INTERVAL_SEC * time.Second)

		usernames := GetConnectedUsernames()
		go func() {
			for _, username := range usernames {
				toClientDispatcher(username, shared.BuildPingPacket())
			}
		}()

	}
	log.Panicln("Pinger returns!")
}
