package main

import (
	"bitknife.se/wtf/server/core"
	"bitknife.se/wtf/server/game"
	"bitknife.se/wtf/server/socketserver"
	flags "github.com/spf13/pflag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	HOST = "localhost"
	PORT = "7777"
	TYPE = "tcp"
)

func startServer() {
	pingIntervalMsec := flags.IntP("ping_interval_msec", "p", 10000, "Interval in milliseconds to ping clients.")
	flags.Parse()

	// Fancy console for the future!
	// go StartConsole()

	/**
	Main serverside game loop
	*/
	go game.Run()

	/**
	This is the layer separating sockets from the game.
	*/
	go core.Run(*pingIntervalMsec)

	/**
	  Handles the TCP connections, moving messages through
	  the channels and to/from each socket.
	*/
	log.Println("Starting server on", HOST, ":", PORT)

	go socketserver.Run(HOST, PORT)

	go PrintStats(5)

}

func stopServer() {
	log.Println("Stopping server.")

	// TODO: call goroutines to shutdown properly
}

func waitForExitSignals() {
	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal

	stopServer()
}

func main() {

	// Spawns everything we need
	startServer()

	// Waits for SIGINT and SIGTERM to perform shutdown
	waitForExitSignals()

}
