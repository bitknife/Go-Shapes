package core

import (
	"bitknife.se/wtf/shared"
	"golang.org/x/sync/syncmap"
	"log"
)

/**
NOTE: The all mighty registry mapping a USERNAME to a CHANNEL

TODO: Rethink if Username is a good key or not, it has its merits (ie if connection lost
	  we could just re-attach the new connection and channel and move on). It would also
	  naturally not allow multiple connections using the same Username (ie. disconnect old,
	  or block the new).
*/

var ToClientChannels = syncmap.Map{}   // make(map[string]chan []byte)
var FromClientChannels = syncmap.Map{} // make(map[string]chan []byte)

func GetConnectedUsernames() []string {

	var keys []string

	ToClientChannels.Range(func(key, value interface{}) bool {
		keys = append(keys, key.(string))
		return true
	})
	return keys
}

func HasChannel(username string) bool {
	if _, ok := ToClientChannels.Load(username); ok {
		return true
	}
	if _, ok := FromClientChannels.Load(username); ok {
		return true
	}
	return false
}

func RegisterToClientChannel(username string, toClient chan []byte) {
	ToClientChannels.Store(username, toClient)
}

func RegisterFromClientChannel(username string, fromClient chan []byte) {
	FromClientChannels.Store(username, fromClient)

	// NOTE: This creates a goroutine for each client
	go fromClientHandler(username, fromClient)
}

func UnRegisterClientChannels(username string) {
	ToClientChannels.Delete(username)
	FromClientChannels.Delete(username)
}

func toClientDispatcher(username string, packet *shared.Packet) {
	// Look up the channel in the registry, and then send message
	toClientChannel, ok := ToClientChannels.Load(username)
	if ok {
		tc := toClientChannel.(chan []byte)
		tc <- shared.PacketToBytes(packet)
	}
}

func SendPacketsToUsername(username string, packets []*shared.Packet) {
	for _, packet := range packets {
		toClientDispatcher(username, packet)
	}
}

func BroadCastPackets(packets []*shared.Packet) {
	/**
	NOTE: Costly!
	*/
	usernames := GetConnectedUsernames()
	for _, username := range usernames {
		// Go routine for each user as they all have their own socket
		go SendPacketsToUsername(username, packets)
	}
}

func fromClientHandler(username string, in chan []byte) {
	for {
		buffer := <-in
		if buffer == nil {
			log.Println("fromClientHandler(): Unregistering client:", username)
			UnRegisterClientChannels(username)
			return
		}

		packet := shared.BytesToPacket(buffer)

		if packet == nil {
			// TODO: handle client inputs

			/*
				Should send username and payload (or packet) to Game etc.
			*/
		}
		// log.Println("Dispatcher got", packet.GetTheMessage(), "from:", username)

		/**
		A possible good pattern would be to publish client events on a topic w. Candidate keys could be:
			- Event type - that way we could ensure listeners are aware of the types coming.
		*/
	}
}
