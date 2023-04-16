package ebiten

import (
	"bitknife.se/wtf/shared"
	"fmt"
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
			gameObjectEvent := packet.GetGameObjectEvent()
			fmt.Println("Received GameObjectEvent: " + gameObjectEvent.String())
		}
	}
}
