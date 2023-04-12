package core

import (
	"bitknife.se/wtf/shared"
	"log"
)

type DispatcherMessage struct {
	// SourceID is for now the Username
	SourceID string
	Packet   *shared.Packet
}

/**
NOTE: The all mighty registry mapping a USERNAME to a CHANNEL

TODO: Rethink if Username is a good key or not, it has its merits (ie if connection lost
	  we could just re-attach the new connection and channel and move on). It would also
	  naturally not allow multiple connections using the same Username (ie. disconnect old,
	  or block the new).
*/

var ToClientChannels = make(map[string]chan DispatcherMessage)
var FromClientChannels = make(map[string]chan DispatcherMessage)

func GetConnectedUsernames() []string {
	return getAllKeysFromMap(ToClientChannels)
}

func getAllKeysFromMap(theMap map[string]chan DispatcherMessage) []string {
	keys := make([]string, 0, len(theMap))
	for k := range theMap {
		keys = append(keys, k)
	}
	return keys
}

func RegisterToClientChannel(username string, toClient chan DispatcherMessage) {
	ToClientChannels[username] = toClient
}

func RegisterFromClientChannel(username string, fromClient chan DispatcherMessage) {
	FromClientChannels[username] = fromClient

	// NOTE: This creates a goroutine for each client
	go fromClientHandler(fromClient)
}

func toClientDispatcher(message DispatcherMessage) {
	// Look up the channel in the registry, and then send message
	toClientChannel := ToClientChannels[message.SourceID]
	toClientChannel <- message
}

func fromClientHandler(in chan DispatcherMessage) {
	for {
		// NOTE: This blocks the routine until next message arrives
		dm := <-in
		log.Println("Dispatcher got", dm.Packet.GetTheMessage(), "from:", dm.SourceID)

		/**
		A possible good pattern would be to publish client events on a topic w. Candidate keys could be:
			- Event type - that way we could ensure listeners are aware of the types coming.
		*/
	}
}
