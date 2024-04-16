package physics

import (
	"fmt"
	"math"
	"testing"
)

func TestVelocity(t *testing.T) {

	result := VelocityCurve(30, 3)

	scale := 10

	for i, v := range result {
		fmt.Println(i, "=", float64(scale)*v)
	}
}

func TestGetCoordinatesAlongLine(t *testing.T) {
	startX := 0
	startY := 0
	angle := math.Pi / 6 // 45 degrees in radians
	endLength := 10.0
	resolution := 10

	coordinates := GetCoordinatesAlongLine(startX, startY, angle, endLength, resolution)

	// Print the coordinates
	for _, coord := range coordinates {
		fmt.Println(coord.X, coord.Y)
	}
}

func TestAnimateParticle(t *testing.T) {
	startCoord := Coordinate64{X: 0, Y: 0}
	length := 100.0
	angle := math.Pi / 4
	acceleration := 1.0
	maxSpeed := 10.0
	fps := 10

	p := Particle{
		StartCoord:   startCoord,
		Length:       length,
		Angle:        angle,
		Acceleration: acceleration,
		MaxSpeed:     maxSpeed,
		FPS:          fps,
	}

	AnimateParticle(p)
}
