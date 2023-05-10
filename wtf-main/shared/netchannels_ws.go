package shared

import (
	"github.com/gorilla/websocket"
	"log"
)

func PacketReceiverWS(c *websocket.Conn, fromClient chan *[]byte) {
	for {
		_, message, err := c.ReadMessage()

		fromClient <- &message

		if err != nil {
			log.Println("read:", err)
			break
		}
	}
}

func PacketSenderWS(c *websocket.Conn, toClient chan *[]byte) {
	for {
		message := <-toClient
		err := c.WriteMessage(websocket.BinaryMessage, *message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
