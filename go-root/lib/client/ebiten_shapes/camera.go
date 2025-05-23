package ebiten_shapes

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/math/f64"
	"math"
)

type Camera struct {
	// TODO: maybe let camera know of world to calc offsets?
	ViewPort   f64.Vec2
	Position   f64.Vec2
	ZoomFactor int // Positive is more zoomed in
	Rotation   int
}

func (c *Camera) String() string {
	return fmt.Sprintf(
		"T: %.1f, R: %d, S: %d",
		c.Position, c.Rotation, c.ZoomFactor,
	)
}

func (c *Camera) viewportCenter() f64.Vec2 {
	return f64.Vec2{
		c.ViewPort[0] * 0.5,
		c.ViewPort[1] * 0.5,
	}
}

func (c *Camera) worldMatrix() ebiten.GeoM {
	m := ebiten.GeoM{}
	m.Translate(-c.Position[0], -c.Position[1])
	// We want to scale and rotate around center of image / screen
	m.Translate(-c.viewportCenter()[0], -c.viewportCenter()[1])
	m.Scale(
		math.Pow(1.01, float64(c.ZoomFactor)),
		math.Pow(1.01, float64(c.ZoomFactor)),
	)
	m.Rotate(float64(c.Rotation) * 2 * math.Pi / 360)
	m.Translate(c.viewportCenter()[0], c.viewportCenter()[1])
	return m
}

func (c *Camera) Render(world, screen *ebiten.Image) {
	screen.DrawImage(
		world,
		&ebiten.DrawImageOptions{
			GeoM: c.worldMatrix(),
		})
}

func (c *Camera) ScreenToWorld(posX, posY int, wofX, wofY float64) (float64, float64) {
	inverseMatrix := c.worldMatrix()
	if inverseMatrix.IsInvertible() {
		inverseMatrix.Invert()
		wX, wY := inverseMatrix.Apply(float64(posX), float64(posY))
		return wX - wofX, wY - wofY
	} else {
		// When scaling it can happened that matrix is not invertable
		return math.NaN(), math.NaN()
	}
}

func (c *Camera) Reset() {
	// TODO: why to world zero?
	c.Position[0] = 0
	c.Position[1] = 0
	c.Rotation = 0
	c.ZoomFactor = 0
}

func (c *Camera) SetCamera(x, y int) {
	panSpeed := 1.0
	if (c.ZoomFactor) < 0 {
		// ZoomFactor is negative when "zoomed out", so inverting will make pan faster (in pixels)
		panSpeed = 1 - float64(c.ZoomFactor)/10
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		c.Position[0] -= panSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		c.Position[0] += panSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		c.Position[1] -= panSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		c.Position[1] += panSpeed
	}

	if ebiten.IsKeyPressed(ebiten.KeyF) {
		if c.ZoomFactor > -2400 {
			c.ZoomFactor -= 1
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		if c.ZoomFactor < 2400 {
			c.ZoomFactor += 1
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyE) {
		c.Rotation += 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		c.Rotation -= 1
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		c.Reset()
	}
}
