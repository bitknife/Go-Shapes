package shared

import (
	"google.golang.org/protobuf/proto"
	"time"
)

func BytesToPacket(buffer []byte) *Packet {
	/*
		Un-marshals a []byte into a Packet.

		This buffer does NOT contain the length as that is stripped of
		during the packet reception.
	*/
	packet := Packet{}
	err := proto.Unmarshal(buffer, &packet)
	if err != nil {
		return nil
	}
	return &packet
}

func PacketToBytes(packet *Packet) []byte {
	/*
		Marshals a core Packet into []bytes and prepends it with length
		this is the Wire-format sent over the socket
	*/
	marshal, err := proto.Marshal(packet)
	if err != nil {
		return nil
	}
	header := make([]byte, 1)
	header[0] = byte(len(marshal))
	wirePacket := append(header, marshal...)
	return wirePacket
}

func BuildPingPacket() *Packet {
	// A Ping message placed in a Packet, works for now
	ping := Ping{}
	ping.Sent = uint64(time.Now().UnixMicro())

	packet := Packet{
		Payload: &Packet_Ping{&ping},
	}
	return &packet
}

func BuildLoginPacket(username string, password string) *Packet {
	login := PlayerLogin{
		Username: username,
		Password: password,
	}
	packet := Packet{
		Payload: &Packet_PlayerLogin{&login},
	}
	return &packet
}
