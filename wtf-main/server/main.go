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
	"shapes"
	"syscall"
)

const (
	HOST = "0.0.0.0"
)

var Commit string = "dev"

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
	fmt.Println("version:", Commit)
	// --- All output after this should be done as logs
}

func startServer(
	gameLoopFps int64,
	nDots int,
	pingIntervalMsec int,
	enableTCP bool,
	enableWebsockets bool,
	printMetrics bool) {

	printSplash()

	// Fancy console for the future!
	// go StartConsole()

	/**
	  Handles the TCP connections, moving messages through
	  the channels and to/from each socket.
	*/
	if enableTCP {
		tcpAddress := HOST + ":" + shared.TCP_PORT
		log.Println("Starting TCP server on", tcpAddress)
		go socketserver.RunTCP(tcpAddress)
	}

	if enableWebsockets {
		wsAddress := HOST + ":" + shared.WS_PORT
		log.Println("Starting WebSocket server on", wsAddress)
		go socketserver.RunWS(wsAddress)
	}

	// go PingAllClients(*pingIntervalMsec)

	/**
	Broadcast packets routine
	*/
	packetBroadCastChannel := make(chan []*shared.Packet)
	packetsSentChannel := make(chan int)
	go core.PacketBroadCaster(packetBroadCastChannel, packetsSentChannel)

	/**
	Main serverside game loop
	*/
	shapesGame := shapes.CreateGame(-500, 500, nDots)
	go game.Run(gameLoopFps, packetBroadCastChannel, packetsSentChannel, shapesGame)

	if printMetrics {
		go CollectAndPrintMetricsRoutine("WTF server", 2)
	}
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
	gameLoopFps := flags.Int64P("fps", "f", 30, "Game loop FPS")
	nDots := flags.IntP("dots", "d", 100, "Dots to spawn.")

	enableTCP := flags.BoolP("tcpServer", "t", true, "Enable TCP server")
	enableWS := flags.BoolP("websocketServer", "w", true, "Enable WebSocket server")
	socketWriteTimeoutMs := flags.IntP("socketWriteTimeoutMs", "s", 10, "Socket write timeout in ms")

	pingIntervalMsec := flags.IntP("pingIntervalMsec", "i", 10000,
		"Interval in milliseconds to ping clients.")

	printMetrics := flags.BoolP("printMetrics", "p", true, "Print metrics to stdout")

	flags.Parse()

	// TODO: global var, not the best.. works for now, singleton config?
	shared.WriteTimeout = *socketWriteTimeoutMs

	// Spawns everything we need
	startServer(*gameLoopFps, *nDots, *pingIntervalMsec, *enableTCP, *enableWS, *printMetrics)

	// Waits for SIGINT and SIGTERM to perform shutdown
	waitForExitSignals()

}
