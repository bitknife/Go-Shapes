package game

import (
	"bitknife.se/wtf/shared"
)

func CreateDotWorld(gameObjects map[string]*shared.GameObject, nDots int) {
	/*
	   Just some initial GameObjects to work with for initial testing
	*/
	x_min := 100
	x_max := 700
	for i := 1; i < nDots; i++ {
		aDot := shared.GameObject{
			X: int32(shared.RandInt(x_min, x_max)),
			Y: int32(shared.RandInt(x_min, x_max)),
		}
		gameObjects[shared.RandName("dot")] = &aDot
	}
}

func ShakeDots(gameObjects map[string]*shared.GameObject, amp int) {
	for _, gameObject := range gameObjects {
		gameObject.X += int32(shared.RandInt(-amp, amp))
		gameObject.Y += int32(shared.RandInt(-amp, amp))
		// fmt.Println("Shook, ", id)
	}
}
