package ebiten

import (
	"bitknife.se/wtf/shared"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	toServer chan []byte

	// Ebiten representation of gameObjects and also non-game objects
	remoteEBObjects map[string]*EBGameObject
	localEBObjects  map[string]*EBGameObject
}

func NewGame(
	toServerChan chan []byte,
) *Game {
	game := Game{
		toServer: toServerChan,
	}
	game.remoteEBObjects = make(map[string]*EBGameObject)
	game.localEBObjects = make(map[string]*EBGameObject)

	// Create a local cursor that is not sent to server
	localDot := EBGameObject{}
	gobj := shared.GameObject{
		Id:   "_localdot",
		Kind: shared.GOK_LOCAL_DOT,
		X:    0,
		Y:    0,
		W:    0,
		H:    0,
		R:    0,
	}
	localDot.Init(&gobj)
	game.localEBObjects["dot"] = &localDot

	return &game
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	// TODO: optimize, maybe no need to send in every tick?
	x, y := ebiten.CursorPosition()

	// Not sure if we want to keep the toServer channel this deep
	// into the game.
	// Also, only send on change etc. much to improve here
	pP := shared.BuildGameObjectEvent(int32(x), int32(y))
	g.toServer <- shared.PacketToBytes(pP)

	// NOTE: All transient UI-elements should be updated here as well
	//
	//	That could be the game UI, notifications etc.
	g.localEBObjects["dot"].gob.X = int32(x)
	g.localEBObjects["dot"].gob.Y = int32(y)

	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	for _, ebitenObject := range g.localEBObjects {
		ebitenObject.Draw(screen)
	}
	for _, ebitenObject := range g.remoteEBObjects {
		ebitenObject.Draw(screen)
	}
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
