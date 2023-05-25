package shapes

import (
	"bitknife.se/wtf/shared"
	"github.com/enriquebris/goconcurrentqueue"
)

func CreateRandomBubble(
	game *ShapesGame,
	min int32, max int32, radius float32) *ShapesDoer {

	id := shared.RandName("dot")
	x := shared.RandInt(min, max)
	y := shared.RandInt(min, max)

	R := shared.RandInt(64, 200)
	G := shared.RandInt(64, 200)
	B := shared.RandInt(64, 200)

	return CreateBubbleGameObject(game, id, x, y, radius, R, G, B)
}

func CreateBubbleGameObject(
	game *ShapesGame,
	id string,
	x int32, y int32, radius float32,
	R int32, G int32, B int32) *ShapesDoer {

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
	return &ShapesDoer{
		Game:       game,
		Id:         id,
		GameObject: gObj,
		Mailbox:    goconcurrentqueue.NewFIFO(),
	}
}

func CreateRandomBox(
	game *ShapesGame,
	min int32,
	max int32) *ShapesDoer {

	id := shared.RandName("box")
	x := shared.RandInt(min, max)
	y := shared.RandInt(min, max)
	w := shared.RandInt(10, 60)
	h := shared.RandInt(10, 60)

	R := shared.RandInt(100, 255)
	G := shared.RandInt(100, 255)
	B := shared.RandInt(100, 255)

	return CreateBoxDoer(game, id, x, y, w, h, R, G, B)
}

func CreateBoxDoer(
	game *ShapesGame,
	id string,
	x int32, y int32,
	w int32, h int32,
	R int32, G int32, B int32) *ShapesDoer {

	gObj := &shared.GameObject{
		Id: id,
		X:  x,
		Y:  y,
		W:  w,
		H:  h,
		IntAttrs: map[string]int32{
			"R": R,
			"G": G,
			"B": B,
			"A": 255,
		},
		Kind: shared.GameObjectKind_BOX,
	}

	doer := &ShapesDoer{
		Id:         id,
		GameObject: gObj,
		// NOTE: The fixed version is faster
		Mailbox:        goconcurrentqueue.NewFIFO(),
		Game:           game,
		AiIntervalMsec: AiBaseIntervalMsec + shared.RandInt(0, 100),
	}
	doer.Start()
	return doer
}
