package core

import (
	"bitknife.se/wtf/shared"
	"fmt"
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

		fmt.Println("Number of clients = ", len(usernames))

		go func() {
			for _, username := range usernames {
				toClientDispatcher(username, shared.BuildPingPacket())
			}
		}()

	}
	log.Panicln("Pinger returns!")
}

func SendPacketsToUsername(username string, packets []*shared.Packet) {
	busy := 0
	busy += toClientDispatcherMulti(username, packets)

	if busy != 0 {
		// TODO: Stats this!
		// log.Println(">>> Busy channel for", username)
	}
}

func broadCastPackets(packets []*shared.Packet) {
	/**
	NOTE: This one is of course costly
	*/
	usernames := GetConnectedUsernames()
	for _, username := range usernames {
		// NOTE: We may need to flow control this one
		go SendPacketsToUsername(username, packets)
	}
}

func PacketBroadCaster(packetBroadCastChannel chan []*shared.Packet) {
	for {
		packets := <-packetBroadCastChannel
		broadCastPackets(packets)
	}
}
