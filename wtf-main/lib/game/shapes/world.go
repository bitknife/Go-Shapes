package shapes

import (
	"bitknife.se/wtf/server/game"
	"bitknife.se/wtf/shared"
	"fmt"
	"log"
	"sync"
)

const (
	PLAYER_GOB_ID_PREFIX = "PLAYER"
)

type ShapesAction struct {
}

type ShapesGame struct {
	// Implements DoerGame
	mutex sync.Mutex

	// For external access
	GameObjects map[string]*shared.GameObject

	// QuadTree quadtree.Tree[*shared.GameObject]

	// Note since Doer is an interface, we should _not_ use *Doer, but just Doer.
	// TODO: Access need to be protected for add/remove etc (but not reading?)
	//		 actually, we maybe need a more managed way of working with different
	//		 data structures for the game objects/Doers
	Doers map[string]game.Doer

	ActionsChannel chan ShapesAction

	Physics game.Physics
}

func CreateGame(min int32, max int32, nDots int) *ShapesGame {
	log.Println("Creating dot world with", nDots, "dots.")

	// Allocate
	shapesGame := ShapesGame{}

	// Acquire before modifying the collections below:
	shapesGame.mutex = sync.Mutex{}
	shapesGame.GameObjects = make(map[string]*shared.GameObject)
	shapesGame.Doers = make(map[string]game.Doer)
	// shapesGame.QuadTree = *(quadtree.New[*shared.GameObject](-1000, 1000, 4))

	shapesGame.ActionsChannel = make(chan ShapesAction)
	shapesGame.buildShapes(min, max, nDots)

	// Listen for actions sent to the game
	go shapesGame.ActionListener()

	return &shapesGame
}

func (shapesGame *ShapesGame) Lock() {
	shapesGame.mutex.Lock()
}

func (shapesGame *ShapesGame) Unlock() {
	shapesGame.mutex.Unlock()
}

func (shapesGame *ShapesGame) ActionListener() {
	for {
		action := <-shapesGame.ActionsChannel
		fmt.Println("Got action: ", action)

		// TODO Send to actions handler logic etc..
	}
}

func (shapesGame *ShapesGame) GetGameObjects() map[string]*shared.GameObject {
	return shapesGame.GameObjects
}

func (shapesGame *ShapesGame) AddDoer(id string, doer game.Doer) {
	shapesGame.Lock()
	shapesGame.Doers[id] = doer
	shapesGame.GameObjects[id] = doer.GetGameObject()
	shapesGame.Unlock()

	// Doer pattern, one go routine for each object
	// go doer.Start()
}

func (shapesGame *ShapesGame) RemoveDoer(id string) {
	shapesGame.Lock()
	shapesGame.Unlock()
	// TODO: remove from both
}

func (shapesGame *ShapesGame) buildShapes(min int32, max int32, nObjs int) {
	// Create a bunch of dots within the bounds
	for i := 1; i <= nObjs; i++ {
		// radius := float32(shared.RandInt(15, 15))
		doerObj := CreateRandomBox(shapesGame, min, max)
		shapesGame.AddDoer(doerObj.Id, doerObj)
	}
}

func (shapesGame *ShapesGame) Update() {

	// TODO use waitGroup?
	doneChan := make(chan string)

	// Lock as to not have objects added or removed during update
	shapesGame.Lock()
	for _, doer := range shapesGame.Doers {
		go doer.UpdateGL(doneChan)
	}
	shapesGame.Unlock()

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
		playerBubble := CreateBoxGameObject(shapesGame, playerGobId, 0, 0, 5, 5, 255, 255, 255)
		playerBubble.Id = playerGobId
		shapesGame.AddDoer(playerBubble.Id, playerBubble)
	}

	player := shapesGame.Doers[playerGobId]
	if packet.GetMouseInput() != nil {
		x := packet.GetMouseInput().MouseX
		y := packet.GetMouseInput().MouseY

		// IDEA: Another solution
		// shapesGame.ActionsChannel

		// A bit special for shapes game, this is not how we would do it
		// when moving a player. That must be done through physics engine etc.
		MoveByMouse(player, x, y)

	} else if packet.GetPlayerLogout() != nil {
		// TODO
	}
}
