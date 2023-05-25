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

var TcpPort = "7777"

/*
BootstrapFromCommandLine The default values are typically hard set from the caller at build time.
*/
func BootstrapFromCommandLine(defaultWTFProtocol, defaultWTFHost, defaultWsPort string) (chan *shared.Packet, chan *shared.Packet, int) {
	protocol := flags.StringP("protocol", "p", defaultWTFProtocol, "Network protocol: [ws|wss|tcp]")
	host := flags.StringP("host", "h", defaultWTFHost, "Server IP or Hostname")

	username := flags.StringP("username", "u", shared.RandName("user"), "Player name")
	password := flags.StringP("password", "w", "welcome", "Password")

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
			}
		}()

		// NOTE: This Runs a local simulation and receiving inputs as well!
		shapesGame := shapes.CreateGame(0, 0, 300)
		go game.Run(30, packetsForFrame, shapesGame)
		go game.UserInputRunner("local", updatesToSimulation)

	} else {

		// Connects and returns two channels for communication to  a remote server
		fromServer, toServer := SetUpClientCommunication(*protocol, *host, TcpPort, defaultWsPort, *username, *password)

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
