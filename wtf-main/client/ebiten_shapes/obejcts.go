package ebiten_shapes

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

func (ebObj *EBGameObject) Draw(layer *ebiten.Image) {
	switch ebObj.gob.Kind {
	case shared.GameObjectKind_DOT:
		DrawDot(ebObj, layer)
	case shared.GameObjectKind_BOX:
		DrawBox(ebObj, layer)
	}
}

func DrawDot(ebObj *EBGameObject, screen *ebiten.Image) {
	c := color.RGBA{
		R: uint8(ebObj.gob.IntAttrs["R"]),
		G: uint8(ebObj.gob.IntAttrs["G"]),
		B: uint8(ebObj.gob.IntAttrs["B"]),
		A: uint8(ebObj.gob.IntAttrs["A"]),
	}

	vector.DrawFilledCircle(
		screen,
		float32(ebObj.gob.X),
		float32(ebObj.gob.Y),
		ebObj.gob.FlAttrs["radius"],
		c,
		false)
}

func DrawBox(ebObj *EBGameObject, screen *ebiten.Image) {
	c := color.RGBA{
		R: uint8(ebObj.gob.IntAttrs["R"]),
		G: uint8(ebObj.gob.IntAttrs["G"]),
		B: uint8(ebObj.gob.IntAttrs["B"]),
		A: uint8(ebObj.gob.IntAttrs["A"]),
	}

	vector.DrawFilledCircle(
		screen,
		float32(ebObj.gob.X),
		float32(ebObj.gob.Y),
		ebObj.gob.FlAttrs["radius"],
		c,
		false)
}
