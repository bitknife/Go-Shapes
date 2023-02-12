package core

import (
	"google.golang.org/protobuf/proto"
	"time"
)

func BytesToPacket(buffer []byte) *Packet {
	packet := Packet{}
	err := proto.Unmarshal(buffer, &packet)
	if err != nil {
		return nil
	}
	return &packet
}

func buildPingPacket() *Packet {
	// A Ping message placed in a Packet, works for now
	ping := Ping{}
	ping.Sent = float32(time.Now().UnixMicro())

	packet := Packet{
		TheMessage: &Packet_Ping{&ping},
	}
	return &packet
}
