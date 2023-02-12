package game

type Immovable struct {
	/**
	Represents a floor, wall, door, pillar etc.
	*/
	entityId string

	kind string

	// Placement (coordinate inside segment)
	fromX int
	toX   int
	fromY int
	toY   int

	// Graphics
	texture   string
	tileSizeX int
	tileSizeY int
}
