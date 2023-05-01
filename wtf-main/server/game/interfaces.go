package game

import "bitknife.se/wtf/shared"

type DoerGame interface {
	Update()

	AddDoer(id string, doer Doer)
	RemoveDoer(id string)
	GetGameObjects() map[string]*shared.GameObject

	HandleUserInputPacket(username string, packet *shared.Packet)
}

type Doer interface {
	Start()
	Update(chan string)
	GetGameObject() *shared.GameObject
}
