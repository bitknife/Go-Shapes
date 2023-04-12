package socketserver

import (
	"bitknife.se/wtf/server/core"
	"bitknife.se/wtf/shared"
	"fmt"
	"log"
	"net"
)

const (
	HOST = "localhost"
	PORT = "7777"
	TYPE = "tcp"
)

func handleConnection(conn net.Conn) {
	// shared.PlayerLogin{}

	/**
	First packet must be a login request.
	*/
	packageData := shared.ReceivePackageDataFromConnection(conn)
	if packageData == nil {
		log.Print("Got nil from connection, closing!")
		conn.Close()
		return
	}

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
	//go receivePacketsRoutine(conn, playerLogin, fromClient)
	go shared.PacketReceiver(conn, fromClient)

	// Main packet sender
	//go sendPackagesRoutine(conn, toClient)
	go shared.PacketSender(conn, toClient)
}

func makeAndRegisterChannels(playerLogin *shared.PlayerLogin) (chan []byte, chan []byte) {
	fromClient := make(chan []byte)
	toClient := make(chan []byte)

	// And register channels on the Dispatcher in the core layer
	core.RegisterToClientChannel(playerLogin.Username, toClient)
	core.RegisterFromClientChannel(playerLogin.Username, fromClient)

	return fromClient, toClient
}

func Run() {

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
