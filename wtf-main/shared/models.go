package shared

const (
	GOK_DOT       = "DOT"
	GOK_LOCAL_DOT = "GOK_LOCAL_DOT"
)

type GameObject struct {
	/*
		GameObject represents all entities (objects with an id) in the game.

		- Represent what is shared between client and server
		- Client and server uses kind to figure out what concrete objects to show.
		- Events creates, destroys and manipulates GameObject.
		- GameObject are manipulated by GameObjectEvents defined in the wire protocol (protobuf v3).
	*/
	Id string

	// Updated nano
	Ts int64

	// Hint for client behaviour
	Kind string

	X int32
	Y int32
	Z int32
	W int32
	H int32
	R int32 // rotation in degrees,

	Attributes map[string]string
}
