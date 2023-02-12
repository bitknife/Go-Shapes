package socketserver

import (
	"bitknife.se/core"
	"fmt"
	"io"
	"log"
	"net"
)

const (
	HOST = "localhost"
	PORT = "7777"
	TYPE = "tcp"
)

func handleConnection(conn net.Conn) {

	/**
	First packet must be a login request.
	*/
	packageData := receivePackageDataFromConnection(conn)
	playerLogin := core.AuthenticateClient(packageData)
	if playerLogin == nil {
		conn.Close()
		return
	}

	/**
	Client is authenticated, now we need to connect the client
	to the game. This is done using Channels that connects to
	the Dispatcher (middle layer), which then in turn connects to
	the game engine (upper layer).

	This separates the socket layer from the game layer completely.
	*/

	// Create and register the needed channels on the dispatcher
	fromClient, toClient := makeAndRegisterChannels(playerLogin)

	// Main packet receiver
	go receivePacketsRoutine(conn, playerLogin, fromClient)

	// Main packet sender
	go sendPackagesRoutine(conn, toClient)

}

func makeAndRegisterChannels(playerLogin *core.PlayerLogin) (chan core.DispatcherMessage, chan core.DispatcherMessage) {
	fromClient := make(chan core.DispatcherMessage)
	toClient := make(chan core.DispatcherMessage)

	// And register channels on the Dispatcher in the core layer
	core.RegisterToClientChannel(playerLogin.Username, toClient)
	core.RegisterFromClientChannel(playerLogin.Username, fromClient)

	return fromClient, toClient
}

func receivePacketsRoutine(conn net.Conn, playerLogin *core.PlayerLogin, fromClient chan core.DispatcherMessage) {
	for {
		packageData := receivePackageDataFromConnection(conn)

		if packageData == nil {
			// Communication error?
			log.Println("ERROR from receivePacketsRoutine()")

			// TODO: Improve, recover or disconnect/cleanup and let player reconnect etc.
			conn.Close()
			return
		}

		// Ok got a valid message, pass that to the dispatcher
		dm := core.DispatcherMessage{SourceID: playerLogin.Username, Data: packageData}
		fromClient <- dm
	}
}

func receivePackageDataFromConnection(conn net.Conn) []byte {
	/**
	Waits for the header and returns the type and []byte representing the package.
	*/
	// printReceivedBuffer(packetData, messageType)

	// Allocate header
	header := make([]byte, 1)

	// First read the two byte header
	_, err := io.ReadAtLeast(conn, header, 1)

	if err != nil {
		// Broken connection, client ugly shutdown etc.
		log.Print("Error reading from:", conn.RemoteAddr(), "reason was: ", err)
		log.Print("Closing!", conn)
		return nil
	}

	packageSize := header[0]
	// messageType := int(header[1])

	// Allocate for packet
	packetData := make([]byte, packageSize)

	// And read the packet
	_, err = io.ReadFull(conn, packetData)

	return packetData
}

func sendPackagesRoutine(conn net.Conn, toClient chan core.DispatcherMessage) {
	for {
		dm := <-toClient

		packet := dm.Data

		_, err := conn.Write(packet)
		if err != nil {
			log.Println("Error writing packet: ")
			return
		}
	}
}

func Start() {

	listen, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		// Note: call to Fatal will do os.Exit(1).
		log.Fatal(err)
	}
	fmt.Println("Server up on", HOST, ":", PORT)
	defer listen.Close()

	for {
		/*
			Accept() blocks
		*/
		log.Println("Waiting for next client to connect...")
		conn, err := listen.Accept()
		if err != nil {
			log.Println("Failed to Accept():", err)
		}

		log.Println("Client connected from", conn.RemoteAddr())
		go handleConnection(conn)
	}
}
