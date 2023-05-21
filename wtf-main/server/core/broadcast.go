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

func SendPacketsToUsername(username string, bPackets []*[]byte) {
	/*
		NOTE: This function is called from each frame of the game loop!
			  Meaning, this must finish before the next frame is ready to send
			  else the client will experience loss of frames from the server.
	*/

	busy := ToClientDispatcherMultiBytes(username, bPackets)

	if busy {
		/*
			When this happens a new frame was sent to the _same_ client
			without having finished delivered the first one.
		*/
		atomic.AddInt64(busyChannelDrops, int64(len(bPackets)))
	}
}

func broadCastPackets(packets []*shared.Packet) {

	usernames := GetConnectedUsernames()

	if len(usernames) == 0 {
		return
	}

	// Important to send in same order to give each client more equal amount of time
	sort.Strings(usernames)

	// OPTIMIZATION: Do this early
	bytePackets := make([]*[]byte, len(packets))
	for n, packet := range packets {
		bytePackets[n] = shared.PacketToBytes(packet)
	}

	for _, username := range usernames {
		// NOTE: We "go" here to do this concurrently and parallel if multiple CPUs
		//       as this is quite work intensive.
		go SendPacketsToUsername(username, bytePackets)
	}
}

/*
var GameLoopSend = new(float32)
GameLoopSend      float32
GameLoopSend:      *GameLoopSend,
*GameLoopSend = float32(aggregatedSendTime) / statsDivideBy
*/
func PacketBroadCaster(packetBroadCastChannel chan []*shared.Packet) {

	for {
		packets := <-packetBroadCastChannel
		broadCastPackets(packets)
	}
}
