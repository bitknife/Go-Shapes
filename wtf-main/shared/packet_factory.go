package shared

import (
	"google.golang.org/protobuf/proto"
	"time"
)

func BytesToPacket(buffer *[]byte) *Packet {
	/*
		Un-marshals a []byte into a Packet.

		This buffer does NOT contain the length byte as that is stripped of
		during the packet reception.
	*/
	packet := Packet{}
	err := proto.Unmarshal(*buffer, &packet)
	if err != nil {
		return nil
	}
	return &packet
}

func PacketToBytes(packet *Packet) *[]byte {
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
	return &wirePacket
}

func BuildPingPacket() *Packet {
	// A Ping message placed in a Packet, works for now
	data := Ping{
		Sent: uint64(time.Now().UnixMicro()),
	}
	return &Packet{
		Payload: &Packet_Ping{&data},
	}
}

func BuildLoginPacket(username string, password string) *Packet {
	return &Packet{
		Payload: &Packet_PlayerLogin{&PlayerLogin{
			Username: username,
			Password: password,
		}},
	}
}

func BuildLogoutPacket(username string) *Packet {
	return &Packet{
		Payload: &Packet_PlayerLogout{&PlayerLogout{
			Username: username,
		}},
	}
}

func BuildMouseInputPacket(mouseInput *MouseInput) *Packet {
	return &Packet{
		Payload: &Packet_MouseInput{mouseInput},
	}
}

func BuildGameObjectPacket(gameObject *GameObject) *Packet {
	return &Packet{
		Payload: &Packet_GameObject{gameObject},
	}
}

func BuildGameObjectPackets(
	tick int64, gameObjects map[string]*GameObject) []*Packet {

	packets := make([]*Packet, len(gameObjects))

	i := 0
	for _, gobj := range gameObjects {
		gobj.Tick = tick
		packet := Packet{
			Payload: &Packet_GameObject{GameObject: gobj},
		}
		packets[i] = &packet
		i++
	}
	return packets
}
