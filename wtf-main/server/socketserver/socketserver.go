package socketserver

import (
	"bitknife.se/wtf/server/core"
	"bitknife.se/wtf/shared"
	"log"
	"net"
	"time"
)

func handleConnection(conn net.Conn) {
	// shared.PlayerLogin{}

	/**
	First packet must be a login request.
	*/
	packageData := shared.ReceivePackageDataFromTCPConnection(conn)
	if packageData == nil {
		log.Print("Got nil from connection, closing!")
		conn.Close()
		return
	}

	// Will initiate client de-registration of old user if connected
	playerLogin := shared.BytesToPacket(packageData).GetPlayerLogin()
	if _, ok := core.ToClientChannelsRegistry.Get(playerLogin.Username); ok {
		frC, _ := core.FromClientChannels.Get(playerLogin.Username)
		frC <- nil

		// TODO: Fix race?
		//       Need to wait for old connection to close before moving on
		time.Sleep(time.Duration(100) * time.Millisecond)
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

	This separates the socket layer from the game layers.
	*/

	// Create and register the needed channels on the dispatcher
	fromClient, toClient := makeAndRegisterChannels(playerLogin)

	// Main packet receiver
	//go receivePacketsRoutine(conn, playerLogin, fromClient)
	go shared.PacketReceiverTCP(conn, fromClient)

	// Main packet sender
	//go sendPackagesRoutine(conn, toClient)
	go shared.PacketSenderTCP(conn, toClient)

	log.Println("User", playerLogin.Username, "accepted and setup!")
}

func makeAndRegisterChannels(playerLogin *shared.PlayerLogin) (chan *[]byte, chan *[]byte) {
	// IDEA: Is this where we want to create them?
	fromClient := make(chan *[]byte)
	toClient := make(chan *[]byte)

	// And register channels on the Dispatcher in the core layer
	core.InitClient(playerLogin.Username, toClient, fromClient)

	return fromClient, toClient
}

func Run(host string, port string) {

	// TODO: Implement udp and websocket
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
		// log.Println("Waiting for next client to connect...")
		conn, err := listen.Accept()
		if err != nil {
			log.Println("Failed to Accept():", err)
		}

		// log.Println("Client connected from", conn.RemoteAddr())
		go handleConnection(conn)
	}
}
