package core

import (
	"bitknife.se/wtf/server/game"
	"bitknife.se/wtf/shared"
	"fmt"
	cmap "github.com/orcaman/concurrent-map/v2"
)

/**
NOTE: The all mighty registry mapping a USERNAME to a CHANNEL

TODO: Rethink if Username is a good key or not, it has its merits (ie if connection lost
	  we could just re-attach the new connection and channel and move on). It would also
	  naturally not allow multiple connections using the same Username (ie. disconnect old,
	  or block the new).
*/

var ToClientChannels = cmap.New[chan []byte]()   // make(map[string]chan []byte)
var FromClientChannels = cmap.New[chan []byte]() // make(map[string]chan []byte)

func GetConnectedUsernames() []string {
	return ToClientChannels.Keys()
}

func HasChannel(username string) bool {
	if _, ok := ToClientChannels.Get(username); ok {
		return true
	}
	if _, ok := FromClientChannels.Get(username); ok {
		return true
	}
	return false
}

func InitClient(username string, toClient chan []byte, fromClient chan []byte) {
	ToClientChannels.Set(username, toClient)
	FromClientChannels.Set(username, fromClient)

	// NOTE: This creates a goroutine for each client
	go fromClientHandler(username, fromClient)
}

func UnRegisterClientChannels(username string) {
	fmt.Println("Unregistering", username)

	ToClientChannels.Pop(username)
	FromClientChannels.Pop(username)

	// NOTE: Not closing channels here, go developers says they are GCed.
}

func toClientDispatcher(username string, packet *shared.Packet) {
	// Look up the channel in the registry, and then send message
	toClientChannel, ok := ToClientChannels.Get(username)
	if ok {
		// NOTE: This blocks until lower layer is done!
		toClientChannel <- shared.PacketToBytes(packet)
	}
}

func fromClientHandler(username string, in chan []byte) {
	// This is OK, core knows of both game and socket layers,
	userInputForGame := make(chan *shared.Packet)
	go game.UserInputRunner(username, userInputForGame)

	for {
		buffer := <-in
		if buffer == nil {
			// Means the underlying layer will not send more packets, unregister and return
			UnRegisterClientChannels(username)
			userInputForGame <- nil
			close(userInputForGame)
			return
		}

		packet := shared.BytesToPacket(buffer)

		if packet == nil {
			// This is no payload from client, ignore for now or maybe disconect?

		} else {
			// OK got a packet, send it to
			userInputForGame <- packet
		}
	}
}
