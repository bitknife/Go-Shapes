package game

import (
	"bitknife.se/wtf/shared"
	"log"
	"time"
)

const (
	// https://daposto.medium.com/game-networking-1-interval-and-ticks-b39bb51ccca9
	TICK_RATE      = 20
	STATS_INTERVAL = 1
)

type GameConfig struct {
	TickRate int32
}

// For calculating avg. send time
var GameLoopLoad = new(float32)

type GameLoopMetrics struct {
	GameLoopLoad float32
}

func GetGameLoopMetrics() *GameLoopMetrics {
	currentStats := GameLoopMetrics{
		GameLoopLoad: *GameLoopLoad,
	}
	return &currentStats
}

var wtfGameGlobal WTFGame

func UserInputRunner(username string, userInputForGame chan *shared.Packet) {
	// TODO: Rethink, the global variable theWtfGame is a singleton

	for {
		packet := <-userInputForGame
		if packet == nil {
			return
		}
		wtfGameGlobal.HandleUserInputPacket(username, packet)
	}
}

func Run(gameLoopFps int, packetBroadCastChannel chan []*shared.Packet, packetsSentChannel chan int, wtfGame WTFGame) {

	// TODO: convert all this to a go struct methods (go "class")
	wtfGameGlobal = wtfGame

	ticTime := time.Duration(1000/gameLoopFps) * time.Millisecond

	log.Println("Game loop at", gameLoopFps, "FPS. Frame time:", ticTime)

	// Server tick number
	tick := int64(0)

	aggregatedSleepTime := 0

	for {
		// Game loop
		loopStartTime := time.Now()

		//-----------------------------------------------------------------

		// Update game logic
		wtfGame.Update()

		// Package and send game objects
		packets := shared.BuildGameObjectPackets(tick, wtfGame.GetGameObjects())

		// Broadcast packets, this will eat all packets
		packetBroadCastChannel <- packets

		// TODO: Implement much smarter "send to clients" strategy! Ie. group by geoHash etc.

		// Wait for completion, we get an int here len(packets)
		<-packetsSentChannel

		// METRICS
		//-----------------------------------------------------------------
		tick = tick + 1
		if tick%(TICK_RATE*STATS_INTERVAL) == 0 {
			// Calculate average headroom
			allPossibleSleepTime := ticTime * TICK_RATE * STATS_INTERVAL
			sleepFraction := float32(aggregatedSleepTime) / float32(allPossibleSleepTime)
			*GameLoopLoad = 100 - 100*sleepFraction
			aggregatedSleepTime = 0
		}

		// Calculate sleep time needed to keep FPS
		loopEndTime := time.Since(loopStartTime)
		sleepTime := ticTime - loopEndTime

		// For stats collection to see if we meet deadlines when collecting/showing stats
		aggregatedSleepTime += int(sleepTime.Nanoseconds())

		// SLEEP
		//-----------------------------------------------------------------
		time.Sleep(sleepTime)
	}
}
