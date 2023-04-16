package ebiten

import (
	"bitknife.se/wtf/shared"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

/*
	EBDot is the graphical representation of a corresponding GameObject. (View)
*/

type EbitenObject struct {
	// Read-only!
	gob *shared.GameObject

	// Other properties _NOT_ synced to server could be kept here (read-write)
}

func (ebObj *EbitenObject) Init(gob *shared.GameObject) {
	ebObj.gob = gob
}

func (ebObj *EbitenObject) Draw(screen *ebiten.Image) {
	if ebObj.gob.Kind == "Dot" {
		DrawDot(ebObj, screen)
	}
}

func DrawDot(ebObj *EbitenObject, screen *ebiten.Image) {
	c := color.RGBA{
		R: uint8(0xff),
		G: uint8(0x00),
		B: uint8(0xff),
		A: 0xff}

	vector.DrawFilledCircle(screen, float32(ebObj.gob.X), float32(ebObj.gob.Y), 5, c, true)
}
