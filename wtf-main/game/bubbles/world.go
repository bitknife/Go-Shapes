package bubbles

import (
	"bitknife.se/wtf/shared"
	"log"
)

const (
	PLAYER_GOB_ID_PREFIX = "PLAYER"
)

// BubbleGameObject
// TODO: Interface?

type BubbleGame struct {
	// Implements DoerGame

	// For external access
	GameObjects map[string]*shared.GameObject

	// Note since Doer is an interface, we should _not_ say *Doer, but just Doer!
	Doers map[string]shared.Doer
}

func CreateBubbleGame(min int32, max int32, nDots int) *BubbleGame {
	log.Println("Creating dot world with", nDots, "dots.")
	// Allocate
	bubbleGame := BubbleGame{}
	bubbleGame.GameObjects = make(map[string]*shared.GameObject)
	bubbleGame.Doers = make(map[string]shared.Doer)

	bubbleGame.buildDotWorld(min, max, nDots)
	return &bubbleGame
}

func (bubbleGame *BubbleGame) GetGameObjects() map[string]*shared.GameObject {
	return bubbleGame.GameObjects
}

func (bubbleGame *BubbleGame) AddDoer(id string, doer shared.Doer) {
	bubbleGame.Doers[id] = doer
	bubbleGame.GameObjects[id] = doer.GetGameObject()
}

func (bubbleGame *BubbleGame) RemoveDoer(id string) {
	// TODO: remove from both
}

func (bubbleGame *BubbleGame) buildDotWorld(min int32, max int32, nDots int) {
	// Create a bunch of dots within the bounds
	for i := 1; i <= nDots; i++ {
		pBubble := CreateRandomBubble(min, max)
		bubbleGame.AddDoer(pBubble.Id, pBubble)
	}
	log.Println("Created", len(bubbleGame.GameObjects), "dots.")
}

func (bubbleGame *BubbleGame) Update() {
	// TODO use waitGroup?
	doneChan := make(chan string)

	for _, doer := range bubbleGame.Doers {
		// Update all objects in parallel works for some kind of updates
		go doer.Update(doneChan)
	}

	// And wait for completion
	for todo := len(bubbleGame.Doers); todo > 0; todo-- {
		// Wait for all clients to complete
		<-doneChan
	}
}

func (bubbleGame *BubbleGame) HandleUserInputPacket(
	username string,
	packet *shared.Packet) {

	if packet.GetMouseInput() != nil {
		x := packet.GetMouseInput().MouseX
		y := packet.GetMouseInput().MouseY
		// fmt.Println("UserInput [MOUSE]: got X =", x, " Y =", y, "from", username)

		playerGobId := PLAYER_GOB_ID_PREFIX + "-" + username

		if _, ok := bubbleGame.GameObjects[playerGobId]; !ok {
			log.Println("===> SPAWNED PLAYER <===")
			doer := CreateBubbleGameObject(playerGobId, x, y, 5, 128, 128, 255)
			bubbleGame.AddDoer(doer.Id, doer)
		}
		// Update
		playerGob := bubbleGame.GameObjects[playerGobId]
		playerGob.X = x
		playerGob.Y = y

	} else if packet.GetPlayerLogout() != nil {
		// TODO
	}
}
