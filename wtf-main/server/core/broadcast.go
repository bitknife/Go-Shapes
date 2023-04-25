package core

import (
	"bitknife.se/wtf/shared"
	"sort"
	"sync/atomic"
)

var busyChannelDrops *int64 = new(int64)
var numberOfClients *int64 = new(int64)

type BroadcastStats struct {
	BusyChannelDrops int64
	NumberOfClients  int64
}

func GetBroadcastStats() *BroadcastStats {
	atomic.StoreInt64(numberOfClients, int64(len(GetConnectedUsernames())))

	currentStats := BroadcastStats{
		BusyChannelDrops: *busyChannelDrops,
		NumberOfClients:  *numberOfClients,
	}
	return &currentStats
}

func SendPacketsToUsername(username string, packets []*shared.Packet) {

	busy := ToClientDispatcherMulti(username, packets)

	if busy != 0 {
		/*
			When this happens a new call was made to the same client
			without having finished delivered the first one.

			Dabbled w. a retry. Better to drop I think.

			This is the bad thing with sending batches this way.
		*/
		// log.Println("Dropping", len(packets), "packets for", username)
		atomic.AddInt64(busyChannelDrops, int64(len(packets)))

		// TODO: Need to throttle server if this happens !
	}
}

func broadCastPackets(packets []*shared.Packet) {

	usernames := GetConnectedUsernames()

	// Important to send in same order to give each client more equal amount of time
	sort.Strings(usernames)

	for _, username := range usernames {
		/*
				NOTE, We opt to NOT "go" this. I guess an "async"
					  could be viable as well. but in that case the
					  broadcaster would/may need to be notified on
			          completion through a back channel.
		*/
		SendPacketsToUsername(username, packets)
	}
}

func PacketBroadCaster(packetBroadCastChannel chan []*shared.Packet, packetsSentChannel chan int) {

	for {
		packets := <-packetBroadCastChannel
		broadCastPackets(packets)

		// Just send the number of packages meant to send, not multiplied by
		// receivers
		packetsSentChannel <- len(packets)
	}
}
