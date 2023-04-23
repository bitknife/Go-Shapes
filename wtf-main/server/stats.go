package main

import (
	"bitknife.se/wtf/server/core"
	"bitknife.se/wtf/shared"
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

		log.Println("------------------------------------------------------")
		log.Println("[Go] Go routines .....", gS["numGoroutines"])
		// log.Println("[MEM] Heap:", gS["heapAllocKB"], "| Total:", gS["TotalAllocKB"])
		// log.Println("[NET] kB Sent:", nS.BytesSent/1000, "| Rec:", nS.BytesReceived/1000)
		// log.Println("[NET] Ps Sent:", nS.PacketsSent, "| Rec:", nS.PacketsReceived)
		log.Println("[NET] Max send (ms) ..", nS.MaxSendTime/1000000, "| Packets loss:", nS.PacketsLost)
		log.Println("[BC] Clients .........", bS.NumberOfClients, "| BusyC. Drops:", bS.BusyChannelDrops)

		time.Sleep(time.Duration(intervalSec) * time.Second)
	}
}
