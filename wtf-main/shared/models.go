package shared

type GameObject struct {
	/*
		GameObject represents all entities (objects with an id) in the game.

		- Represent what is shared between client and server
		- Client and server uses kind to figure out what concrete objects to show.
		- Events creates, destroys and manipulates GameObject.
	*/
	id int64

	// Hint for client behaviour
	kind string

	x   int
	y   int
	w   int
	h   int
	rot int // rotation in degrees,
}
