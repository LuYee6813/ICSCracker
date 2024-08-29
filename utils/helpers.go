package utils

import "log"

func IntToBytes(value int) []byte {
	return []byte{byte(value >> 8), byte(value)}
}

func HandleError(err error) {
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
