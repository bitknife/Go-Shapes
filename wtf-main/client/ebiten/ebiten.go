package ebiten

import (
	"bitknife.se/wtf/shared"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

const (
	screenWidth  = 800
	screenHeight = 600
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
	toSimulation chan *shared.Packet,
	fromSimulation chan *shared.Packet,
) {

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("WTF!?")
	ebiten.SetFullscreen(false)

	ebitenGame := CreateGame(toSimulation, 2000, 2000)

	// Start the controller
	ebitenController := EbitenController{gameObjects, fromSimulation, ebitenGame}

	go ebitenController.Run()

	//ebiten.SetCursorMode(ebiten.CursorModeHidden)

	// NOTE: Blocks!
	if err := ebiten.RunGame(ebitenGame); err != nil {
		log.Fatal(err)
	}
}
