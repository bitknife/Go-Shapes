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
	messageType, message := receivePackageFromConnection(conn)
	playerLogin := core.AuthenticateClient(messageType, message)
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
		messageType, messageData := receivePackageFromConnection(conn)

		if messageType == 0 {
			// Communication error?
			log.Println("ERROR from receivePackageFromConnection()!")

			// TODO: Improve
			conn.Close()
			return
		}

		// Ok got a valid message, pass that to the dispatcher
		dm := core.DispatcherMessage{SourceID: playerLogin.Username, Type: messageType, Data: messageData}
		fromClient <- dm
	}
}

func receivePackageFromConnection(conn net.Conn) (int, []byte) {
	/**
	Waits for the header and returns the type and []byte representing the package.
	*/
	// printReceivedBuffer(messageData, messageType)

	// Allocate header
	header := make([]byte, 2)

	// First read the two byte header
	_, err := io.ReadAtLeast(conn, header, 2)

	if err != nil {
		// Broken connection, client ugly shutdown etc.
		log.Print("Error reading from:", conn.RemoteAddr(), "reason was: ", err)
		log.Print("Closing!", conn)
		return 0, nil
	}

	messageSize := header[0]
	messageType := int(header[1])

	// Allocate for packet
	messageData := make([]byte, messageSize)

	// And read the packet
	_, err = io.ReadFull(conn, messageData)

	return messageType, messageData
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
