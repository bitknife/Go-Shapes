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
var GameLoopActualFPS = new(float32)

type GameLoopMetrics struct {
	GameLoopSim       float32
	GameLoopSend      float32
	GameLoopSleep     float32
	GameLoopActualFPS float32
}

func GetGameLoopMetrics() *GameLoopMetrics {
	currentStats := GameLoopMetrics{
		GameLoopSim:       *GameLoopSim,
		GameLoopSend:      *GameLoopSend,
		GameLoopSleep:     *GameLoopSleep,
		GameLoopActualFPS: *GameLoopActualFPS,
	}
	return &currentStats
}

var doerGameGlobal shared.DoerGame

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

func Run(gameLoopFps int64, packetsForFrame chan []*shared.Packet, allComplete chan int, doerGame shared.DoerGame) {

	// TODO: convert all this to a go struct methods (go "class")
	doerGameGlobal = doerGame

	ticTime := FPSToDuration(int(gameLoopFps))

	log.Println("Game loop at", gameLoopFps, "FPS. Frame time:", ticTime)

	// Server tick number
	tick := int64(0)

	statsIntervalMsec := ticTime * STATS_INTERVAL
	statsDivideBy := float32(statsIntervalMsec) * float32(gameLoopFps)

	aggregatedSimTime := 0
	aggregatedSendTime := 0
	aggregatedSleepTime := 0

	// For measuring avg FPS between stats
	statsStart := time.Now()

	// Loop is mysteriously statically off by 1 ms, round-off errors or error in measurement
	fpsSleepAdj := time.Duration(1) * time.Millisecond

	for {
		// Game loop
		loopStartTime := time.Now()
		t1 := loopStartTime

		//--- Simulation ---------------------------------------------------

		// Update game logic
		doerGame.Update()

		t2 := time.Now()
		simTime := t2.Sub(t1)
		aggregatedSimTime += int(simTime.Nanoseconds())

		//--- Send ---------------------------------------------------

		// Package and send game objects
		packets := shared.BuildGameObjectPackets(tick, doerGame.GetGameObjects())

		// Broadcast packets, this will eat all packets
		packetsForFrame <- packets

		// TODO: Implement much smarter "send to clients" strategy! Ie. group by geoHash etc.

		// Wait for completion, we get an int here len(packets)
		<-allComplete

		t3 := time.Now()
		sendTime := t3.Sub(t2)
		aggregatedSendTime += int(sendTime.Nanoseconds())

		//--- Sleep ---------------------------------------------------

		sleepTime := ticTime - sendTime - simTime
		aggregatedSleepTime += int(sleepTime.Nanoseconds())

		// METRICS
		//-----------------------------------------------------------------
		tick = tick + 1
		if tick%(gameLoopFps*STATS_INTERVAL) == 0 {
			setGLLMetrics(statsDivideBy, aggregatedSimTime, aggregatedSendTime, aggregatedSleepTime)
			aggregatedSimTime, aggregatedSendTime, aggregatedSleepTime = 0, 0, 0

			// Measure actual FPS, if it misses by a lot, it means the load is too high
			// in either the simulation phase or the send phase.
			statsCycleTime := time.Since(statsStart)
			*GameLoopActualFPS = float32(gameLoopFps*STATS_INTERVAL) * float32(time.Second) / float32(statsCycleTime)
			statsStart = time.Now()

			// NOTE: Negative aggregatedSleepTime means we do not meet deadlines
			// 	     and simulation starts to slow down. Would need to handle this
			//		 gracefully, ideally upgrade hardware (or software) before this
			//		 happens.
			//		 Negative sleep can happen due to too high simulation load and
			//	     or too many clients connected or unable to parallelize enough.

			// Calculate sleep adjustment needed to better hit FPS
			// fpsDiff := float32(gameLoopFps) - *GameLoopActualFPS
			// log.Println("fpsDiff:", fpsDiff, "FPS")
			//newFpsSleepAdj := FloatFPSToDuration(float32(gameLoopFps) / fpsDiff)
			// log.Println("fpsSleepAdj:", FloatFPSToDuration(float32(gameLoopFps)/fpsDiff))
			// Slowly adjust it
			//fpsSleepAdj = (fpsSleepAdj + newFpsSleepAdj) / 2
		}

		// Calculate sleep time needed to keep FPS
		sleepDur := ticTime - time.Since(loopStartTime) - fpsSleepAdj
		time.Sleep(sleepDur)
	}
}

func setGLLMetrics(statsDivideBy float32, aggregatedSimTime int, aggregatedSendTime int, aggregatedSleepTime int) {
	*GameLoopSim = float32(aggregatedSimTime) / statsDivideBy
	*GameLoopSend = float32(aggregatedSendTime) / statsDivideBy
	*GameLoopSleep = float32(aggregatedSleepTime) / statsDivideBy
}

func FPSToDuration(fps int) time.Duration {
	return time.Duration(1000000000/fps) * time.Nanosecond
}

func FloatFPSToDuration(fps float32) time.Duration {
	return time.Duration(1000000000/fps) * time.Nanosecond
}
