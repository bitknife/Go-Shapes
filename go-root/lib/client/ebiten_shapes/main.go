package ebiten_shapes

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

toSimulation is used by the client to notify the server of user-inputs etc.

https://github.com/sedyh/awesome-ebitengine
*/
func RunEbitenApplication(
	toSimulation chan *shared.Packet,
	fromSimulation chan *shared.Packet,
) {

	// Local game objects cache
	gameObjects := make(map[string]*shared.GameObject)

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("WTF!?")
	ebiten.SetFullscreen(false)

	ebitenGame := CreateGame(toSimulation, 2000, 2000)

	// Start the controller
	ebitenController := Controller{gameObjects, fromSimulation, ebitenGame}
	go ebitenController.Run()

	// ebiten.SetCursorMode(ebiten.CursorModeHidden)

	// NOTE: Blocks!
	if err := ebiten.RunGame(ebitenGame); err != nil {
		log.Fatal(err)
	}
}
