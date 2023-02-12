package socketserver

import (
	"encoding/hex"
	"fmt"
)

func printBuffer(buffer []byte) {
	encodedStr := hex.EncodeToString(buffer)
	fmt.Printf("%s\n", encodedStr)
}

func toGameMessage(buffer []byte, messageType byte) {
}
