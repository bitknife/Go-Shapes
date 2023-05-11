package shared

import (
	"context"
	"encoding/hex"
	"log"
	"net"
	"nhooyr.io/websocket"
	"sync/atomic"
)

import (
	"fmt"
	"os"
)

const (
	TCP_PORT        = "7777"
	WS_PORT         = "8888"
	WS_PACKETS_PATH = "/packets"
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

func ConnectClient(protocol string, host string,
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

	if protocol == "tcp" {
		// Connect to server
		tcpAddr := host + ":" + TCP_PORT
		log.Println("Connecting to TCP game server using", tcpAddr)

		conn, err := net.Dial(protocol, tcpAddr)
		if err != nil {
			fmt.Println("Error connecting:", err.Error())
			os.Exit(1)
		}

		// Separate send and receive routines works as conn is thread-safe
		go PacketReceiverTCP(conn, fromServer)
		go PacketSenderTCP(conn, toServer)

	} else if protocol == "udp" {
		// Same Dial as TCP , but would need other logic due to no connection

	} else if protocol == "websocket" {
		wsAddr := host + ":" + WS_PORT
		log.Println("Connecting to Websocket game server using", wsAddr)

		// u := url.URL{Scheme: "ws", Host: wsAddr, Path: WS_PACKETS_PATH}
		conn, _, err := websocket.Dial(context.TODO(), "ws://"+wsAddr+WS_PACKETS_PATH, nil)
		if err != nil {
			log.Fatal("dial:", err)
		}

		// WS connections are not thread safe, see comments in the method
		go WSPacketWorker(conn, fromServer, toServer)
	}
}

func PrintBuffer(buffer []byte) {
	encodedStr := hex.EncodeToString(buffer)
	fmt.Printf("%s\n", encodedStr)
}
