package objects

import (
	"bitknife.se/wtf/server/game"
	"bitknife.se/wtf/shared"
	"fmt"
	"github.com/enriquebris/goconcurrentqueue"
	"time"
)

const (
	FPS = 20
)

type ShapesGameObject struct {
	Id         string
	GameObject *shared.GameObject
	Mailbox    *goconcurrentqueue.FIFO
}

func (dwg *ShapesGameObject) PostMail(mail game.Mail) {
	// Note the Mailbox can be of any type
	dwg.Mailbox.Enqueue(mail)
}

func (dwg *ShapesGameObject) readMail() {
	// Handle all incoming mail
	for {
		if dwg.Mailbox.GetLen() == 0 {
			break
		}
		mail, err := dwg.Mailbox.Dequeue()
		if err == nil {
			fmt.Println(dwg.Id, "got Mail:", mail)
		}
	}
}

func (dwg *ShapesGameObject) Start() {
	// TODO: Call a Generic Doer loop method?
	ticTime := game.FPSToDuration(FPS)

	go func() {
		for {
			loopStartTime := time.Now()

			//--- WORK HERE ----
			dwg.Update()

			// dwg.readMail()

			//--- SLEEP HERE ---
			sleepDur := ticTime - time.Since(loopStartTime)
			time.Sleep(sleepDur)
		}
	}()
}

func (dwg *ShapesGameObject) Update() {
	// Artificial load to see how game engine reacts to lengthy calculations
	// during game loop, comparing job scheduling algorithms etc.
	//
	// Increasing this and observing CPU load and game loop load figures
	// is interesting. A couple of things can be observed when we max out CPUs
	// during Simulation, preferably all cores! This can be achieved if using
	// go-routines for simulations:
	//
	//  - Sim time will be > 1.00 and Sleep as negative (negative sleep == no sleep!)
	//    We could start to skip frames here instead of slowing things down for example.
	//  - And actual FPS will fall below its target.
	//
	// million := int64(1000000)
	// shared.BurnCPU(1 * million)

	// The actual "job"
	dwg.shake(1)
}

func (dwg *ShapesGameObject) UpdateGL(doneChan chan string) {
	/**
	UpdateGL is the syncronized Update() version called from the Game Loop
	*/
	dwg.Update()

	// And report done to game loop
	doneChan <- "done"
}

func (dwg *ShapesGameObject) GetGameObject() *shared.GameObject {
	return dwg.GameObject
}

func (dwg *ShapesGameObject) shake(amp int32) {
	dwg.GameObject.X += shared.RandInt(-amp, amp)
	dwg.GameObject.Y += shared.RandInt(-amp, amp)
	// gameObject.FlAttrs["radius"] = gameObject.FlAttrs["radius"] + float32(shared.RandInt(-1, 1))
}

func CreateRandomBubble(min int32, max int32, radius float32) *ShapesGameObject {
	id := shared.RandName("dot")
	x := shared.RandInt(min, max)
	y := shared.RandInt(min, max)

	R := shared.RandInt(64, 200)
	G := shared.RandInt(64, 200)
	B := shared.RandInt(64, 200)

	return CreateBubbleGameObject(id, x, y, radius, R, G, B)
}

func CreateBubbleGameObject(
	id string,
	x int32, y int32, radius float32,
	R int32, G int32, B int32) *ShapesGameObject {

	gObj := &shared.GameObject{
		Id: id,
		X:  x,
		Y:  y,
		FlAttrs: map[string]float32{
			"radius": radius,
		},
		IntAttrs: map[string]int32{
			"R": R,
			"G": G,
			"B": B,
			"A": 255,
		},
		Kind: shared.GameObjectKind_DOT,
	}
	return &ShapesGameObject{
		Id:         id,
		GameObject: gObj,
		// NOTE: The fixed version is faster
		Mailbox: goconcurrentqueue.NewFIFO(),
	}
}

func CreateRandomBox(min int32, max int32) *ShapesGameObject {
	id := shared.RandName("box")
	x := shared.RandInt(min, max)
	y := shared.RandInt(min, max)
	w := shared.RandInt(10, 30)
	h := shared.RandInt(10, 30)

	R := shared.RandInt(64, 200)
	G := shared.RandInt(64, 200)
	B := shared.RandInt(64, 200)

	return CreateBoxGameObject(id, x, y, w, h, R, G, B)
}

func CreateBoxGameObject(
	id string,
	x int32, y int32,
	w int32, h int32,
	R int32, G int32, B int32) *ShapesGameObject {

	gObj := &shared.GameObject{
		Id: id,
		X:  x,
		Y:  y,
		W:  w,
		H:  h,
		IntAttrs: map[string]int32{
			"R": R,
			"G": G,
			"B": B,
			"A": 255,
		},
		Kind: shared.GameObjectKind_BOX,
	}
	return &ShapesGameObject{
		Id:         id,
		GameObject: gObj,
		// NOTE: The fixed version is faster
		Mailbox: goconcurrentqueue.NewFIFO(),
	}
}
