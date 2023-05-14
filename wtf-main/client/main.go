package main

/*
	https://github.com/spf13/pflag
*/
import (
	"bitknife.se/wtf/client/ebiten_shapes"
	"bitknife.se/wtf/server/game"
	"bitknife.se/wtf/shared"
	flags "github.com/spf13/pflag"
	"log"
	"os"
	"os/signal"
	"shapes"
	"syscall"
	"time"
)

const (
	WTFLocalHost     = "localhost"
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
	username := flags.StringP("username", "u", shared.RandName("user"), "Player name")
	password := flags.StringP("password", "w", "welcome", "Password")
	protocol := flags.StringP("protocol", "p", "websocket", "Network protocol: websocket or tcp")
	lifetimeSec := flags.IntP("lifetime_sec", "t", 0, "Terminate client after this many seconds")
	localSim := flags.BoolP("localsim", "l", false, "Run game locally, no server needed.")
	flags.Parse()

	updatesFromSimulation := make(chan *shared.Packet)
	updatesToSimulation := make(chan *shared.Packet)

	// TODO: This could be streamlined even further, maybe combine server and client even?
	if *localSim == true {

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
		shapesGame := shapes.CreateGame(-500, 500, 500)
		go game.Run(30, packetsForFrame, allComplete, shapesGame)
		go game.UserInputRunner("local", updatesToSimulation)

	} else {
		// Connects and returns two channels for communication to  a remote server
		fromServer, toServer := SetUpNetworking(*protocol, *host, *username, *password)

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
		ebiten_shapes.RunEbitenApplication(updatesToSimulation, updatesFromSimulation)
	}
}
