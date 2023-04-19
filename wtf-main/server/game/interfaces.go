package game

import "bitknife.se/wtf/shared"

type WTFGame interface {
	Update()
	GetGameObjects() map[string]*shared.GameObject
	HandleUserInputPacket(username string, packet *shared.Packet)
}
