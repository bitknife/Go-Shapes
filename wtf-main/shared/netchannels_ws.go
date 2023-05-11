package shared

import (
	"github.com/gorilla/websocket"
	"log"
)

func PacketReceiverWS(conn *websocket.Conn, incoming chan *[]byte) {
	for {
		_, message, err := conn.ReadMessage()

		log.Println("Got", len(message), "bytes from websocket")

		incoming <- &message

		if err != nil {
			log.Println("read:", err)
			break
		}
	}
}

func PacketSenderWS(conn *websocket.Conn, outgoing chan *[]byte) {
	for {
		message := <-outgoing

		err := conn.WriteMessage(websocket.BinaryMessage, *message)

		log.Println("Wrote", len(*message), "bytes to websocket")

		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func wsPacketsToChannel(conn *websocket.Conn, receiveBuf chan *[]byte) {

	for {
		_, message, err := conn.ReadMessage()

		log.Println("Got", len(message), "bytes from websocket")

		receiveBuf <- &message

		if err != nil {
			log.Println("read:", err)
			break
		}
	}
}
func WSPacketWorker(conn *websocket.Conn, incoming chan *[]byte, outgoing chan *[]byte) {
	// We use the unbuffered channel and the select below
	// to ensure we don't write to the connection while reading
	receiveBuf := make(chan *[]byte)
	go wsPacketsToChannel(conn, receiveBuf)

	for {
		select {
		case packet := <-outgoing:

			err := conn.WriteMessage(websocket.BinaryMessage, *packet)

			log.Println("Wrote", len(*packet), "bytes to websocket")

			if err != nil {
				log.Println("write:", err)
				break
			}
		case packet := <-receiveBuf:
			incoming <- packet
		}
	}
}
