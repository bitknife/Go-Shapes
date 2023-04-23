package shared

import (
	"io"
	"log"
	"net"
	"sync/atomic"
	"time"
)

import (
	"fmt"
	"os"
)

const (
	/*
		WRITE_TIMEOUT_MS If send time takes longer than this
		the send operation will be aborted and a packet loss
		will be noted.

		Value? 50-90% of game frame duration is a good start.

	*/
	WRITE_TIMEOUT_MS = 30
)

var bytesSent *int64 = new(int64)
var bytesReceived *int64 = new(int64)
var maxSendTime *int64 = new(int64)
var packetsLost *int64 = new(int64)

var packetsSent *int64 = new(int64)
var packetsReceived *int64 = new(int64)

type NetStats struct {
	BytesSent       int64
	BytesReceived   int64
	PacketsSent     int64
	PacketsReceived int64
	PacketsLost     int64

	MaxSendTime int64
}

func GetNetStats() *NetStats {
	currentStats := NetStats{
		BytesSent:       *bytesSent,
		BytesReceived:   *bytesReceived,
		PacketsSent:     *packetsSent,
		PacketsReceived: *packetsReceived,
		PacketsLost:     *packetsLost,
		MaxSendTime:     *maxSendTime,
	}
	return &currentStats
}

func ConnectClient(protocol string, host string, port string,
	fromServer chan []byte, toServer chan []byte) {

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
		// TODO: Implement websocket Dial, PacketReceive and PacketSend
	}
}

func PacketReceiverTCP(conn net.Conn, incoming chan []byte) {

	for {
		// Blocks
		packageData := ReceivePackageDataFromTCPConnection(conn)

		if packageData == nil {
			// Communication error, broken pipe etc
			// log.Println("Broken pipe (got nil packet)... disconnecting and forcing cleanup.")

			// Will trigger cleanup in above layers
			incoming <- nil

			// NOTE: Writer/Sender closes channels in Go!
			close(incoming)

			conn.Close()

			return
		}

		// "Nice" disconnect will be handeled by above layer

		// Ok got a valid message, pass that to the dispatcher
		incoming <- packageData

		// packet := BytesToPacket(packageData)
		// dm := core.DispatcherMessage{SourceID: playerLogin.Username, Packet: packet}
		// fromClient <- dm
	}
}

func ReceivePackageDataFromTCPConnection(conn net.Conn) []byte {
	/*
		Helper that waits for the header and returns the type and []byte representing the package.

		This can be used stand-alone
	*/

	// printReceivedBuffer(packetData, messageType)

	// Allocate header
	header := make([]byte, 1)

	// First read the two byte header
	_, err := io.ReadAtLeast(conn, header, 1)

	if err != nil {
		// Broken connection, client ugly shutdown etc.
		// log.Print("Error reading from:", conn.RemoteAddr(), "reason was: ", err)
		return nil
	}

	packageSize := header[0]

	// Allocate for packet
	packetData := make([]byte, packageSize)

	// And read the packet
	_, err = io.ReadFull(conn, packetData)

	// Stats
	atomic.AddInt64(packetsReceived, 1)
	atomic.AddInt64(bytesReceived, int64(len(packetData)+1))

	return packetData
}

func PacketSenderTCP(conn net.Conn, outgoing chan []byte) {

	for {
		// Wait for packets
		wirePacket := <-outgoing

		//------------------------------
		start := time.Now()

		if wirePacket == nil {
			log.Println("PacketSenderTCP(): Nil packet from channel. Aborting ")
			conn.Close()
			return
		}

		if WRITE_TIMEOUT_MS > 0 {
			conn.SetWriteDeadline(time.Now().Add(time.Duration(WRITE_TIMEOUT_MS) * time.Millisecond))
		}
		_, err := conn.Write(wirePacket)

		if err != nil {
			// NOTE: Packet-loss !
			atomic.AddInt64(packetsLost, 1)

			// Writing to closed socket
			// log.Println("PacketSenderTCP(): Error writing packet ")
			// conn.Close()
			return
		}

		// Stats
		atomic.AddInt64(packetsSent, 1)
		atomic.AddInt64(bytesSent, int64(len(wirePacket)))

		sendTime := time.Since(start)

		if int64(sendTime) > *maxSendTime {
			// fmt.Println("New Max send time", sendTime)
			atomic.StoreInt64(maxSendTime, int64(sendTime))
		}
	}
}
