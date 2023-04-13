package shared

import (
	"io"
	"log"
	"net"
)

import (
	"fmt"
	"os"
)

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
			log.Println("Broken pipe (got nil packet)... disconnecting and forcing cleanup.")

			// Will trigger cleanup in above layers
			incoming <- nil

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

	return packetData
}

func PacketSender(conn net.Conn, outgoing chan []byte) {

	for {
		wirePacket := <-outgoing
		_, err := conn.Write(wirePacket)
		if err != nil {
			log.Println("Error writing packet: ")
			conn.Close()
		}
	}
}
