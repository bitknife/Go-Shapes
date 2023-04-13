package main

/*
	https://github.com/spf13/pflag
*/
import (
	"bitknife.se/wtf/client/ebiten"
	flags "github.com/spf13/pflag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const (
	HOST = "localhost"
	PORT = "7777"
)

func waitForExitSignals(conn net.Conn) {
	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal

	log.Print("Closing connection.")
	err := conn.Close()
	if err != nil {
		log.Println("Failed to close connection.")
		return
	}
	log.Print("Exiting.")
}

func main() {
	headless := flags.Bool("headless", false, "Start a client headless.")
	host := flags.StringP("host", "h", HOST, "Server IP or Hostname")
	port := flags.StringP("port", "p", PORT, "Server Port")
	username := flags.StringP("username", "u", "goClient", "Player name")
	password := flags.StringP("password", "w", "welcome", "Password")
	flags.Parse()

	// Connects and returns two channels for communication
	fromServer, toServer, conn := SetUpNetworking(*host, *port, *username, *password)

	go HandlePacketsFromServer(fromServer, toServer)

	// Starts the UI, this blocks
	if *headless == true {
		log.Println("Starting headless client")
		waitForExitSignals(conn)
	} else {
		ebiten.RunEbitenApplication()
	}

}
