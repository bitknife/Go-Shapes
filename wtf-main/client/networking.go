package main

import (
	"bitknife.se/wtf/shared"
	"log"
)

func SetUpNetworking(username string, password string) (chan []byte, chan []byte) {

	// Connects
	conn := shared.Connect(HOST, PORT)

	fromServer := make(chan []byte)
	toServer := make(chan []byte)

	// Login using the connection directly
	pPacket := shared.BuildLoginPacket(username, password)
	wirePacket := shared.PacketToBytes(pPacket)
	_, err := conn.Write(wirePacket)
	if err != nil {
		log.Println("Error writing packet: ")
		conn.Close()
	}

	log.Println("Login successful!")

	// Start send and receive routines
	go shared.PacketSender(conn, toServer)
	go shared.PacketReceiver(conn, fromServer)

	return fromServer, toServer
}

func HandlePacketsFromServer(fromServer chan []byte, toServer chan []byte) {
	for {
		receivedData := <-fromServer
		packet := shared.BytesToPacket(receivedData)

		if packet.GetPing() != nil {
			sent := packet.GetPing().Sent
			log.Println("Got ", sent)
			pP := shared.BuildPingPacket()
			toServer <- shared.PacketToBytes(pP)
		} else if packet.GetPlayerLogin() != nil {
			// TODO for each packet type and
		}
	}
}
