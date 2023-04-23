package game

import (
	"bitknife.se/wtf/shared"
	"log"
	"time"
)

const (
	// https://daposto.medium.com/game-networking-1-interval-and-ticks-b39bb51ccca9
	TICK_RATE      = 20
	STATS_INTERVAL = 3
)

var wtfGameGlobal WTFGame

func UserInputRunner(username string, userInputForGame chan *shared.Packet) {
	// TODO: Rethink, the global variable theWtfGame is a singleton

	for {
		packet := <-userInputForGame
		wtfGameGlobal.HandleUserInputPacket(username, packet)
	}
}

func Run(packetBroadCastChannel chan []*shared.Packet, wtfGame WTFGame) {

	// TODO: convert all this to a go struct methods (go "class")
	wtfGameGlobal = wtfGame

	ticTimeNano := time.Second / TICK_RATE

	// Server tick number
	tick := int64(0)

	aggregatedSleepTime := 0

	for {
		// Game loop
		start := time.Now()
		//-----------------------------------------------------------------

		// Update game logic
		wtfGame.Update()

		// Package and send game objects
		packets := shared.BuildGameObjectPackets(tick, wtfGame.GetGameObjects())

		packetBroadCastChannel <- packets

		// METRICS
		//-----------------------------------------------------------------
		tick = tick + 1
		if tick%(TICK_RATE*STATS_INTERVAL) == 0 {
			// Would be nice to collect average headroom
			// fmt.Println("Server tics: ", tics)

			// Calculate average headroom
			allPossibleSleepTime := ticTimeNano * TICK_RATE * STATS_INTERVAL
			sleepFraction := float32(aggregatedSleepTime) / float32(allPossibleSleepTime)
			log.Printf("Game-loop load: %.2f %%", 100-100*sleepFraction)
			aggregatedSleepTime = 0
		}

		// Calculate sleep time needed to keep FPS
		workTime := time.Since(start)
		sleepTime := ticTimeNano - workTime

		// For stats collection to see if we meet deadlines
		aggregatedSleepTime += int(sleepTime.Nanoseconds())

		// SLEEP
		//-----------------------------------------------------------------
		time.Sleep(sleepTime)
	}
}
