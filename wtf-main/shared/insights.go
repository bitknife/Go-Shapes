package shared

import (
	"log"
	"runtime"
	"time"
)

const (
	STATS_INTERVAL_SEC = 10
)

func collectStats() map[string]interface{} {
	var stats map[string]interface{}
	stats = make(map[string]interface{})

	stats["numCpus"] = runtime.NumCPU()
	stats["numGoroutines"] = runtime.NumGoroutine()

	var memStats = runtime.MemStats{}
	runtime.ReadMemStats(&memStats)

	stats["heapAllocKB"] = memStats.HeapAlloc / (1000)
	stats["TotalAllocKB"] = memStats.TotalAlloc / (1000)

	return stats
}

func PrintStats() {
	stats := collectStats()
	log.Println("numCpus", stats["numCpus"])

	log.Println("numGoroutines", "|", "heapAllocKB", "|", "TotalAllocKB")
	log.Println("--------------------------------------------------------")
	for {
		stats := collectStats()

		log.Println(stats["numGoroutines"], "|", stats["heapAllocKB"], "kB", "|", stats["TotalAllocKB"], "kB")
		time.Sleep(STATS_INTERVAL_SEC * time.Second)
	}
}
