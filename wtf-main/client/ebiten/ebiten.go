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
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

const (
	screenWidth  = 600
	screenHeight = 400
	scale        = 64
)

/*
RunEbitenApplication takes as argument the shared structure gameObjects (for now). That
structure is updated by the Server. Ebiten reads what it needs from that.

toServer is used by the client to notify the server of user-inputs etc.
*/
func RunEbitenApplication(
	gameObjects map[string]*shared.GameObject,
	toServer chan []byte,
	fromServerChan chan *shared.Packet,
) {

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("WTF!?")
	ebiten.SetFullscreen(false)

	ebitenGame := NewGame(toServer)

	// Start the controller
	ebitenController := EbitenController{gameObjects, fromServerChan, ebitenGame}
	go ebitenController.Run()

	if err := ebiten.RunGame(ebitenGame); err != nil {
		log.Fatal(err)
	}
}
