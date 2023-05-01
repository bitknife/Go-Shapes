package shared

import (
	"math"
	"math/rand"
	"runtime"
	"strconv"
	"time"
)

func GetKeysFromStringInterfaceMap(theMap map[string]interface{}) []string {
	keys := make([]string, 0, len(theMap))
	for k := range theMap {
		keys = append(keys, k)
	}
	return keys
}

func RandInt(min int32, max int32) int32 {
	return rand.Int31n(max-min+1) + min
}

func RandName(base string) string {
	s := strconv.FormatInt(time.Now().UnixNano()+rand.Int63n(math.MaxInt32), 32)
	return base + "-" + s
}

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

func BurnCPU(num int64) {
	// For creating artificial CPU load, num should be in the million range
	for i := int64(0); i <= num; i++ {
		math.Sqrt(float64(i))
	}
}
