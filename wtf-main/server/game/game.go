package game

import (
	"bitknife.se/wtf/server/core"
	game "bitknife.se/wtf/server/game/test_worlds"
	"bitknife.se/wtf/shared"
	"log"
	"time"
)

const (
	// https://daposto.medium.com/game-networking-1-interval-and-ticks-b39bb51ccca9
	TICK_RATE      = 20
	STATS_INTERVAL = 3
)

func Run() {

	/**

	This is where the server-side game code would be implemented.

	*/
	gameObjects := make(map[string]*shared.GameObject)

	game.CreateDotWorld(gameObjects, 100, 300, 10)

	tic_time_nano := time.Second / TICK_RATE

	// Server tick number
	tick := int64(0)

	cum_sleep_time := 0
	for {
		// Game loop
		start := time.Now()
		// Code to measure

		/*
			START WORK
		*/

		// Update game logic
		game.ShakeDots(gameObjects, 4)

		// Build events to broadcast
		packets := buildgameObjectPackets(tick, gameObjects)
		core.BroadCastPackets(packets)

		/*
			END WORK
		*/

		tick = tick + 1
		if tick%(TICK_RATE*STATS_INTERVAL) == 0 {
			// Would be nice to collect average headroom
			// fmt.Println("Server tics: ", tics)

			// Calculate average headroom
			all_possible_sleep_time := tic_time_nano * TICK_RATE * STATS_INTERVAL
			sleep_fraction := float32(cum_sleep_time) / float32(all_possible_sleep_time)
			log.Printf("Game-loop load: %.2f %%", 100-100*sleep_fraction)
			cum_sleep_time = 0
		}

		// Calculate sleep time to keep FPS
		work_time := time.Since(start)
		sleep_time := tic_time_nano - work_time
		time.Sleep(sleep_time)

		// For stats collection to see if we meet deadlines
		cum_sleep_time += int(sleep_time.Nanoseconds())
	}
}

func buildgameObjectPackets(
	tick int64, gameObjects map[string]*shared.GameObject) []*shared.Packet {

	packets := make([]*shared.Packet, len(gameObjects))

	for _, gobj := range gameObjects {
		packet := shared.Packet{
			Payload: &shared.Packet_GameObject{GameObject: gobj},
		}
		packets = append(packets, &packet)
	}
	return packets
}
