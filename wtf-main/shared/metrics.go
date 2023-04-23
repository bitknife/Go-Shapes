package shared

import (
	"log"
	"runtime"
	"time"
)

func CollectGoStats() map[string]interface{} {
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

func CollectAndPrintMetricsRoutine(label string, interval_sec int) {
	/*
		For keeping an eye on server and game performance.

		Will be expanded to meet need. Just print to screen for now

	*/
	gS := CollectGoStats()

	log.Println("numCpus", gS["numCpus"])

	log.Println("numGoroutines | heapAlloc kB | totalAlloc kB | sent kB | Rec. kB | pS | pR | mST ms | nPL")
	log.Println("--", label, "------------------------------------------------------")
	for {
		gS := CollectGoStats()
		nS := GetStats()

		log.Println(gS["numGoroutines"], "|", gS["heapAllocKB"], "|", gS["TotalAllocKB"], "|",
			nS.BytesSent/1000, "|", nS.BytesReceived/1000, "|", nS.PacketsSent, "|", nS.PacketsReceived, "|",
			nS.MaxSendTime/1000000, "|", nS.PacketsLost)

		time.Sleep(time.Duration(interval_sec) * time.Second)
	}
}
