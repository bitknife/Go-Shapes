package core

import (
	"bitknife.se/wtf/shared"
	"log"
)

const (
	STATIC_PASSWORD = "welcome"
)

func AuthenticateClient(playerLogin *shared.PlayerLogin) bool {
	/**
	  Called from the socketserver layer during initial connection.
	*/

	// TODO: Check vs player-db etc.
	if playerLogin == nil {
		log.Println("ACCESS DENIED: First packet was not a playerLogin!")
		return false
	} else {
		if playerLogin.Password == STATIC_PASSWORD {
			// log.Println("ACCESS GRANTED FOR Username:", playerLogin.Username)
			return true
		} else {
			log.Println("ACCESS DENIED FOR Username:", playerLogin.Username, ": Invalid password:", playerLogin.Password)
		}
	}

	//	log.Println("ACCESS DENIED: Invalid message type", messageType, "when authenticating.")
	return false
}
