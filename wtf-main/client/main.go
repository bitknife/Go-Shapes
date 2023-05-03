package main

/*
	https://github.com/spf13/pflag
*/
import (
	"bitknife.se/wtf/client/ebiten"
	"bitknife.se/wtf/shared"
	flags "github.com/spf13/pflag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	WTFDevServerHost = "wtf-dev-server.bitknife.se"
	WTFDevServerPort = "7777"
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

func main() {
	standalone := flags.BoolP("standalone", "s", false, "Standalone, ie. single player mode.")
	headless := flags.Bool("headless", false, "Start a client headless.")
	host := flags.StringP("host", "h", WTFDevServerHost, "Server IP or Hostname")
	port := flags.StringP("port", "p", WTFDevServerPort, "Server Port")
	username := flags.StringP("username", "u", shared.RandName("user"), "Player name")
	password := flags.StringP("password", "w", "welcome", "Password")
	lifetime_sec := flags.IntP("lifetime_sec", "l", 0, "Terminate client after this many seconds")
	flags.Parse()

	// Central objects shared between game engine (server or local) and view, keep it simple for now
	gameObjects := make(map[string]*shared.GameObject)

	updatesFromSimulation := make(chan *shared.Packet)
	updatesToSimulation := make(chan *shared.Packet)

	if *standalone == true {

	} else {
		// Connects and returns two channels for communication
		fromServer, toServer := SetUpNetworking("tcp", *host, *port, *username, *password)

		go DeliverPacketsToServer(toServer, updatesToSimulation)

		// Isolates the []byte channels from the
		go ReceivePacketsFromServer(fromServer, updatesFromSimulation)
	}

	// Starts the UI, this blocks
	if *headless == true {

		// For scripted runs of the client typically
		if *lifetime_sec > 0 {
			setupExitTimer(*lifetime_sec)
		}

		log.Println("Starting headless client")
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

	} else {
		/* Runs on Main thread
		DOC:
			https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2

			https://ebitengine.org/en/documents/cheatsheet.html
		*/
		ebiten.RunEbitenApplication(gameObjects, updatesToSimulation, updatesFromSimulation)
	}
}
