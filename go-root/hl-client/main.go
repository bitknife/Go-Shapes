package main

/*
	https://github.com/spf13/pflag
*/
import (
	"bitknife.se/wtf/client"
	"bitknife.se/wtf/shared"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var Protocol = "tcp"
var WTFHost = "localhost"
var WsPort = "8888"

func waitForExitSignals(toServer chan *shared.Packet) {
	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal

	shared.BuildLogoutPacket("")
	toServer <- nil

	log.Print("Exiting.")
}

func setupExitTimer(lifetimeSec int) {

	log.Println("Kill timer set to", lifetimeSec, "seconds.")
	killTimer := time.NewTimer(time.Duration(lifetimeSec) * time.Second)
	go func() {
		<-killTimer.C
		log.Println("Exiting due to kill timer fired after", lifetimeSec, "sec.")
		os.Exit(0)
	}()
}

func startHeadlessClient(
	lifetimeSec int,
	updatesFromSimulation chan *shared.Packet,
	updatesToSimulation chan *shared.Packet) {

	// For scripted runs of the client typically
	if lifetimeSec > 0 {
		setupExitTimer(lifetimeSec)
	}

	go func() {
		packetCounter := 0
		for {
			// Juste read packets for now.
			packet := <-updatesFromSimulation
			if packet == nil {
				log.Println("Server closed connection, exiting.")
				syscall.Exit(0)
			}
			packetCounter++
		}
	}()
	// NOTE: Blocks
	waitForExitSignals(updatesToSimulation)
}

func main() {
	log.Println("Starting headless client")
	updatesToSimulation, updatesFromSimulation, lifetimeSec := client.BootstrapFromCommandLine(Protocol, WTFHost, WsPort)
	startHeadlessClient(lifetimeSec, updatesFromSimulation, updatesToSimulation)
}
