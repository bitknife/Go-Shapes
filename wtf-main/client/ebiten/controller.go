package ebiten

import (
	"bitknife.se/wtf/shared"
)

type EbitenController struct {
	gameObjects    map[string]*shared.GameObject
	fromServerChan chan *shared.Packet
	ebitenGame     *Game
}

func (controller *EbitenController) Run() {

	/*
		Connects and manages packets received from server
		and updates the local game model, so in some sense
		this is the "controller" from an MVC perspective
	*/
	for {

		// Blocks on the stream of packets
		packet := <-controller.fromServerChan

		controller.hardUpdateStrategy(packet)
	}
}

/*
NOTE: This replaces the gameObject completely.

IDEA: Update single properties etc. look at the action property etc.
*/
func (controller *EbitenController) hardUpdateStrategy(packet *shared.Packet) {

	var inGob *shared.GameObject

	// Create if client does not have it
	if packet.GetGameObject() != nil {
		inGob = packet.GetGameObject()
	} else {
		// TODO: Handle other packets as well!
		return
	}

	controller.gameObjects[inGob.Id] = inGob

	if _, ok := controller.ebitenGame.remoteEBObjects[inGob.Id]; !ok {
		controller.ebitenGame.remoteEBObjects[inGob.Id] = &EBGameObject{gob: inGob}
	} else {
		// Update pointer
		ebGob := controller.ebitenGame.remoteEBObjects[inGob.Id]
		ebGob.gob = inGob
	}
}
