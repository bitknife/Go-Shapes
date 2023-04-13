package shared

import (
	"runtime"
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
