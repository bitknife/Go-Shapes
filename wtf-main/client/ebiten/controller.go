package ebiten

import (
	"bitknife.se/wtf/shared"
)

func EbitenController(
	gameObjects map[string]*shared.GameObject,
	fromServerChan chan *shared.Packet,
	ebitenGame *Game,
) {
	/*
		Connects and manages packets received from server
		and updates the local game model, so in some sense
		this is the "controller" from an MVC perspective
	*/
	for {
		packet := <-fromServerChan
		if packet.GetGameObjectEvent() != nil {
			goe := packet.GetGameObjectEvent()

			if _, ok := gameObjects[goe.Id]; !ok {
				// Create GameObject
				gob := &shared.GameObject{}
				gameObjects[goe.Id] = gob

				// And create Ebiten representation to draw
				ebitenGame.remoteEBObjects[goe.Id] = &EBGameObject{gob: gob}
			}
			gob := gameObjects[goe.Id]
			gob.X = goe.X
			gob.Y = goe.Y
			gob.Z = goe.Z

			gob.Ts = goe.Tick
			gob.Kind = goe.Kind
			gob.W = goe.W
			gob.H = goe.H
			gob.R = goe.R

			// TODO: merge?
			gob.Attributes = goe.Attributes
		}
	}
}
