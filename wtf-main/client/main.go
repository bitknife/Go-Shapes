package main

import (
	"bitknife.se/wtf/client/ebiten"
	"os"
	"os/signal"
	"syscall"
)

func startClient() {

	// Fancy console for the future!
	// go StartConsole()

	go ebiten.RunEbitenApplication()
}

func waitForExitSignals() {
	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal

	// TODO: Stop client loop, disconnect etc.
}

func main() {

	// Spawns everything we need
	startClient()

	// Waits for SIGINT and SIGTERM to perform shutdown
	waitForExitSignals()

}
