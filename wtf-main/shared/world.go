package shared

type Immovable struct {
	/**
	Represents a floor, wall, door, pillar etc.
	*/
	entityId string

	kind string

	// Placement (coordinate inside segment)
	x int
	y int
	w int
	h int

	// Graphics
	texture   string
	tileSizeX int
	tileSizeY int
}
