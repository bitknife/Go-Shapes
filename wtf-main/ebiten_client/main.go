package main

import (
	"bitknife.se/wtf/client"
	"bitknife.se/wtf/client/ebiten_shapes"
	"log"
)

func main() {
	log.Println("Starting Ebitengine client")
	updatesToSimulation, updatesFromSimulation, _, _ := client.BootstrapFromCommandLine()
	ebiten_shapes.RunEbitenApplication(updatesToSimulation, updatesFromSimulation)
}
