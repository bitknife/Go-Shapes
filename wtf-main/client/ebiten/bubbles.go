package ebiten

import (
	"bitknife.se/wtf/shared"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
	"golang.org/x/image/math/f64"
	"image"
)

const (
	LOCAL_DOT_ID = "_localdot"
)

type Game struct {
	world  *ebiten.Image
	camera Camera

	// TODO: Not sure if this type is best, it is generic though
	toSimulationPackets chan *shared.Packet

	// Ebitengine representation of gameObjects
	remoteEBObjects map[string]*EBGameObject

	// Local game objects
	localEBObjects map[string]*EBGameObject
}

func CreateGame(toSimulationPackets chan *shared.Packet, worldWidth int, worldHeight int) *Game {
	// Create a util layer that is volatile

	// Create World with given width, centered at 0,0
	halfW := worldWidth / 2
	halfH := worldHeight / 2
	bounds := image.Rectangle{
		Min: image.Point{X: -halfW, Y: -halfH},
		Max: image.Point{X: halfW, Y: halfH},
	}
	options := &ebiten.NewImageOptions{Unmanaged: false}
	game := Game{
		//world: ebiten.NewImage(worldWidth, worldHeight),
		world: ebiten.NewImageWithOptions(bounds, options),

		camera: Camera{
			ViewPort: f64.Vec2{0, 0},
			Position: f64.Vec2{float64(halfW / 2), float64(halfH / 2)},
		},

		toSimulationPackets: toSimulationPackets,

		remoteEBObjects: make(map[string]*EBGameObject),
		localEBObjects:  make(map[string]*EBGameObject),
	}

	// Create a local cursor that is not sent to server
	localDot := EBGameObject{}
	gobj := shared.GameObject{
		Id:   LOCAL_DOT_ID,
		Kind: shared.GameObjectKind_DOT,
		X:    0,
		Y:    0,
		FlAttrs: map[string]float32{
			"radius": 5,
		},
		IntAttrs: map[string]int32{
			"R": 0x00,
			"G": 0xff,
			"B": 0x00,
			"A": 0x99,
		},
	}
	localDot.Init(&gobj)
	game.localEBObjects[gobj.Id] = &localDot

	return &game
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {

	// TODO: optimize, maybe no need to send in every tick?
	scrX, scrY := ebiten.CursorPosition()

	// Adjust viewport
	g.camera.SetCamera(scrX, scrY)

	// TODO: not the nicest way to adjust coordinates
	wX, wY := g.camera.ScreenToWorld(scrX, scrY, float64(g.world.Bounds().Max.X), float64(g.world.Bounds().Max.Y))

	newX := int32(wX)
	newY := int32(wY)

	posChanged := false
	localDot := g.localEBObjects[LOCAL_DOT_ID]
	if newX != localDot.gob.X {
		localDot.gob.X = newX
		posChanged = true
	}
	if newY != localDot.gob.Y {
		localDot.gob.Y = newY
		posChanged = true
	}

	// Send to server only if changed
	if posChanged {
		// fmt.Println("X", localDot.gob.X, "Y", localDot.gob.Y)
		mouseInput := &shared.MouseInput{
			MouseX:     newX,
			MouseY:     newY,
			RightClick: ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight),
			LeftClick:  ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft),
		}

		// TODO: Refactor, should send an "Action" to game loop (local or remote!)

		// TODO: Move, this is network centric
		g.toSimulationPackets <- shared.BuildMouseInputPacket(mouseInput)
	}
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	// TODO Use correct Image option instead
	g.world.Clear()

	drawGrid(g.world, colornames.Darkgray, 100)

	// Draw on World (or maybe layers?)
	for _, ebitenObject := range g.localEBObjects {
		ebitenObject.Draw(g.world)
	}
	for _, ebitenObject := range g.remoteEBObjects {
		ebitenObject.Draw(g.world)
	}

	// Render the updated World onto our screen
	g.camera.Render(g.world, screen)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
