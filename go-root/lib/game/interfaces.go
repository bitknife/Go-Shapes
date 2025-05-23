package game

import "bitknife.se/wtf/shared"

type DoerGame interface {
	Update()

	Lock()
	AddDoer(id string, doer Doer)
	RemoveDoer(id string)
	GetGameObjects() map[string]*shared.GameObject
	Unlock()

	HandleUserInputPacket(username string, packet *shared.Packet)
}

type Physics interface {
	BoxCollider(a shared.GameObject, b shared.GameObject)
}

/*
Doer Implements the Actor pattern and is self-sustaining using its own go-routine.
Other Actors post mail to each others mailboxes.
Actors are responsible for emptying incoming mail.
*/
type Doer interface {
	Start()
	Update()
	UpdateGL(chan string)
	GetGameObject() *shared.GameObject

	PostMail(*Mail)
}

/*
Mail contains data for informing or manipulating the receiving Actor
in some way. The contents of the Mail is domain specific.
*/
type Mail struct {
	Subject string
	Data    map[string]any
}

func CreateMail(subject string) *Mail {
	mail := Mail{
		Subject: subject,
		Data:    make(map[string]any),
	}
	return &mail
}
