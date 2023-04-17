package shared

import (
	"math/rand"
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

func RandInt(min int, max int) int {
	return rand.Intn(max-min) + min
}

func RandName(base string) string {
	return base + "-" + strconv.Itoa(int(time.Now().UnixNano()))
}
