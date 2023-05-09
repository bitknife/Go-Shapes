package game

import (
	"bitknife.se/wtf/shared"
)

type NaivePhysics struct {
}

func (cmPhys *NaivePhysics) BoxCollider(a shared.GameObject, b shared.GameObject) bool {

	if (a.X+a.W > b.X) && (a.Y+a.H > b.Y) {
		return true
	}
	if (b.X+b.W > a.X) && (b.Y+b.H > a.Y) {
		return true
	}

	// TODO More?
	return false
}
