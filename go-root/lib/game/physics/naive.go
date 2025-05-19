package physics

import (
	"bitknife.se/wtf/shared"
	"fmt"
	"math"
)

func BoxCollider(a *shared.GameObject, b *shared.GameObject) bool {

	if a == b {
		return false
	}
	return boxesOverlap(a.X, a.Y, a.W, a.H, b.X, b.Y, b.W, b.H)
}

func boxesOverlap(x1, y1, w1, h1, x2, y2, w2, h2 int32) bool {
	if x1+w1 < x2 || x2+w2 < x1 || y1+h1 < y2 || y2+h2 < y1 {
		return false
	}
	return true
}

func LineAnimator2D(x1, y1, x2, y2, fps, seconds int) {

}

func VelocityCurve(fps, seconds float64) []float64 {
	/*
		Returns slice of velocity for a certain time for:

			v(t) = -((t-M)^2)/M^2 + 1

		This is an upside down quadratic function with A-max = 1
		where M is the middle of the curve on the t-axis.

		First and last item is 0, the mid is 1.

		Example: fps=10 and seconds 2, will return 20 values
	*/
	nFrames := int(fps * seconds)

	M := float64(nFrames) / 2
	Msq := M * M
	tArr := makeRange(0, nFrames)
	vArr := makeRange(0, nFrames)

	for i, t := range tArr {
		v := -(((t - M) * (t - M)) / Msq) + 1
		vArr[i] = v
	}
	return vArr
}

func makeRange(min, max int) []float64 {
	a := make([]float64, max-min+1)
	for i := range a {
		a[i] = float64(min + i)
	}
	return a
}

type CoordinateI64 struct {
	X float64
	Y float64
}

func GetCoordinatesAlongLineInt64(x1, y1, angle, length float64, resolution int) []CoordinateI64 {
	coordinates := make([]CoordinateI64, resolution)

	deltaX := length * math.Cos(angle) / float64(resolution-1)
	deltaY := length * math.Sin(angle) / float64(resolution-1)

	for i := 0; i < resolution; i++ {
		x := x1 + float64(i)*deltaX
		y := y1 + float64(i)*deltaY

		coordinates[i] = CoordinateI64{x, y}
	}

	return coordinates
}

type Coordinate struct {
	X int
	Y int
}

func GetCoordinatesAlongLine(x1, y1 int, angle, length float64, nItems int) []Coordinate {
	coordinates := make([]Coordinate, nItems)

	deltaX := length * math.Cos(angle) / float64(nItems-1)
	deltaY := length * math.Sin(angle) / float64(nItems-1)

	for i := 0; i < nItems; i++ {
		x := x1 + int(math.Round(float64(i)*deltaX))
		y := y1 + int(math.Round(float64(i)*deltaY))

		coordinates[i] = Coordinate{x, y}
	}

	return coordinates
}

type Coordinate64 struct {
	X float64
	Y float64
}

type Particle struct {
	StartCoord   Coordinate64
	Length       float64
	Angle        float64
	Acceleration float64
	MaxSpeed     float64
	FPS          int
}

func AnimateParticle(p Particle) {
	// Calculate the end coordinate
	endCoord := Coordinate64{
		X: p.StartCoord.X + p.Length*math.Cos(p.Angle),
		Y: p.StartCoord.Y + p.Length*math.Sin(p.Angle),
	}

	// Calculate the midpoint and time it takes to reach the midpoint
	midpoint := Coordinate64{
		X: (p.StartCoord.X + endCoord.X) / 2,
		Y: (p.StartCoord.Y + endCoord.Y) / 2,
	}
	midpointDistance := math.Sqrt(math.Pow(midpoint.X-p.StartCoord.X, 2) + math.Pow(midpoint.Y-p.StartCoord.Y, 2))
	midpointTime := midpointDistance / (p.MaxSpeed / 2)

	// Calculate the total animation duration and number of frames
	totalTime := midpointTime * 2
	numFrames := totalTime * float64(p.FPS)

	// Calculate the time increment per frame
	dt := 1.0 / float64(p.FPS)

	// Animation loop
	for i := 0; i <= int(numFrames); i++ {
		t := float64(i) * dt

		var currCoord Coordinate64
		var currSpeed float64

		if t <= midpointTime { // Accelerating phase
			currSpeed = p.Acceleration * t
			currCoord.X = p.StartCoord.X + currSpeed*math.Cos(p.Angle)*t
			currCoord.Y = p.StartCoord.Y + currSpeed*math.Sin(p.Angle)*t
		} else { // Decelerating phase
			remainingTime := totalTime - t
			currSpeed = p.Acceleration * remainingTime
			currCoord.X = endCoord.X - currSpeed*math.Cos(p.Angle)*remainingTime
			currCoord.Y = endCoord.Y - currSpeed*math.Sin(p.Angle)*remainingTime
		}

		// Print the current coordinates
		fmt.Printf("Frame %d: (%f, %f)\n", i+1, currCoord.X, currCoord.Y)

		// Wait for the next frame
		// time.Sleep(time.Second / time.Duration(p.FPS))
	}
}
