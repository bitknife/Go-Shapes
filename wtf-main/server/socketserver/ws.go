package socketserver

import (
	"bitknife.se/wtf/shared"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func RunWS(address string) {
	wc := NewWebsocketChannels(address)
	go wc.Run()
}

type WebsocketChannels struct {
	address  string
	upgrader websocket.Upgrader
}

func NewWebsocketChannels(
	address string) *WebsocketChannels {

	wc := WebsocketChannels{
		address:  address,
		upgrader: websocket.Upgrader{},
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
	// This is similar to handleConnection() of the TCP variant
	conn, err := wc.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()

	// This is the first package
	_, message, err := conn.ReadMessage()

	// Login and setup channels
	fromClient, toClient := HandleFirstPacket(&message)

	// TODO: Return correct HTTP status code upon invalid login?

	go shared.PacketReceiverWS(conn, fromClient)

	go shared.PacketSenderWS(conn, toClient)
}
