package shared

import (
	"github.com/gorilla/websocket"
	"log"
)

func PacketReceiverWS(conn *websocket.Conn, fromClient chan *[]byte) {
	for {
		_, message, err := conn.ReadMessage()

		fromClient <- &message

		if err != nil {
			log.Println("read:", err)
			break
		}
	}
}

func PacketSenderWS(conn *websocket.Conn, toClient chan *[]byte) {
	for {
		message := <-toClient
		err := conn.WriteMessage(websocket.BinaryMessage, *message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
