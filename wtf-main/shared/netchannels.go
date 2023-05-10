package shared

import (
	"encoding/hex"
	"net"
	"sync/atomic"
)

import (
	"fmt"
	"os"
)

/*
WRITE_TIMEOUT_MS If send time takes longer than this
the send operation will be aborted and a packet loss
will be noted.

Value? 50-90% of game frame duration is a good start.
*/
var WriteTimeout = 5

// --- METRICS ---
var bytesSent *int64 = new(int64)
var bytesReceived *int64 = new(int64)
var minSendTimeMs *int64 = new(int64)
var maxSendTimeMs *int64 = new(int64)
var packetsLost *int64 = new(int64)

var packetsSent *int64 = new(int64)
var packetsReceived *int64 = new(int64)

// For calculating avg. send time
var totalSendTimeMs = new(int64)

type NetChannelsMetrics struct {
	BytesSent       int64
	BytesReceived   int64
	PacketsSent     int64
	PacketsReceived int64
	PacketsLost     int64
	MinSendTimeMs   int64
	MaxSendTimeMs   int64
	TotalSendTimeMs int64
}

func GetNetChannelsStats() *NetChannelsMetrics {
	currentStats := NetChannelsMetrics{
		BytesSent:       *bytesSent,
		BytesReceived:   *bytesReceived,
		PacketsSent:     *packetsSent,
		PacketsReceived: *packetsReceived,
		PacketsLost:     *packetsLost,
		MinSendTimeMs:   *minSendTimeMs,
		MaxSendTimeMs:   *maxSendTimeMs,
		TotalSendTimeMs: *totalSendTimeMs,
	}
	return &currentStats
}

func ConnectClient(protocol string, host string, port string,
	fromServer chan *[]byte, toServer chan *[]byte) {
	/*
		Purpose: Connect and set up the provided channels
		 		vs the underlying network protocol
	*/

	if minSendTimeMs == nil {
		// Just initialize to something big
		atomic.StoreInt64(minSendTimeMs, 10000000)
		atomic.StoreInt64(maxSendTimeMs, 0)
		atomic.StoreInt64(totalSendTimeMs, 0)
	}

	if protocol == "tcp" || protocol == "udp" {
		// Connect to server
		conn, err := net.Dial(protocol, host+":"+port)
		if err != nil {
			fmt.Println("Error connecting:", err.Error())
			os.Exit(1)
		}

		// Start send and receive routines
		go PacketReceiverTCP(conn, fromServer)
		go PacketSenderTCP(conn, toServer)

	} else if protocol == "websocket" {
	}
}

func PrintBuffer(buffer []byte) {
	encodedStr := hex.EncodeToString(buffer)
	fmt.Printf("%s\n", encodedStr)
}
