package main

import (
	"bitknife.se/wtf/client/ebiten"
	"bitknife.se/wtf/shared"
	"log"
)

const (
	HOST = "localhost"
	PORT = "7777"
)

func setUpNetworking(username string, password string) {
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

	go func() {
		for {
			receivedData := <-fromServer
			sent := shared.BytesToPacket(receivedData).GetPing().Sent
			log.Println("Got ", sent)
		}
	}()
}

func main() {

	setUpNetworking("goClient", "welcome")

	// Starts the UI, this blocks
	ebiten.RunEbitenApplication()

}
