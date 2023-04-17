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

	newX := int32(x)
	newY := int32(y)

	posChanged := false
	if newX != g.localEBObjects["dot"].gob.X {
		g.localEBObjects["dot"].gob.X = newX
		posChanged = true
	}
	if newY != g.localEBObjects["dot"].gob.Y {
		g.localEBObjects["dot"].gob.Y = newY
		posChanged = true
	}

	// Send to server only if changed
	if posChanged {
		pP := shared.BuildMouseInputPacket(&shared.MouseInput{
			MouseX:     int32(x),
			MouseY:     int32(y),
			RightClick: ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight),
			LeftClick:  ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft),
		})
		g.toServer <- shared.PacketToBytes(pP)
	}
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
