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

var Commit string = "dev"
var Host string = "0.0.0.0"
var TcpPort = "7777"
var WsPort = "8888"

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
	metricsInterval int) {

	printSplash()

	// Fancy console for the future!
	// go StartConsole()

	/**
	  Handles the TCP connections, moving messages through
	  the channels and to/from each socket.
	*/
	if enableTCP {
		tcpAddress := Host + ":" + TcpPort
		log.Println("Starting TCP server on", tcpAddress)
		go socketserver.ServeTCP(tcpAddress)
	}

	if enableWebsockets {
		wsAddress := Host + ":" + WsPort
		log.Println("Starting WebSocket server on", wsAddress)
		go socketserver.RunWS(wsAddress)
	}

	// go PingAllClients(*pingIntervalMsec)

	/**
	Broadcast packets routine
	*/
	packetBroadCastChannel := make(chan []*shared.Packet)
	go core.PacketBroadCaster(packetBroadCastChannel)

	/**
	Main serverside game loop
	*/
	shapesGame := shapes.CreateGame(-500, 500, nDots)
	go game.Run(gameLoopFps, packetBroadCastChannel, shapesGame)

	if metricsInterval > 0 {
		go CollectAndPrintMetricsRoutine("WTF server", metricsInterval)
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

	metricsInterval := flags.IntP("metricsInterval", "p", 0, "Print metrics to stdout, 0 = disabled")

	flags.Parse()

	// TODO: global var, not the best.. works for now, singleton config?
	shared.WriteTimeout = *socketWriteTimeoutMs

	// Spawns everything we need
	startServer(*gameLoopFps, *nDots, *pingIntervalMsec, *enableTCP, *enableWS, *metricsInterval)

	// Waits for SIGINT and SIGTERM to perform shutdown
	waitForExitSignals()

}
