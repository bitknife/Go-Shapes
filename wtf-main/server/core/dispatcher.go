package core

import (
	"bitknife.se/wtf/shared"
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

func RegisterToClientChannel(username string, toClient chan []byte) {
	ToClientChannels.Set(username, toClient)
}

func RegisterFromClientChannel(username string, fromClient chan []byte) {
	FromClientChannels.Set(username, fromClient)

	// NOTE: This creates a goroutine for each client
	go fromClientHandler(username, fromClient)
}

func UnRegisterClientChannels(username string) {
	ToClientChannels.Pop(username)
	FromClientChannels.Pop(username)

	// NOTE: Not closing channels here, go developers says they are GCed.
}

func toClientDispatcher(username string, packet *shared.Packet) {
	// Look up the channel in the registry, and then send message
	toClientChannel, ok := ToClientChannels.Get(username)
	if ok {
		toClientChannel <- shared.PacketToBytes(packet)
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
			// Maybe not..
			// log.Println("fromClientHandler(): Unregistering client:", username)
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
