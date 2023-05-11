package shared

import (
	"context"
	"nhooyr.io/websocket"
)

func WSPacketWorker(conn *websocket.Conn, incoming chan *[]byte, outgoing chan *[]byte) {
	// NOTE We use the unbuffered channel and the select below
	// 	    to ensure we don't write to the connection while reading.
	//		This works, but still looks like a bit fishy.. :)
	unbufRecChan := make(chan *[]byte)
	go wsPacketsToChannel(conn, unbufRecChan)

	for {
		select {
		case packet := <-outgoing:

			err := conn.Write(context.TODO(), websocket.MessageBinary, *packet)

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
		// NOTE: Would have preferred to to "canRead" on this, or select directly
		//		 dunno how to do that, but this model seem to work anyway.
		_, message, err := conn.Read(context.TODO())

		unbufRecChan <- &message

		if err != nil {
			// log.Println("read:", err)
			unbufRecChan <- nil
			break
		}
	}
}
