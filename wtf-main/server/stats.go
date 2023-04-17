package main

import (
	"bitknife.se/wtf/shared"
	"log"
	"time"
)

func MetricsManager(interval_sec int) {
	/*
		For keeping an eye on server and game performance.

		Will be expanded to meet need. Just print to screen for now

	*/
	gS := shared.CollectGoStats()

	log.Println("numCpus", gS["numCpus"])

	log.Println("numGoroutines | heapAlloc kB | totalAlloc kB | sent kB | Rec. kB | pS | pR")
	log.Println("--------------------------------------------------------")
	for {
		gS := shared.CollectGoStats()
		nS := shared.GetStats()

		log.Println(gS["numGoroutines"], "|", gS["heapAllocKB"], "|", gS["TotalAllocKB"], "|",
			nS.BytesSent/1000, "|", nS.BytesReceived/1000, "|", nS.PacketsSent, "|", nS.PacketsReceived)
		time.Sleep(time.Duration(interval_sec) * time.Second)
	}
}
