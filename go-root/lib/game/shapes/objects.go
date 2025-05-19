package shapes

import (
	"bitknife.se/wtf/server/game"
	"bitknife.se/wtf/server/game/physics"
	"bitknife.se/wtf/shared"
	"github.com/enriquebris/goconcurrentqueue"
	"strings"
	"time"
)

const (
	FPS                = 20
	AiBaseIntervalMsec = 100

	MaxDistance = 500
)

// ShapesDoer Implements Doer
type ShapesDoer struct {
	Id string
	// NOTE, this is game-state, Doer can also keep other states for simulation purposes
	GameObject *shared.GameObject
	Mailbox    *goconcurrentqueue.FIFO

	// Back-reference to Game for interaction with other parts of the game
	// from inside the Doers
	Game *ShapesGame

	AiIntervalMsec int32
}

func (dwg *ShapesDoer) PostMail(mail *game.Mail) {
	dwg.Mailbox.Enqueue(mail)
}

func (dwg *ShapesDoer) readMail() {
	// Read through all of mailbox
	for {
		if dwg.Mailbox.GetLen() == 0 {
			break
		}
		mail, err := dwg.Mailbox.Dequeue()
		if err == nil {
			m := mail.(*game.Mail)
			dwg.handleMail(m)
		}
	}
}

func (dwg *ShapesDoer) Start() {
	// This is a variant, instead of using a single game loop, let
	// each doer maintain its own life through their own loop. This
	// would decouple all objects action from each other by frame and
	// instead only by time.

	// TODO: Call a Generic Doer loop method?
	ticTime := game.FPSToDuration(FPS)

	// Start the AI
	dwg.AI()

	go func() {
		for {
			loopStartTime := time.Now()

			//--- WORK HERE ----
			dwg.Update()

			//--- SLEEP HERE ---
			sleepDur := ticTime - time.Since(loopStartTime)
			time.Sleep(sleepDur)
		}
	}()
}

func (dwg *ShapesDoer) AI() {
	// Whenever called, will deliver a decision based on the current state
	// NOTE: This method never updates the state of the doer.
	// Also, this does not have to occur in every frame

	// IDEA:
	// 	- move at direction until not colliding.
	if dwg.checkCollisions(false) {
		dwg.shake(shared.RandInt(1, 100))
	}

	// --- And call again at a later time ---
	aiTimer := time.NewTimer(100 * time.Millisecond)
	go func() {
		<-aiTimer.C
		dwg.AI()
	}()
}

func (dwg *ShapesDoer) Update() {
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
	if !strings.Contains(dwg.Id, "PLAYER") {
		dwg.shake(0)

		dwg.returnToCenter()
	}
	dwg.readMail()
}

func (dwg *ShapesDoer) UpdateGL(doneChan chan string) {
	/**
	UpdateGL is the syncronized Update() version called from the Game Loop
	*/
	dwg.Update()

	// And report done to game loop
	doneChan <- "done"
}

func (dwg *ShapesDoer) GetGameObject() *shared.GameObject {
	return dwg.GameObject
}

func (dwg *ShapesDoer) shake(amp int32) {
	dwg.GameObject.X += shared.RandInt(-amp, amp)
	dwg.GameObject.Y += shared.RandInt(-amp, amp)
}

func (dwg *ShapesDoer) returnToCenter() {

	if dwg.GameObject.X > MaxDistance ||
		dwg.GameObject.X < -MaxDistance ||
		dwg.GameObject.Y > MaxDistance ||
		dwg.GameObject.Y < -MaxDistance {
		dwg.ToCenter()
	}
}

func (dwg *ShapesDoer) ToCenter() {
	dwg.GameObject.X = 0
	dwg.GameObject.Y = 0
}

func (dwg *ShapesDoer) checkCollisions(notifyOther bool) bool {
	collides := false
	collidedOnce := false
	for _, other := range dwg.Game.Doers {
		collides = physics.BoxCollider(dwg.GameObject, other.GetGameObject())
		if collides && notifyOther {
			// Message other object
			mailOut := game.CreateMail("COLLIDE")
			other.PostMail(mailOut)
		}
		if collides {
			collidedOnce = true
		}
	}
	return collidedOnce
}

func (dwg *ShapesDoer) handleMail(mail *game.Mail) {

	switch mail.Subject {

	case "SET_XY":
		dwg.GameObject.X = mail.Data["x"].(int32)
		dwg.GameObject.Y = mail.Data["y"].(int32)
		dwg.checkCollisions(true)

	case "COLLIDE":
		dwg.ToCenter()
	}
}
