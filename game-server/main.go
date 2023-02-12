package main

import (
	"bitknife.se/core"
	"bitknife.se/game"
	"bitknife.se/socketserver"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func startServer() {

	// Fancy console for the future!
	// go StartConsole()

	/**
	Main serverside game loop
	*/
	go game.Start()

	/**
	  Handles the TCP connections, moving messages through
	  the channels and to/from each socket.
	*/
	go socketserver.Start()

	/**
	This is the layer separating sockets from the game.
	*/
	go core.Start()

}

func stopServer() {
	log.Println("Stopping server.")

	// TODO: call goroutines to shutdown properly
}

func mainLoop() {
	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal

	stopServer()
}

func main() {

	// Spawns everything we need
	startServer()

	// Waits for SIGINT and SIGTERM to perform shutdown
	mainLoop()

}
