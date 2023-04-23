package main

import (
	"bitknife.se/wtf/shared"
	"log"
	"syscall"
)

func SetUpNetworking(protocol string, host string, port string, username string, password string) (chan []byte, chan []byte) {

	fromServer := make(chan []byte)
	toServer := make(chan []byte)

	// Connects
	log.Println("Connecting to game server at", host+":"+port, "as", username)

	shared.ConnectClient(protocol, host, port, fromServer, toServer)

	pPacket := shared.BuildLoginPacket(username, password)
	wirePacket := shared.PacketToBytes(pPacket)
	toServer <- wirePacket
	log.Println("Login successful!")

	return fromServer, toServer
}

func HandlePacketsFromServer(
	fromServer chan []byte,
	toServer chan []byte,
	gamePacketsToUpperLayers chan *shared.Packet) {
	for {
		receivedData := <-fromServer

		if receivedData == nil {
			// This is server disconnecting, raise SIGINT to trigger exit handler
			syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		}
		packet := shared.BytesToPacket(receivedData)

		/*
			Unpack and handle or re-route "Low-level" packets directly.
				(At time of writing, this is Ping only).

			Rest is routed upwards to game layer.

		*/
		if packet.GetPing() != nil {
			sent := packet.GetPing().Sent
			log.Println("Got Ping from server:", sent)
			pP := shared.BuildPingPacket()
			toServer <- shared.PacketToBytes(pP)

		} else {
			// Onwards an upwards!
			gamePacketsToUpperLayers <- packet
		}
	}
}
