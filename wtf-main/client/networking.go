package main

import (
	"bitknife.se/wtf/shared"
	"log"
)

func SetUpNetworking(
	protocol string,
	host string,
	tcpPort string,
	wsPort string,
	username string,
	password string) (chan *[]byte, chan *[]byte) {

	fromServer := make(chan *[]byte)
	toServer := make(chan *[]byte)

	// Connects
	shared.ConnectClient(protocol, host, tcpPort, wsPort, fromServer, toServer)

	// Login
	pPacket := shared.BuildLoginPacket(username, password)
	wirePacket := shared.PacketToBytes(pPacket)
	toServer <- wirePacket
	log.Println("Login successful!")

	return fromServer, toServer
}

func DeliverPacketsToServer(
	toServer chan *[]byte,
	updatesToServer chan *shared.Packet) {

	for {
		packet := <-updatesToServer
		toServer <- shared.PacketToBytes(packet)
	}
}

func ReceivePacketsFromServer(
	fromServer chan *[]byte,
	packetsToView chan *shared.Packet) {

	for {
		receivedData := <-fromServer

		if receivedData == nil {
			// This is server disconnecting, send nil to signal this
			packetsToView <- nil
			return
		}
		packet := shared.BytesToPacket(receivedData)

		/*
			Unpack and handle or re-route "Low-level" packets directly.
				(At time of writing, this is Ping only).

			Rest is routed upwards to game layer.

		*/
		// TODO: Move elsewhere? Should be ping initiated from client also
		if packet.GetPing() != nil {
			sent := packet.GetPing().Sent
			log.Println("Got Ping from server:", sent)

			// TODO: Move to better place
			// pP := shared.BuildPingPacket()
			// toServer <- shared.PacketToBytes(pP)

		} else {

			packetsToView <- packet
		}
	}
}
