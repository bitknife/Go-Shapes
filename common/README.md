# Common
The server and clients communicate using some common data interchange format. 

The following libraries were evaluated: 

- Protobuf (v2 and v3), 
- Cap´n proto
- FlatBuffers were evaluated.

A comparison made by the developer of Protobuf v2 and Cap´n proto: https://capnproto.org/news/2014-06-17-capnproto-flatbuffers-sbe.html

Cap´n proto came close, but in the end FlatBuffers were chosen as it seems to enjoy the broadest support, has
a rich feature set and was actually developed for gaming in mind.

This statement on a Forum made a point:

"FlatBuffers, like Protobuf, has the ability to leave out arbitrary fields from a table (better yet, it will automatically leave them out if they happen to have the default value). Experience with ProtoBuf has shown me that as data evolves, cases where lots of field get added/removed, or where a lot of fields are essentially empty are extremely common, making this a vital feature for efficiency, and for allowing people to add fields without instantly bloating things.

Cap'n'Proto doesn't have this functionality, and instead just writes zeroes for these fields, taking up space. They recommend using an additional zero-compression on top of the resulting buffer, but this kills the ability to read data without first unpacking it, which is a big selling point of both FlatBuffers and Cap'n'Proto. So Cap'n'Proto is either not truly "zero-parse", or it is rather inefficient/inflexible, take your pick."

(https://github.com/google/flatbuffers/issues/2#issuecomment-215203333)

## Terminology
A packet is an 2 byte *header* that is concatenated with an N byte *message*. Thats it.

        [PACKET] = [SIZE, RESERVED][MESSAGE]

The message is whatever our serialization outputs. The packet is sent over the socket.

## Flatbuffers compiler

## Messages over the wire
Whatever data interchange format is used, there is a need to distinguish what packet is being sent over whatever
socket we are employing (TCP, KDP, UDP etc), and also how many bytes are being sent, since the packets are not of
fixed length.

When designing a more complex application like a game server, we should also put some consideration into if we need
messages on different layers and/or for different purposes. I.e. a high-level packet for the game-domain and some other
kind for logging in, ping etc.

### Packet size
This is the simplest one, a common pattern and often seen recommendation is to just send one initial byte telling how
many bytes to read from the socket.


### Packet type
This is an answer by the Cap'n proto developer:
https://stackoverflow.com/questions/47402232/how-to-distinguish-between-multiple-messages-types-in-capn-proto

