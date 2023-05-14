package ebiten_shapes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

func DrawGrid(layer *ebiten.Image, color color.Color, spacing int, axisColor color.Color) {
	// Draw vertical lines
	bounds := layer.Bounds()
	xMin := bounds.Min.X
	xMax := bounds.Max.X
	yMin := bounds.Min.Y
	yMax := bounds.Max.Y
	for i := xMin; i < xMax; i += spacing {
		vector.StrokeLine(layer, float32(i), float32(yMin), float32(i), float32(layer.Bounds().Dy()), 1, color, true)
	}

	// Draw horizontal lines
	for i := yMin; i < yMax; i += spacing {
		vector.StrokeLine(layer, float32(xMin), float32(i), float32(layer.Bounds().Dx()), float32(i), 1, color, true)
	}

	// Draw Axis
	vector.StrokeLine(layer, float32(0), float32(yMin), float32(0), float32(layer.Bounds().Dy()), 3, axisColor, true)
	vector.StrokeLine(layer, float32(xMin), float32(0), float32(layer.Bounds().Dx()), float32(0), 3, axisColor, true)

}
