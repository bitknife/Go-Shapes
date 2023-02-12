package core

import (
	"log"
)

const (
	STATIC_PASSWORD = "welcome"
)

func AuthenticateClient(packageData []byte) *PlayerLogin {
	/**
	  Called from the socketserver layer during initial connection.
	*/

	playerLogin := bytesToPlayerLogin(packageData)

	// TODO: Check vs player-db etc.
	if playerLogin.Password == STATIC_PASSWORD {
		log.Println("ACCESS GRANTED FOR Username:", playerLogin.Username)
		return playerLogin
	} else {
		log.Println("ACCESS DENIED FOR Username:", playerLogin.Username, ": Invalid password:", playerLogin.Password)
	}
	//	log.Println("ACCESS DENIED: Invalid message type", messageType, "when authenticating.")
	return nil
}

func bytesToPlayerLogin(messageData []byte) *PlayerLogin {
	/**
	Works for now.
	*/
	packet := BytesToPacket(messageData)
	return packet.GetPlayerLogin()
}
