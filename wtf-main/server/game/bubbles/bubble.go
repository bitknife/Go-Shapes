package bubbles

import (
	"bitknife.se/wtf/shared"
)

type BubbleGameObject struct {
	Id         string
	GameObject *shared.GameObject

	// NOTE: And its methods!
}

func (dwg *BubbleGameObject) Update() {
	dwg.shake(2)
}

func (dwg *BubbleGameObject) GetGameObject() *shared.GameObject {
	return dwg.GameObject
}

func (dwg *BubbleGameObject) shake(amp int32) {
	dwg.GameObject.X += shared.RandInt(-amp, amp)
	dwg.GameObject.Y += shared.RandInt(-amp, amp)
	// gameObject.FlAttrs["radius"] = gameObject.FlAttrs["radius"] + float32(shared.RandInt(-1, 1))
}

func CreateRandomBubble(min int32, max int32) *BubbleGameObject {
	id := shared.RandName("dot")
	x := shared.RandInt(min, max)
	y := shared.RandInt(min, max)

	R := shared.RandInt(64, 200)
	G := shared.RandInt(64, 200)
	B := shared.RandInt(64, 200)

	return CreateBubbleGameObject(id, x, y, 4, R, G, B)
}

func CreateBubbleGameObject(
	id string,
	x int32, y int32, radius float32,
	R int32, G int32, B int32) *BubbleGameObject {

	gObj := &shared.GameObject{
		Id: id,
		X:  x,
		Y:  y,
		FlAttrs: map[string]float32{
			"radius": radius,
		},
		IntAttrs: map[string]int32{
			"R": R,
			"G": G,
			"B": B,
			"A": 255,
		},
		Kind: shared.GameObjectKind_DOT,
	}
	return &BubbleGameObject{
		Id:         id,
		GameObject: gObj,
	}
}
