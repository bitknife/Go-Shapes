package main

import (
	"bitknife.se/wtf/client"
	"bitknife.se/wtf/client/ebiten_shapes"
	"log"
)

var Protocol = "ws"
var WTFHost = "localhost"
var WsPort = "8888"

func main() {
	log.Println("Starting Ebitengine client")
	updatesToSimulation, updatesFromSimulation, _ := client.BootstrapFromCommandLine(Protocol, WTFHost, WsPort)
	ebiten_shapes.RunEbitenApplication(updatesToSimulation, updatesFromSimulation)
}
