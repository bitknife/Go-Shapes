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

func RandInt(min int32, max int32) int32 {
	return rand.Int31n(max-min+1) + min
}

func RandName(base string) string {
	return base + "-" + strconv.Itoa(int(time.Now().UnixNano()))
}
