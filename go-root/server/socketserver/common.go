package socketserver

import (
	"bitknife.se/wtf/server/core"
	"bitknife.se/wtf/shared"
	"log"
	"time"
)

func HandleFirstPacket(packageData *[]byte) (chan *[]byte, chan *[]byte, *shared.PlayerLogin) {

	// Will initiate client de-registration of old user if connected
	packet := shared.BytesToPacket(packageData)

	if packet == nil {
		log.Println("Invalid first packet.")
		return nil, nil, nil
	}

	playerLogin := packet.GetPlayerLogin()

	if _, ok := core.ToClientChannelsRegistry.Get(playerLogin.Username); ok {
		log.Println("User", playerLogin.Username, "already logged in.")
		frC, _ := core.FromClientChannels.Get(playerLogin.Username)
		frC <- nil

		// TODO: Fix race?
		//       Need to wait for old connection to close before moving on
		time.Sleep(time.Duration(100) * time.Millisecond)
	}

	accessGranted := core.AuthenticateClient(playerLogin)
	if accessGranted == false {
		log.Println("Access denied for", playerLogin.Username)
		return nil, nil, nil
	}

	/**
	Client is authenticated, now we need to connect the client
	to the game. This is done using Channels that connects to
	the Dispatcher (middle layer), which then in turn connects to
	the game engine (upper layer).

	This separates the socket layer from the game layers.
	*/

	// Create and register the needed channels on the dispatcher

	fromClient := make(chan *[]byte)
	toClient := make(chan *[]byte)

	// And register channels on the Dispatcher in the core layer
	core.InitClient(playerLogin.Username, toClient, fromClient)

	return fromClient, toClient, playerLogin
}
