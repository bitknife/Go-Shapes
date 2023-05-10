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
	FPS = 20
)

// Implements Doer
type ShapesDoer struct {
	Id         string
	GameObject *shared.GameObject
	Mailbox    *goconcurrentqueue.FIFO

	// Reference to Game for interaction with other parts of the game
	// from inside the Doers
	Game *ShapesGame
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
	// TODO: Call a Generic Doer loop method?
	ticTime := game.FPSToDuration(FPS)

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
		dwg.shake(1)
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
	// gameObject.FlAttrs["radius"] = gameObject.FlAttrs["radius"] + float32(shared.RandInt(-1, 1))
}

func (dwg *ShapesDoer) handleMail(mail *game.Mail) {

	if mail.Subject == "SET_XY" {
		collides := false
		for _, other := range dwg.Game.Doers {
			collides = physics.BoxCollider(dwg.GameObject, other.GetGameObject())
			if collides {
				// Message other object
				mailOut := game.CreateMail("COLLIDE")
				other.PostMail(mailOut)
			}
		}
		if !collides {
			dwg.GameObject.X = mail.Data["x"].(int32)
			dwg.GameObject.Y = mail.Data["y"].(int32)
		}
	}

	if mail.Subject == "COLLIDE" {
		dwg.GameObject.IntAttrs["R"] = shared.RandInt(0, 255)
		dwg.GameObject.IntAttrs["G"] = shared.RandInt(0, 255)
		dwg.GameObject.IntAttrs["B"] = shared.RandInt(0, 255)
		// dwg.shake(3)
	}
}
