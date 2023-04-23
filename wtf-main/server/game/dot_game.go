package game

import (
	"bitknife.se/wtf/shared"
	"fmt"
)

const (
	PLAYER_GOB_ID_PREFIX = "PLAYER"
)

type DotWorldGame struct {
	GameObjects map[string]*shared.GameObject
}

func CreateDotWorldGame(min int32, max int32, nDots int) *DotWorldGame {
	dotWorldGame := DotWorldGame{}
	dotWorldGame.GameObjects = make(map[string]*shared.GameObject)
	dotWorldGame.buildDotWorld(min, max, nDots)
	return &dotWorldGame
}

func (dotWorldGame *DotWorldGame) GetGameObjects() map[string]*shared.GameObject {
	return dotWorldGame.GameObjects
}

func (dotWorldGame *DotWorldGame) buildDotWorld(min int32, max int32, nDots int) {
	// Create a bunch of dots within the bounds
	for i := 1; i < nDots; i++ {
		id := shared.RandName("dot")
		x := shared.RandInt(min, max)
		y := shared.RandInt(min, max)

		R := shared.RandInt(64, 200)
		G := shared.RandInt(64, 200)
		B := shared.RandInt(64, 200)

		dotWorldGame.GameObjects[id] = createDot(id, x, y, 4, R, G, B)
	}
}

func (dotWorldGame *DotWorldGame) shakeDots(amp int32) {
	for _, gameObject := range dotWorldGame.GameObjects {
		gameObject.X += shared.RandInt(-amp, amp)
		gameObject.Y += shared.RandInt(-amp, amp)

		// gameObject.FlAttrs["radius"] = gameObject.FlAttrs["radius"] + float32(shared.RandInt(-1, 1))
	}
}

func (dotWorldGame *DotWorldGame) Update() {
	dotWorldGame.shakeDots(2)
}

func (dotWorldGame *DotWorldGame) HandleUserInputPacket(
	username string,
	packet *shared.Packet) {

	if packet.GetMouseInput() != nil {
		x := packet.GetMouseInput().MouseX
		y := packet.GetMouseInput().MouseY
		// fmt.Println("UserInput [MOUSE]: got X =", x, " Y =", y, "from", username)

		playerGobId := PLAYER_GOB_ID_PREFIX + "-" + username

		if _, ok := dotWorldGame.GameObjects[playerGobId]; !ok {
			fmt.Println("===> SPAWNED PLAYER <===")
			dotWorldGame.GameObjects[playerGobId] = createDot(playerGobId, x, y, 5, 128, 128, 255)
		}
		// Update
		playerGob := dotWorldGame.GameObjects[playerGobId]
		playerGob.X = x
		playerGob.Y = y

	} else if packet.GetWasdInput() != nil {
		// TODO
	} else if packet.GetPlayerLogout() != nil {
		// TODO
	}
}

func createDot(
	id string,
	x int32, y int32, radius float32,
	R int32, G int32, B int32) *shared.GameObject {

	return &shared.GameObject{
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
			"A": 0,
		},
		Kind: shared.GameObjectKind_DOT,
	}
}
