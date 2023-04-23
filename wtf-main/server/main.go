package main

import (
	"bitknife.se/wtf/server/core"
	"bitknife.se/wtf/server/game"
	"bitknife.se/wtf/server/socketserver"
	"bitknife.se/wtf/shared"
	"fmt"
	flags "github.com/spf13/pflag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	HOST = "0.0.0.0"
	PORT = "7777"
	TYPE = "tcp"
)

var version = "0.1a"

func printSplash() {
	// Read entire file content, giving us little control but
	// making it very simple. No need to close the file.
	content, err := os.ReadFile("motif.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Convert []byte to string and print to screen
	text := string(content)
	fmt.Println(text)
	fmt.Println("                                                               version", version)
}

func startServer() {
	pingIntervalMsec := flags.IntP("ping_interval_msec", "p", 10000,
		"Interval in milliseconds to ping clients.")
	flags.Parse()

	printSplash()

	// Fancy console for the future!
	// go StartConsole()

	/**
	  Handles the TCP connections, moving messages through
	  the channels and to/from each socket.
	*/
	log.Println("Starting server on", HOST, ":", PORT)
	go socketserver.Run(HOST, PORT)

	/**
	This is the layer separating sockets from the game.
	*/
	log.Println("Ping interval", *pingIntervalMsec, "msec.")
	go core.Run(*pingIntervalMsec)

	packetBroadCastChannel := make(chan []*shared.Packet)
	go core.PacketBroadCaster(packetBroadCastChannel)

	/**
	Main serverside game loop
	*/
	dotWorldGame := game.CreateDotWorldGame(250, 250, 100)
	go game.Run(packetBroadCastChannel, dotWorldGame)

	go shared.CollectAndPrintMetricsRoutine("WTF server", 5)

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
