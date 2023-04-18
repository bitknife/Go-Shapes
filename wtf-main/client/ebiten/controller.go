package ebiten

import (
	"bitknife.se/wtf/shared"
)

type EbitenController struct {
	gameObjects    map[string]*shared.GameObject
	fromServerChan chan *shared.Packet
	ebitenGame     *Game
}

func (c *EbitenController) Run() {

	/*
		Connects and manages packets received from server
		and updates the local game model, so in some sense
		this is the "controller" from an MVC perspective
	*/
	for {

		// Blocks on the stream of packets
		packet := <-c.fromServerChan

		c.hardUpdateStrategy(packet)
	}
}

/*
NOTE: This replaces the gameObject completely.

IDEA: Update single properties etc. look at the action property etc.
*/
func (c *EbitenController) hardUpdateStrategy(packet *shared.Packet) {
	if packet.GetGameObject() != nil {
		inGob := packet.GetGameObject()

		// Create if client does not have it
		c.gameObjects[inGob.Id] = inGob

		if _, ok := c.ebitenGame.remoteEBObjects[inGob.Id]; !ok {
			c.ebitenGame.remoteEBObjects[inGob.Id] = &EBGameObject{gob: inGob}
		} else {
			// Update pointer
			ebGob := c.ebitenGame.remoteEBObjects[inGob.Id]
			ebGob.gob = inGob
		}
	}
}
