package socketserver

import (
	"bitknife.se/wtf/shared"
	"context"
	"log"
	"net/http"
	"nhooyr.io/websocket"
)

func RunWS(address string) {
	wc := NewWebsocketChannels(address)
	go wc.Run()
}

type WebsocketChannels struct {
	address string
}

func NewWebsocketChannels(
	address string) *WebsocketChannels {

	wc := WebsocketChannels{
		address: address,
	}
	return &wc
}

func (wc *WebsocketChannels) Run() {
	// TODO: Set Write timeout

	// Register a handler function for given pattern
	http.HandleFunc(shared.WS_PACKETS_PATH, wc.packetsHandler)

	// NOTE: Blocks!
	log.Fatal(http.ListenAndServe(wc.address, nil))
}

func (wc *WebsocketChannels) packetsHandler(w http.ResponseWriter, r *http.Request) {
	// This is similar to handleConnection() of the TCP variant
	ao := websocket.AcceptOptions{
		Subprotocols:         nil,
		InsecureSkipVerify:   true,
		OriginPatterns:       nil,
		CompressionMode:      0,
		CompressionThreshold: 0,
	}
	conn, err := websocket.Accept(w, r, &ao)
	if err != nil {
		// ...
	}

	// This is the first package
	_, message, err := conn.Read(context.TODO())

	// Login and setup channels
	fromClient, toClient := HandleFirstPacket(&message)

	if fromClient == nil {
		conn.Close(websocket.StatusAbnormalClosure, "First packet failed")
	}
	// TODO: Return correct HTTP status code upon invalid login?

	// NOTE: We do this slighly different than TCP due to the nature of the Websocket connection
	go shared.WSPacketWorker(conn, fromClient, toClient)

}
