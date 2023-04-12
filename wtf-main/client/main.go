package main

import (
	"bitknife.se/wtf/client/ebiten"
)

const (
	HOST = "localhost"
	PORT = "7777"
)

func main() {

	fromServer, toServer := SetUpNetworking("goClient", "welcome")

	go HandlePacketsFromServer(fromServer, toServer)

	// Starts the UI, this blocks
	ebiten.RunEbitenApplication()

}
