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
		dotWorldGame.GameObjects[id] = createDot(id, x, y, 2, 255, 0, 255)
	}
}

func (dotWorldGame *DotWorldGame) shakeDots(amp int32) {
	for _, gameObject := range dotWorldGame.GameObjects {
		gameObject.X += int32(shared.RandInt(-amp, amp))
		gameObject.Y += int32(shared.RandInt(-amp, amp))
		// fmt.Println("Shook, ", id)
	}
}

func (dotWorldGame *DotWorldGame) Update() {
	dotWorldGame.shakeDots(1)
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

	} else if packet.GetPlayerLogout() != nil {

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
		},
		Kind: shared.GameObjectKind_DOT,
	}
}
