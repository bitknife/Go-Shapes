package ebiten

import (
	"bitknife.se/wtf/shared"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

const (
	screenWidth  = 500
	screenHeight = 500
	scale        = 64
)

/*
RunEbitenApplication takes as argument the shared structure gameObjects (for now). That
structure is updated by the Server. Ebiten reads what it needs from that.

toServer is used by the client to notify the server of user-inputs etc.

https://github.com/sedyh/awesome-ebitengine
*/
func RunEbitenApplication(
	gameObjects map[string]*shared.GameObject,
	toServerChan chan *[]byte,
	fromServerChan chan *shared.Packet,
) {

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("WTF!?")
	ebiten.SetFullscreen(false)

	ebitenGame := CreateGame(toServerChan)

	// Start the controller
	ebitenController := EbitenController{gameObjects, fromServerChan, ebitenGame}
	go ebitenController.Run()

	if err := ebiten.RunGame(ebitenGame); err != nil {
		log.Fatal(err)
	}
}
