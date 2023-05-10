package socketserver

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func RunWS(address string) {
	wc := NewWebsocketChannels(address)
	go wc.Run()
}

type WebsocketChannels struct {
	address     string
	upgrader    websocket.Upgrader
	messageType int
}

func NewWebsocketChannels(
	address string) *WebsocketChannels {

	wc := WebsocketChannels{
		address:     address,
		upgrader:    websocket.Upgrader{},
		messageType: websocket.BinaryMessage,
	}
	return &wc
}

func (wc *WebsocketChannels) Run() {
	// TODO: Set Write timeout

	// Register a handler function for given pattern
	http.HandleFunc("/packets", wc.packets)

	// NOTE: Blocks!
	log.Fatal(http.ListenAndServe(wc.address, nil))
}

func (wc *WebsocketChannels) packets(w http.ResponseWriter, r *http.Request) {
	// "CONNECT"-ish
	c, err := wc.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	_, message, err := c.ReadMessage()

	fromClient, toClient := HandleFirstPacket(&message)

	// Receiver loop
	go func() {
		for {
			_, message, err := c.ReadMessage()

			fromClient <- &message

			if err != nil {
				log.Println("read:", err)
				break
			}
		}
	}()

	// Sender loop
	go func() {
		for {
			message := <-toClient
			err = c.WriteMessage(wc.messageType, *message)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	}()
}
