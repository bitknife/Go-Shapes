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

type ShapesGame struct {
	// Implements DoerGame

	// For external access
	GameObjects map[string]*shared.GameObject

	// Note since Doer is an interface, we should _not_ use *Doer, but just Doer.
	Doers map[string]game.Doer
}

func CreateGame(min int32, max int32, nDots int) *ShapesGame {
	log.Println("Creating dot world with", nDots, "dots.")
	// Allocate
	shapesGame := ShapesGame{}
	shapesGame.GameObjects = make(map[string]*shared.GameObject)
	shapesGame.Doers = make(map[string]game.Doer)

	shapesGame.buildShapesGame(min, max, nDots)
	return &shapesGame
}

func (shapesGame *ShapesGame) GetGameObjects() map[string]*shared.GameObject {
	return shapesGame.GameObjects
}

func (shapesGame *ShapesGame) AddDoer(id string, doer game.Doer) {
	shapesGame.Doers[id] = doer
	shapesGame.GameObjects[id] = doer.GetGameObject()

	// Doer pattern, one go routine for each object
	go doer.Start()
}

func (shapesGame *ShapesGame) RemoveDoer(id string) {
	// TODO: remove from both
}

func (shapesGame *ShapesGame) buildShapesGame(min int32, max int32, nDots int) {
	// Create a bunch of dots within the bounds
	for i := 1; i <= nDots; i++ {
		radius := float32(shared.RandInt(15, 15))
		pBubble := objects.CreateRandomBubble(min, max, radius)
		shapesGame.AddDoer(pBubble.Id, pBubble)
	}
	log.Println("Created", len(shapesGame.GameObjects), "dots.")
}

func (shapesGame *ShapesGame) Update() {

	// TODO use waitGroup?
	doneChan := make(chan string)

	// IDEA: Just do global work here!
	// 		 Actions upon objects should be posted to them instead
	//		 And each object will handle its own update.
	// for _, doer := range shapesGame.Doers {
	// IDEA: This could be the "Systems" in ECS maybe?
	// go doer.UpdateGL(doneChan)
	// }

	// And wait for completion
	for todo := len(shapesGame.Doers); todo > 0; todo-- {
		// Wait for all clients to complete
		<-doneChan
	}
}

func (shapesGame *ShapesGame) HandleUserInputPacket(
	username string,
	packet *shared.Packet) {

	playerGobId := PLAYER_GOB_ID_PREFIX + "-" + username

	if _, ok := shapesGame.GameObjects[playerGobId]; !ok {
		log.Println("===> SPAWNED PLAYER <===")
		playerBubble := objects.CreateShapesGameObject(playerGobId, 0, 0, 3, 255, 255, 255)
		playerBubble.Id = playerGobId
		shapesGame.AddDoer(playerBubble.Id, playerBubble)
	}

	if packet.GetMouseInput() != nil {
		x := packet.GetMouseInput().MouseX
		y := packet.GetMouseInput().MouseY

		// TODO: Limit movement, should look at distance from
		//		 current pos and determine an acceleration up to a limit
		//		 if this was the real game, such concepts would come
		//		 into play.
		playerGob := shapesGame.GameObjects[playerGobId]
		playerGob.X = x
		playerGob.Y = y

		// TODO: Continue affect game state elsewhere in code

	} else if packet.GetPlayerLogout() != nil {
		// TODO
	}
}
