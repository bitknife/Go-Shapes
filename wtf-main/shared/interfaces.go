package shared

type DoerGame interface {
	Update()

	AddDoer(id string, doer Doer)
	RemoveDoer(id string)
	GetGameObjects() map[string]*GameObject

	HandleUserInputPacket(username string, packet *Packet)
}

type Doer interface {
	Start()
	Update(chan string)
	GetGameObject() *GameObject
}
