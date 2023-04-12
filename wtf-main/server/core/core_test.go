package core

import (
	"bitknife.se/wtf/shared"
	"encoding/hex"
	"fmt"
	"log"
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestPackets(t *testing.T) {
	playerLogin := shared.PlayerLogin{
		Username: "Tester",
		Password: "s33cret!",
	}
	out, err := proto.Marshal(&playerLogin)

	if err != nil {
		log.Fatalln("Failed to Marshal PlayerLogin object:", err)
	}
	encodedStr := hex.EncodeToString(out)
	fmt.Println("Encoded to: ", encodedStr)

	playerLoginCopy := shared.PlayerLogin{}
	err = proto.Unmarshal(out, &playerLoginCopy)

	fmt.Println("Username: ", playerLoginCopy.Username)
	fmt.Println("Password: ", playerLoginCopy.Password)
}

func TestTypeSwitch(t *testing.T) {
	// Given: create a packet

	// When: Switch on the type

	// Then: Ensure correct type

}
