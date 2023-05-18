package physics

import (
	"bitknife.se/wtf/shared"
)

func BoxCollider(a *shared.GameObject, b *shared.GameObject) bool {

	if a == b {
		return false
	}
	return boxesOverlap(a.X, a.Y, a.W, a.H, b.X, b.Y, b.W, b.H)
}

func boxesOverlap(x1, y1, w1, h1, x2, y2, w2, h2 int32) bool {
	if x1+w1 < x2 || x2+w2 < x1 || y1+h1 < y2 || y2+h2 < y1 {
		return false
	}
	return true
}
