package shared

import (
	"github.com/gorilla/websocket"
)

func WSPacketWorker(conn *websocket.Conn, incoming chan *[]byte, outgoing chan *[]byte) {
	// NOTE We use the unbuffered channel and the select below
	// 	    to ensure we don't write to the connection while reading.
	//		This works, but still looks like a bit fishy.. :)
	unbufRecChan := make(chan *[]byte)
	go wsPacketsToChannel(conn, unbufRecChan)

	defer conn.Close()

	for {
		select {
		case packet := <-outgoing:

			err := conn.WriteMessage(websocket.BinaryMessage, *packet)

			if err != nil {
				// Client disconnect most likely, send nil to clean up
				incoming <- nil
				return
			}
		case packet := <-unbufRecChan:
			incoming <- packet
		}
	}
}

func wsPacketsToChannel(conn *websocket.Conn, unbufRecChan chan *[]byte) {

	for {
		_, message, err := conn.ReadMessage()

		unbufRecChan <- &message

		if err != nil {
			// log.Println("read:", err)
			unbufRecChan <- nil
			break
		}
	}
}
