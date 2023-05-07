package shapes

import (
	"bitknife.se/wtf/server/game"
	"bitknife.se/wtf/shared"
	"log"
	"shapes/objects"
)

const (
	PLAYER_GOB_ID_PREFIX = "PLAYER"
)

type BubbleGame struct {
	// Implements DoerGame

	// For external access
	GameObjects map[string]*shared.GameObject

	// Note since Doer is an interface, we should _not_ use *Doer, but just Doer.
	Doers map[string]game.Doer
}

func CreateGame(min int32, max int32, nDots int) *BubbleGame {
	log.Println("Creating dot world with", nDots, "dots.")
	// Allocate
	bubbleGame := BubbleGame{}
	bubbleGame.GameObjects = make(map[string]*shared.GameObject)
	bubbleGame.Doers = make(map[string]game.Doer)

	bubbleGame.buildDotWorld(min, max, nDots)
	return &bubbleGame
}

func (bubbleGame *BubbleGame) GetGameObjects() map[string]*shared.GameObject {
	return bubbleGame.GameObjects
}

func (bubbleGame *BubbleGame) AddDoer(id string, doer game.Doer) {
	bubbleGame.Doers[id] = doer
	bubbleGame.GameObjects[id] = doer.GetGameObject()

	// Doer pattern, one go routine for each object
	go doer.Start()
}

func (bubbleGame *BubbleGame) RemoveDoer(id string) {
	// TODO: remove from both
}

func (bubbleGame *BubbleGame) buildDotWorld(min int32, max int32, nDots int) {
	// Create a bunch of dots within the bounds
	for i := 1; i <= nDots; i++ {
		radius := float32(shared.RandInt(1, 6))
		pBubble := objects.CreateRandomBubble(min, max, radius)
		bubbleGame.AddDoer(pBubble.Id, pBubble)
	}
	log.Println("Created", len(bubbleGame.GameObjects), "dots.")
}

func (bubbleGame *BubbleGame) Update() {

	// TODO use waitGroup?
	doneChan := make(chan string)

	// IDEA: Just do global work here!
	// 		 Actions upon objects should be posted to them instead
	//		 And each object will handle its own update.
	// for _, doer := range bubbleGame.Doers {
	// IDEA: This could be the "Systems" in ECS maybe?
	// go doer.UpdateGL(doneChan)
	// }

	// And wait for completion
	for todo := len(bubbleGame.Doers); todo > 0; todo-- {
		// Wait for all clients to complete
		<-doneChan
	}
}

func (bubbleGame *BubbleGame) HandleUserInputPacket(
	username string,
	packet *shared.Packet) {

	playerGobId := PLAYER_GOB_ID_PREFIX + "-" + username

	if _, ok := bubbleGame.GameObjects[playerGobId]; !ok {
		log.Println("===> SPAWNED PLAYER <===")
		playerBubble := objects.CreateRandomBubble(0, 0, 15)
		playerBubble.Id = playerGobId
		bubbleGame.AddDoer(playerBubble.Id, playerBubble)
	}

	if packet.GetMouseInput() != nil {
		x := packet.GetMouseInput().MouseX
		y := packet.GetMouseInput().MouseY

		// TODO: Limit movement, should look at distance from
		//		 current pos and determine an acceleration up to a limit
		//		 if this was the real game, such concepts would come
		//		 into play.
		playerGob := bubbleGame.GameObjects[playerGobId]
		playerGob.X = x
		playerGob.Y = y

		// TODO: Continue affect game state elsewhere in code

	} else if packet.GetPlayerLogout() != nil {
		// TODO
	}
}
