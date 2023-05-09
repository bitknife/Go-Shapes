package shapes

import (
	"bitknife.se/wtf/server/game"
)

const (
	ACTION_SHOOT = 0
)

/**
NOTE: All Actions acting on Doers should just post a Mail to the doer if
      its intention and data and have the Doer handle all updates.
*/

func MoveByMouse(player game.Doer, x int32, y int32) {

	mail := game.CreateMail("SET_XY")
	mail.Data["x"] = x
	mail.Data["y"] = y
	player.PostMail(mail)
}

func Shoot(x int32, y int32) {

}
