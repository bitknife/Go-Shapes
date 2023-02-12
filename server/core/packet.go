package core

import (
	"google.golang.org/protobuf/proto"
	"time"
)

func BytesToPacket(buffer []byte) Packet {
	packet := Packet{}
	proto.Unmarshal(buffer, &packet)
	return packet
}

func buildPingWireBytes() []byte {
	ping := Ping{}
	ping.Sent = float32(time.Now().UnixMicro())

	packet := Packet{
		TheMessage: &Packet_Ping{&ping},
	}

	marshal, err := proto.Marshal(&packet)
	if err != nil {
		return nil
	}
	header := make([]byte, 1)
	header[0] = byte(len(marshal))
	wirePacket := append(header, marshal...)
	return wirePacket
}
