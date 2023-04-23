package main

import (
	"bitknife.se/wtf/server/core"
	"bitknife.se/wtf/shared"
	"fmt"
	"log"
	"time"
)

func CollectAndPrintMetricsRoutine(label string, intervalSec int) {
	/*
		For keeping an eye on server and game performance.

		Will be expanded to meet need. Just print to screen for now

	*/
	gS := shared.CollectGoStats()

	log.Println("numCpus", gS["numCpus"])

	log.Println("numGoroutines | heapAlloc kB | totalAlloc kB | sent kB | Rec. kB | pS | pR | mST ms | nPL")
	for {
		gS = shared.CollectGoStats()
		nS := shared.GetNetStats()
		bS := core.GetBroadcastStats()

		avgSendTime := string("?")

		if nS.PacketsSent > 0 {
			avgSendTime = fmt.Sprintf("%.2f", float64(nS.TotalSendTimeMs)/float64(nS.PacketsSent))
		}

		log.Println("------------------------------------------------------")
		log.Println("Go routines ....................", gS["numGoroutines"])
		log.Println("Clients ........................", bS.NumberOfClients)
		// log.Println("[MEM] Heap:", gS["heapAllocKB"], "| Total:", gS["TotalAllocKB"])
		// log.Println("[NET] kB Sent:", nS.BytesSent/1000, "| Rec:", nS.BytesReceived/1000)
		log.Println("Net: Send Min / Avg / Max (ms)..", nS.MinSendTimeMs, "/", avgSendTime, "/", nS.MaxSendTimeMs)
		log.Println("Net: Packets sent ..............", nS.PacketsSent)
		log.Println("Net: Packets loss ..............", nS.PacketsLost)
		log.Println("Net: Busy drops ................", bS.BusyChannelDrops)

		time.Sleep(time.Duration(intervalSec) * time.Second)
	}
}
