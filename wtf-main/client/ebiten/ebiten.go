// Copyright 2021 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ebiten

import (
	"bitknife.se/wtf/shared"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

const (
	screenWidth  = 640
	screenHeight = 480
	scale        = 64
)

type Game struct {
	// READ-ONLY: This becomes populated by network events.
	gameObjects map[string]*shared.GameObject

	toServer chan []byte

	gobjEvents chan *shared.GameObjectEvent

	// Ebiten representation of gameObjects and also non-game objects
	ebitenObjects map[string]*EbitenObject
}

func NewGame(
	gameObjects map[string]*shared.GameObject,
	toServerChan chan []byte,
) *Game {
	game := Game{
		gameObjects: gameObjects,
		toServer:    toServerChan,
	}
	game.ebitenObjects = make(map[string]*EbitenObject)

	// Set up monitoring for GameObjectEvents
	return &game
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	// TODO: optimize, maybe no need to send in every tick?
	x, y := ebiten.CursorPosition()

	// Not sure if we want to keep the toServer channel this deep
	// into the game.
	// Also, only send on change etc. much to improve here
	pP := shared.BuildGameObjectEvent(int32(x), int32(y))
	g.toServer <- shared.PacketToBytes(pP)

	//
	// NOTE: All transient UI-elements should be updated here as well
	//		 That could be the game UI, notifications etc.
	//
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	for _, ebitenObject := range g.ebitenObjects {
		ebitenObject.Draw(screen)
	}
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

/*
RunEbitenApplication takes as argument the shared structure gameObjects (for now). That
structure is updated by the Server. Ebiten reads what it needs from that.

toServer is used by the client to notify the server of user-inputs etc.
*/
func RunEbitenApplication(
	gameObjects map[string]*shared.GameObject,
	toServer chan []byte,
	fromServerChan chan *shared.Packet) {

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("WTF!?")
	ebiten.SetFullscreen(false)

	// NOTE: gameObjects are READ-ONLY for the VIEW and is updated from the server
	//		 not all sure if we need to supply a reference to that structure
	//		 as the Ebitengine application should observe
	theGame := NewGame(gameObjects, toServer)

	// Receives packets and updates the Ebitengine objects
	go ebitenObjectManager(fromServerChan, theGame)

	if err := ebiten.RunGame(theGame); err != nil {
		log.Fatal(err)
	}
}

func ebitenObjectManager(fromServerChan chan *shared.Packet, theGame *Game) {

	for {
		packet := <-fromServerChan
		if packet.GetGameObjectEvent() != nil {
			gameObjectEvent := packet.GetGameObjectEvent()
			fmt.Println("Received GameObjectEvent: " + gameObjectEvent.String())
		}
	}
}
