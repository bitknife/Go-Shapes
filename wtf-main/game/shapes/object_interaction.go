package shapes

import (
	"bitknife.se/wtf/server/game"
	"bitknife.se/wtf/server/game/physics"
	"bitknife.se/wtf/shared"
)

func HandleMail(dwg *ShapesDoer, mail *game.Mail) {

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
