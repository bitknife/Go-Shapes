package client

/*
	https://github.com/spf13/pflag
*/
import (
	"bitknife.se/wtf/server/game"
	"bitknife.se/wtf/shared"
	flags "github.com/spf13/pflag"
	"shapes"
)

// Typically modified compiler LD-flags, ie:
// -ldflags="-X main.WsPort=888 -X main.WTFHost=wtf-dev-server.bitknife.se"

var WTFHost = "localhost"
var TcpPort = "7777"
var WsPort = "8888"
var Protocol = "ws"

func BootstrapFromCommandLine() (chan *shared.Packet, chan *shared.Packet, int) {
	host := flags.StringP("host", "h", WTFHost, "Server IP or Hostname")
	username := flags.StringP("username", "u", shared.RandName("user"), "Player name")
	password := flags.StringP("password", "w", "welcome", "Password")
	protocol := flags.StringP("protocol", "p", Protocol, "Network protocol: websocket or tcp")

	lifetimeSec := flags.IntP("lifetime_sec", "t", 0, "Terminate client after this many seconds")

	localSim := flags.BoolP("localsim", "l", false, "Run game locally, no server needed.")
	flags.Parse()

	updatesFromSimulation := make(chan *shared.Packet)
	updatesToSimulation := make(chan *shared.Packet)

	// TODO: This could be streamlined even further, maybe combine server and client even?
	if *localSim {

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
				// NOTE: must use if synced game loop (signal that the local client received all)
				//       we sometimes experiemnt with this not used.
				// allComplete <- 1
			}
		}()

		// NOTE: This Runs a local simulation and receiving inputs as well!
		shapesGame := shapes.CreateGame(-500, 500, 500)
		go game.Run(30, packetsForFrame, allComplete, shapesGame)
		go game.UserInputRunner("local", updatesToSimulation)

	} else {

		// Connects and returns two channels for communication to  a remote server
		fromServer, toServer := SetUpClientCommunication(*protocol, *host, TcpPort, WsPort, *username, *password)

		// Connects the packets to/from a remote server based simulation
		go DeliverPacketsToServer(toServer, updatesToSimulation)
		go ReceivePacketsFromServer(fromServer, updatesFromSimulation)
	}

	/* Runs on Main thread
	DOC:
		https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2
		https://ebitengine.org/en/documents/cheatsheet.html
	*/
	return updatesToSimulation, updatesFromSimulation, *lifetimeSec
}
