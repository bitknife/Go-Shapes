package shared

func GetKeysFromStringInterfaceMap(theMap map[string]interface{}) []string {
	keys := make([]string, 0, len(theMap))
	for k := range theMap {
		keys = append(keys, k)
	}
	return keys
}
