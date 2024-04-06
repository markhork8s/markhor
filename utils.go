package main

// Convert map[string]string to map[string][]byte
func convertMapStringToMapByte(inputMap map[string]string) map[string][]byte {
	outputMap := make(map[string][]byte)
	for key, value := range inputMap {
		outputMap[key] = []byte(value)
	}
	return outputMap
}

// Convert map[string][]byte to map[string]string
func convertMapByteToMapString(inputMap map[string][]byte) map[string]string {
	outputMap := make(map[string]string)
	for key, value := range inputMap {
		outputMap[key] = string(value)
	}
	return outputMap
}
