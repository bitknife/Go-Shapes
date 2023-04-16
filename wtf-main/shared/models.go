package shared

type GameObject struct {
	/*
		GameObject represents all entities (objects with an id) in the game.

		- Represent what is shared between client and server
		- Client and server uses kind to figure out what concrete objects to show.
		- Events creates, destroys and manipulates GameObject.
		- GameObject are manipulated by GameObjectEvents defined in the wire protocol (protobuf v3).
	*/
	Id int64

	// Hint for client behaviour
	Kind string

	X   int
	Y   int
	W   int
	H   int
	Rot int // rotation in degrees,
}
