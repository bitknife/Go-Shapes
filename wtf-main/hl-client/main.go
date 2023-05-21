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

func waitForExitSignals(toServer chan *shared.Packet) {
	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal

	shared.BuildLogoutPacket("")
	toServer <- nil

	log.Print("Exiting.")
}

func setupExitTimer(lifetime_sec int) {

	log.Println("Kill timer set to", lifetime_sec, "seconds.")
	kill_timer := time.NewTimer(time.Duration(lifetime_sec) * time.Second)
	go func() {
		<-kill_timer.C
		log.Println("Exiting due to kill timer fired after", lifetime_sec, "sec.")
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
	updatesToSimulation, updatesFromSimulation, lifetimeSec := client.BootstrapFromCommandLine()
	startHeadlessClient(lifetimeSec, updatesFromSimulation, updatesToSimulation)
}
