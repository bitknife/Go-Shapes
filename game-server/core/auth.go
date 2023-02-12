package core

import (
	"log"
)

const (
	STATIC_PASSWORD = "welcome"
)

func AuthenticateClient(messageType int, message []byte) *PlayerLogin {
	/**
	  Called from the socketserver layer during initial connection.
	*/

	if messageType == int(MType_PLAYER_LOGIN) {
		playerLogin := bytesToPlayerLogin(messageType, message)

		// TODO: Check vs player-db etc.
		if playerLogin.Password == STATIC_PASSWORD {
			log.Println("ACCESS GRANTED FOR Username:", playerLogin.Username)
			return playerLogin
		} else {
			log.Println("ACCESS DENIED FOR Username:", playerLogin.Username, ": Invalid password:", playerLogin.Password)
		}
	} else {
		log.Println("ACCESS DENIED: Invalid message type", messageType, "when authenticating.")
	}
	return nil
}

func bytesToPlayerLogin(messageType int, messageData []byte) *PlayerLogin {
	/**
	Works for now.
	*/
	gameMessage := PacketToGameMessage(messageData, messageType)
	playerLogin := gameMessage.(PlayerLogin)
	return &playerLogin
}
