package game

import (
	"bitknife.se/wtf/shared"
)

func CreateDotWorld(gameObjects map[string]*shared.GameObject, min int, max int, nDots int) {
	/*
	   Just some initial GameObjects to work with for initial testing

	*/
	for i := 1; i < nDots; i++ {
		id := shared.RandName("dot")
		dot := shared.GameObject{
			Id: id,
			X:  int32(shared.RandInt(min, max)),
			Y:  int32(shared.RandInt(min, max)),
			FlAttrs: map[string]float32{
				"radius": 8,
			},
			IntAttrs: map[string]int32{
				"R": 0xff,
				"G": 0x00,
				"B": 0xff,
			},
			Kind: shared.GameObjectKind_DOT,
		}
		gameObjects[id] = &dot
	}
}

func ShakeDots(gameObjects map[string]*shared.GameObject, amp int) {
	for _, gameObject := range gameObjects {
		gameObject.X += int32(shared.RandInt(-amp, amp))
		gameObject.Y += int32(shared.RandInt(-amp, amp))
		// fmt.Println("Shook, ", id)
	}
}
