package core

import (
	"google.golang.org/protobuf/proto"
	"time"
)

func PacketToGameMessage(buffer []byte, mType int) interface{} {

	// fmt.Println("Got Game Message of type: ", mType)

	switch mType {

	case int(MType_PLAYER_LOGIN):
		{
			message := PlayerLogin{}
			proto.Unmarshal(buffer, &message)
			return message
		}

	case int(MType_MOUSE_EVENT):
		{
			message := MouseEvent{}
			proto.Unmarshal(buffer, &message)
			return message
		}
	}
	return nil
}

// TODO, generalize and optimize as this is a crucial function
func buildPingPacket() []byte {

	ping := Ping{}
	ping.SentEpoch = uint64(time.Now().UnixMicro())

	marshal, err := proto.Marshal(&ping)
	if err != nil {
		return nil
	}
	header := make([]byte, 2)
	header[0] = byte(len(marshal))
	header[1] = byte(MType_PING_EVENT)
	packet := append(header, marshal...)

	return packet
}
