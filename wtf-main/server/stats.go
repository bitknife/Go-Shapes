package main

import (
	"bitknife.se/wtf/shared"
	"log"
	"time"
)

func PrintStats(interval_sec int) {
	stats := shared.CollectStats()
	log.Println("numCpus", stats["numCpus"])

	log.Println("numGoroutines", "|", "heapAllocKB", "|", "TotalAllocKB")
	log.Println("--------------------------------------------------------")
	for {
		stats := shared.CollectStats()

		log.Println(stats["numGoroutines"], "|", stats["heapAllocKB"], "kB", "|", stats["TotalAllocKB"], "kB")
		time.Sleep(time.Duration(interval_sec) * time.Second)
	}
}
