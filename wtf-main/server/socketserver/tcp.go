package socketserver

import (
	"bitknife.se/wtf/shared"
	"log"
	"net"
)

func ServeTCP(address string) {
	// TODO: Refactor (TCP version?) to align more with WS variant
	//		 containing all data on a struct and then "Run()"

	// TODO: Implement udp and websocket
	listen, err := net.Listen("tcp", address)

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

	fromClient, toClient := HandleFirstPacket(packageData)

	// Main packet receiver
	go shared.PacketReceiverTCP(conn, fromClient)

	// Main packet sender
	go shared.PacketSenderTCP(conn, toClient)
}
