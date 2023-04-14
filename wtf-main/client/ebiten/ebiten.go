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
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"log"
)

const (
	screenWidth  = 640
	screenHeight = 480
	scale        = 64
)

/*
NOTE:

MousePosition is a the game object shared with server. (Model)
EBMousePosition is the graphical Ebiten representation of the Game Object. (View)

Thoughts on multi player approach.

All game objects must be created by the server and given an ID.
Updates to game objects are done by certain related events:

	For example, spatial stuff

	Rotate {
		id, deg
	}

	Translate {
		id, x, y
	}

	Other events?
	Create {
		id
	}

	Destroy {
		id
	}
*/
type EBMousePosition struct {
	x, y int
}

func (mp *EBMousePosition) Init() {
	x, y := ebiten.CursorPosition()
	mp.x = x
	mp.y = y
}

func (mp *EBMousePosition) Update(x, y int) {
	mp.x = x
	mp.y = y
}

func (mp *EBMousePosition) Draw(screen *ebiten.Image) {
	c := color.RGBA{
		R: uint8(0xff),
		G: uint8(0x00),
		B: uint8(0xff),
		A: 0xff}

	vector.DrawFilledCircle(screen, float32(mp.x), float32(mp.y), 5, c, true)
}

type Game struct {
	mp EBMousePosition
}

// NewGame is the constructor
func NewGame() *Game {
	g := &Game{}
	g.mp.Init()
	return g
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	x, y := ebiten.CursorPosition()

	// Send all updates from THIS client to server here (async!)

	// Get all updates from server from the current local model

	// Update all
	g.mp.Update(x, y)
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	g.mp.Draw(screen)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func RunEbitenApplication() {

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("We The Forsaken")
	ebiten.SetFullscreen(false)

	theGame := NewGame()

	if err := ebiten.RunGame(theGame); err != nil {
		log.Fatal(err)
	}
}
