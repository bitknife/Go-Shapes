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
	HOST = "localhost"
	PORT = "7777"
)

func waitForExitSignals(toServer chan []byte) {
	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal

	shared.BuildLogoutPacket("")
	toServer <- nil

	log.Print("Exiting.")
}

func setupKillTimer(lifetime_sec int) {

	log.Println("Kill timer set to", lifetime_sec, "seconds.")
	kill_timer := time.NewTimer(time.Duration(lifetime_sec) * time.Second)
	go func() {
		<-kill_timer.C
		log.Println("Exiting due to kill timer fired after", lifetime_sec, "sec.")
		syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	}()
}

func main() {
	headless := flags.Bool("headless", false, "Start a client headless.")
	host := flags.StringP("host", "h", HOST, "Server IP or Hostname")
	port := flags.StringP("port", "p", PORT, "Server Port")
	username := flags.StringP("username", "u", shared.RandName("user"), "Player name")
	password := flags.StringP("password", "w", "welcome", "Password")
	lifetime_sec := flags.IntP("lifetime_sec", "l", 0, "Terminate client after this many seconds")
	flags.Parse()

	// Connects and returns two channels for communication
	fromServer, toServer := SetUpNetworking("tcp", *host, *port, *username, *password)

	// Channel from network layer up to UI
	gamePacketsFromServerChannel := make(chan *shared.Packet)
	go HandlePacketsFromServer(fromServer, toServer, gamePacketsFromServerChannel)

	// Central objects shared between network and game engine, keep it simple for now
	gameObjects := make(map[string]*shared.GameObject)

	// For scripted runs of the client typically
	if *lifetime_sec > 0 {
		setupKillTimer(*lifetime_sec)
	}

	// Starts the UI, this blocks
	if *headless == true {
		log.Println("Starting headless client")
		go func() {
			packetCounter := 0
			for {
				// Juste read packets for now.
				<-gamePacketsFromServerChannel
				packetCounter++
			}
		}()
		// NOTE: Blocks
		waitForExitSignals(toServer)

	} else {
		/* Runs on Main thread
		DOC:
			https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2

			https://ebitengine.org/en/documents/cheatsheet.html
		*/
		ebiten.RunEbitenApplication(gameObjects, toServer, gamePacketsFromServerChannel)
	}
}
