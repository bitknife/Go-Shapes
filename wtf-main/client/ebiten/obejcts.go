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
type EBGameObject struct {
	// Read-only!
	gob *shared.GameObject

	// Other properties _NOT_ synced to server could be kept here (read-write)
}

func (ebObj *EBGameObject) Init(gob *shared.GameObject) {
	ebObj.gob = gob
}

func (ebObj *EBGameObject) Draw(screen *ebiten.Image) {
	switch ebObj.gob.Kind {
	case shared.GOK_DOT:
		DrawDot(ebObj, screen)
	case shared.GOK_LOCAL_DOT:
		DrawLocalDot(ebObj, screen)
	}
}

func DrawDot(ebObj *EBGameObject, screen *ebiten.Image) {
	c := color.RGBA{
		R: uint8(0xff),
		G: uint8(0x00),
		B: uint8(0xff),
		A: 0xff}

	vector.DrawFilledCircle(screen, float32(ebObj.gob.X), float32(ebObj.gob.Y), 4, c, true)
}

func DrawLocalDot(ebObj *EBGameObject, screen *ebiten.Image) {
	c := color.RGBA{
		R: uint8(0x00),
		G: uint8(0xff),
		B: uint8(0x00),
		A: 0xff}

	vector.DrawFilledCircle(screen, float32(ebObj.gob.X), float32(ebObj.gob.Y), 4, c, true)
}
