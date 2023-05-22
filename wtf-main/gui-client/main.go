package main

import (
	"bitknife.se/wtf/client"
	"bitknife.se/wtf/client/ebiten_shapes"
	"log"
)

var Protocol = "ws"
var WTFHost = "localhost"

func main() {
	log.Println("Starting Ebitengine client")
	updatesToSimulation, updatesFromSimulation, _ := client.BootstrapFromCommandLine(Protocol, WTFHost)
	ebiten_shapes.RunEbitenApplication(updatesToSimulation, updatesFromSimulation)
}
