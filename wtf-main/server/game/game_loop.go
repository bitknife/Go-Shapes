package game

import (
	"bitknife.se/wtf/shared"
	"log"
	"time"
)

const (
	// https://daposto.medium.com/game-networking-1-interval-and-ticks-b39bb51ccca9
	// TICK_RATE      = 20
	STATS_INTERVAL = 2
)

type GameConfig struct {
	TickRate int32
}

// For calculating avg. send time
var GameLoopSim = new(float32)
var GameLoopSend = new(float32)
var GameLoopSleep = new(float32)

type GameLoopMetrics struct {
	GameLoopSim   float32
	GameLoopSend  float32
	GameLoopSleep float32
}

func GetGameLoopMetrics() *GameLoopMetrics {
	currentStats := GameLoopMetrics{
		GameLoopSim:   *GameLoopSim,
		GameLoopSend:  *GameLoopSend,
		GameLoopSleep: *GameLoopSleep,
	}
	return &currentStats
}

var doerGameGlobal DoerGame

func UserInputRunner(username string, userInputForGame chan *shared.Packet) {
	// TODO: Rethink, the global variable theWtfGame is a singleton

	for {
		packet := <-userInputForGame
		if packet == nil {
			return
		}
		doerGameGlobal.HandleUserInputPacket(username, packet)
	}
}

func Run(gameLoopFps int64, packetBroadCastChannel chan []*shared.Packet, packetsSentChannel chan int, doerGame DoerGame) {

	// TODO: convert all this to a go struct methods (go "class")
	doerGameGlobal = doerGame

	ticTime := time.Duration(1000/gameLoopFps) * time.Millisecond

	log.Println("Game loop at", gameLoopFps, "FPS. Frame time:", ticTime)

	// Server tick number
	tick := int64(0)

	statsIntervalMsec := ticTime * STATS_INTERVAL
	statsDivideBy := float32(statsIntervalMsec) * float32(gameLoopFps)

	aggregatedSimTime := 0
	aggregatedSendTime := 0
	aggregatedSleepTime := 0

	for {
		// Game loop
		loopStartTime := time.Now()
		t1 := loopStartTime

		//--- Simulation ---------------------------------------------------

		// Update game logic
		doerGame.Update()

		// Package and send game objects
		packets := shared.BuildGameObjectPackets(tick, doerGame.GetGameObjects())

		t2 := time.Now()
		simTime := t2.Sub(t1)
		aggregatedSimTime += int(simTime.Nanoseconds())

		//--- Send ---------------------------------------------------

		// Broadcast packets, this will eat all packets
		packetBroadCastChannel <- packets

		// TODO: Implement much smarter "send to clients" strategy! Ie. group by geoHash etc.

		// Wait for completion, we get an int here len(packets)
		<-packetsSentChannel

		t3 := time.Now()
		sendTime := t3.Sub(t2)
		aggregatedSendTime += int(sendTime.Nanoseconds())

		//--- Sleep ---------------------------------------------------

		// Calculate sleep time needed to keep FPS
		sleepDur := ticTime - time.Since(loopStartTime)
		time.Sleep(sleepDur)

		sleepTime := ticTime - sendTime - simTime
		aggregatedSleepTime += int(sleepTime.Nanoseconds())

		// METRICS
		//-----------------------------------------------------------------
		tick = tick + 1
		if tick%(gameLoopFps*STATS_INTERVAL) == 0 {
			setGLLMetrics(statsDivideBy, aggregatedSimTime, aggregatedSendTime, aggregatedSleepTime)
			aggregatedSimTime, aggregatedSendTime, aggregatedSleepTime = 0, 0, 0
		}
	}
}

func setGLLMetrics(statsDivideBy float32, aggregatedSimTime int, aggregatedSendTime int, aggregatedSleepTime int) {
	*GameLoopSim = float32(aggregatedSimTime) / statsDivideBy
	*GameLoopSend = float32(aggregatedSendTime) / statsDivideBy
	*GameLoopSleep = float32(aggregatedSleepTime) / statsDivideBy
}
