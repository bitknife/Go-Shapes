package shared

import (
	"io"
	"log"
	"net"
	"sync/atomic"
)

import (
	"fmt"
	"os"
)

var bytesSent *int64 = new(int64)
var bytesReceived *int64 = new(int64)
var packetsSent *int64 = new(int64)
var packetsReceived *int64 = new(int64)

type NetStats struct {
	BytesSent       int64
	BytesReceived   int64
	PacketsSent     int64
	PacketsReceived int64
}

func GetStats() *NetStats {
	currentStats := NetStats{
		BytesSent:       *bytesSent,
		BytesReceived:   *bytesReceived,
		PacketsSent:     *packetsSent,
		PacketsReceived: *packetsReceived,
	}
	return &currentStats
}

func Connect(host string, port string) net.Conn {

	// Connect to server
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}
	return conn
}

func PacketReceiver(conn net.Conn, incoming chan []byte) {

	for {
		// Blocks
		packageData := ReceivePackageDataFromConnection(conn)

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

func ReceivePackageDataFromConnection(conn net.Conn) []byte {
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

func PacketSender(conn net.Conn, outgoing chan []byte) {

	for {
		wirePacket := <-outgoing
		if wirePacket == nil {
			log.Println("PacketSender(): Nil packet from channel. Aborting ")
			conn.Close()
			return
		}
		_, err := conn.Write(wirePacket)
		if err != nil {
			// Writing to closed socket
			// log.Println("PacketSender(): Error writing packet ")
			conn.Close()
			return
		}

		// Stats
		atomic.AddInt64(packetsSent, 1)
		atomic.AddInt64(bytesSent, int64(len(wirePacket)))
	}
}
