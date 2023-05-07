package ebiten_shapes

import (
	"bitknife.se/wtf/shared"
)

type Controller struct {
	GameObjects    map[string]*shared.GameObject
	FromServerChan chan *shared.Packet
	EbitenGame     *Game
}

func (controller *Controller) Run() {

	/*
		Connects and manages packets received from server
		and updates the local game model, so in some sense
		this is the "controller" from an MVC perspective
	*/
	for {

		// Blocks on the stream of packets
		packet := <-controller.FromServerChan

		controller.hardUpdateStrategy(packet)
	}
}

/*
NOTE: This replaces the gameObject completely.

IDEA: Update single properties etc. look at the action property etc.
*/
func (controller *Controller) hardUpdateStrategy(packet *shared.Packet) {

	var inGob *shared.GameObject

	// Create if client does not have it
	if packet.GetGameObject() != nil {
		inGob = packet.GetGameObject()
	} else {
		// TODO: Handle other packets as well!
		return
	}

	controller.GameObjects[inGob.Id] = inGob

	if _, ok := controller.EbitenGame.remoteEBObjects[inGob.Id]; !ok {
		controller.EbitenGame.remoteEBObjects[inGob.Id] = &EBGameObject{gob: inGob}
	} else {
		// Update pointer
		ebGob := controller.EbitenGame.remoteEBObjects[inGob.Id]
		ebGob.gob = inGob
	}
}
