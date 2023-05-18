package shared

import (
	"io"
	"log"
	"net"
	"sync/atomic"
	"time"
)

func PacketReceiverTCP(conn net.Conn, incoming chan *[]byte) {

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

func ReceivePackageDataFromTCPConnection(conn net.Conn) *[]byte {
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

	return &packetData
}

func PacketSenderTCP(conn net.Conn, outgoing chan *[]byte) {

	for {
		// Wait for packets
		packet := <-outgoing

		//------------------------------
		start := time.Now()

		if packet == nil {
			// TODO: A bit harsh?
			log.Println("PacketSenderTCP(): Nil packet from channel. Closing conn. ")
			conn.Close()
			continue
		}

		// Append the length as a single byte
		if len(*packet) > 256 {
			panic("Packet size larger than 256 bytes!")
		}

		header := make([]byte, 1)
		header[0] = byte(len(*packet))
		wirePacket := append(header, *packet...)

		conn.SetWriteDeadline(time.Now().Add(time.Duration(WriteTimeout) * time.Millisecond))

		_, err := conn.Write(wirePacket)

		if err != nil {
			// NOTE: Packet-loss, note half a packet may have been sent
			//		 in which case the client would need to re-sync w. first-byte len.
			atomic.AddInt64(packetsLost, 1)
			continue
		}
		// log.Println("conn.Write() ok!")

		sendTimeMs := time.Since(start) / 1000000

		// Metrics
		atomic.AddInt64(packetsSent, 1)
		atomic.AddInt64(bytesSent, int64(len(wirePacket)))

		if int64(sendTimeMs) < *minSendTimeMs {
			// fmt.Println("New Max send time", sendTime)
			atomic.StoreInt64(minSendTimeMs, int64(sendTimeMs))
		}
		if int64(sendTimeMs) > *maxSendTimeMs {
			// fmt.Println("New Max send time", sendTime)
			atomic.StoreInt64(maxSendTimeMs, int64(sendTimeMs))
		}
		atomic.AddInt64(totalSendTimeMs, int64(sendTimeMs))
	}
}
