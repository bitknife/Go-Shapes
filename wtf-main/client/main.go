package main

/*
	https://github.com/spf13/pflag
*/
import (
	"bitknife.se/wtf/client/ebiten"
	"bitknife.se/wtf/server/game"
	"bitknife.se/wtf/shared"
	"bubbles"
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
	headless := flags.Bool("headless", false, "Start a client headless.")
	host := flags.StringP("host", "h", WTFDevServerHost, "Server IP or Hostname")
	port := flags.StringP("port", "p", WTFDevServerPort, "Server Port")
	username := flags.StringP("username", "u", shared.RandName("user"), "Player name")
	password := flags.StringP("password", "w", "welcome", "Password")
	lifetimeSec := flags.IntP("lifetime_sec", "t", 0, "Terminate client after this many seconds")
	localSim := flags.BoolP("localsim", "l", true, "Run game locally, no server needed.")
	flags.Parse()

	// Central objects shared between game engine (server or local) and view, keep it simple for now
	gameObjects := make(map[string]*shared.GameObject)

	updatesFromSimulation := make(chan *shared.Packet)
	updatesToSimulation := make(chan *shared.Packet)

	if *localSim == true {
		// Create a local game
		bubbleGame := bubbles.CreateBubbleGame(-100, 100, 500)

		// Game returns all updates needed for each frame
		// This is instead of the serverside broadcaster (list of packets to many clients)
		packetsForFrame := make(chan []*shared.Packet)
		allComplete := make(chan int)

		// This code adapts the game.Run() logic that is built to
		// send updates to all clients for each frame.
		go func() {
			for {
				packets := <-packetsForFrame
				for _, packet := range packets {
					// TODO: Look at each packet to see if we should receive it
					//		 this is for a later optimization, sending different
					//		 batches of packets to different clients dep. on position
					//		 or even what they are suppose to see!
					//		 For now: we receive all as that is the strategy
					updatesFromSimulation <- packet
				}
				// Signal that the local client received all
				allComplete <- 1
			}
		}()

		// NOTE: This Runs a local simulation and receiving inputs as well!
		go game.Run(30, packetsForFrame, allComplete, bubbleGame)
		go game.UserInputRunner("local", updatesToSimulation)

	} else {
		// Connects and returns two channels for communication to  a remote server
		fromServer, toServer := SetUpNetworking("tcp", *host, *port, *username, *password)

		// Connects the packets to/from a remote server based simulation
		go DeliverPacketsToServer(toServer, updatesToSimulation)
		go ReceivePacketsFromServer(fromServer, updatesFromSimulation)
	}

	// Starts the UI, this blocks
	if *headless == true {

		// For scripted runs of the client typically
		if *lifetimeSec > 0 {
			setupExitTimer(*lifetimeSec)
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
