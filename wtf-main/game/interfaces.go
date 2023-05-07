package game

import "bitknife.se/wtf/shared"

type DoerGame interface {
	Update()

	AddDoer(id string, doer Doer)
	RemoveDoer(id string)
	GetGameObjects() map[string]*shared.GameObject

	HandleUserInputPacket(username string, packet *shared.Packet)
}

/*
Doer: Implements the Actor pattern and is self-sustaining using

	its own go-routine.

	Other Actors post mail to each others mailboxes.
	Actors are responsible for emptying incoming mail.
	Mail.
*/
type Doer interface {
	Start()
	Update()
	UpdateGL(chan string)
	GetGameObject() *shared.GameObject

	PostMail(Mail)
}

/*
Mail: Contains data for informing or manipulating the receiving Actor

	in some way. The contents of the Mail can be very Domain specific.
*/
type Mail interface {
}
