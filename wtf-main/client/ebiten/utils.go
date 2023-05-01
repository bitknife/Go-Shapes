package ebiten

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

func drawGrid(layer *ebiten.Image, color color.Color, spacing int) {
	// Draw vertical lines
	for i := 0; i < layer.Bounds().Dx(); i += spacing {
		vector.StrokeLine(layer, float32(i), 0, float32(i), float32(layer.Bounds().Dy()), 1, color, true)
	}

	// Draw horizontal lines
	for i := 0; i < layer.Bounds().Dy(); i += spacing {
		vector.StrokeLine(layer, 0, float32(i), float32(layer.Bounds().Dx()), float32(i), 1, color, true)
	}
}
