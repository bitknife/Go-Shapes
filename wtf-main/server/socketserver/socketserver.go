package socketserver

import (
	"bitknife.se/wtf/server/core"
	"bitknife.se/wtf/shared"
	"log"
	"net"
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

	// Disconnect any existing
	playerLogin := shared.BytesToPacket(packageData).GetPlayerLogin()
	if core.HasChannel(playerLogin.Username) {
		// TODO: Could close the existing, or close this new one
		log.Println("Username", playerLogin.Username, "already connected, denying access!")
		conn.Close()
		return

		//core.UnRegisterClientChannels(playerLogin.Username)
	}

	accessGranted := core.AuthenticateClient(playerLogin)
	if accessGranted == false {
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

func Run(host string, port string) {

	listen, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		// Note: call to Fatal will do os.Exit(1).
		log.Fatal(err)
	}
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
