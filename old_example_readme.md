
# Example 1: Threaded python client and Go server
Cocos 2D client and a Go-based server pinging the client(s). The focus was to get some practice
on how to implement a multiplayer server using Go.

It builds on Example 4) of the python-cocos2d example. This one adds a bunch to that:

- Threaded Python client. Very simple, just to be able to send and receive using a simple Queue.
- Different languages in front and backend makes for a new challenge.
- Proto Buffers generating message stubs in both Python and Go as Wire format is one way of coping with that.
- Go socket server basics.
- Showcasing Go routines for handling multiple clients.
- Showcasing Go routines for healthy separating different parts of the application.
- Player authentication.

One difference between MessagePack (MP) and ProtoBuffers (PB) is that PB can generate "model stubs"
for the supported languages, based on a **common schema**. Sometimes this is NOT what you want (and
in that case I would recommend MessagePack), but in our case, this is exactly what want. 

Why?  Because we expect or game message data structures to evolve frequently and we want to be able to
update both the server and the client quickly without having to manually update the message source
code on each side for each change we make.

We will to define our message protocol using .proto files, and from that we will generate stubs in
the language used by the server (Go) and for the client(s), Python and Possibly javascript.

Read more on Protocol Buffers here:

https://protobuf.dev/

For a new project, I would evaluate the promising library Cap'n Proto: https://capnproto.org/, developed
by one of developers behind Protobuf 2.

The server and client(s) can communicate with each other, but connections are not gracefully shut
down in the example, so expect some broken pipes when shutting down the server early or killing
the client abruptly etc. Left as an exercise to the reader to correct! :)

## Reading the code
This README does not explain the entire architecture. I would advise any reader to follow the code
from the [./server/main.go](./server/main.go) file.

## Running the server and client
Start the go server, for example lik this

    cd example1/server
    go run main.go
    Server up on localhost : 7777
    2023/02/11 11:56:06 Waiting for next client to connect...

And in a similar manner, you should be able to start the client that will connect
to the server using the correct username and password in the PlayerLogin packet.
    
    cd example1/client
    pyenv activate pyg
    python ./main.py

On the server console we will see:

    2023/02/11 11:56:13 Client connected from 127.0.0.1:51586
    2023/02/11 11:56:13 Waiting for next client to connect...
    2023/02/11 11:56:13 ACCESS GRANTED FOR Username: player1
    2023/02/11 11:56:33 Dispatcher got message of type 2 from: player1
    2023/02/11 11:56:33 Dispatcher got message of type 2 from: player1
    2023/02/11 11:56:33 Dispatcher got message of type 2 from: player1

Moving your mouse on the Cocos window will send MouseEvents to the server (poacket type 2)

The server will Ping the client every second (packet type 3), on the client:

    Received packet of type: 3, len: 9
    Received packet of type: 3, len: 9
    Received packet of type: 3, len: 9
    

## Message packet structure
While the protocol buffers specify the messages and allows for effective serial- and deserialization
between the client and server, we still need to consider how to send the messages effectively.

A common strategy is to delimit each packet with the packet size, and since we are sending packets
of different types, we also add a single byte to tell the de-serializing side what type to create:

**[1 byte packet length] + [1 byte packet type] + [N bytes packet data]**

A receiver read strategy will then be:

1. Read 1 byte determining packet size.
2. Read 1 byte determining packet type.
3. Read [packet size] bytes into a buffer.
4. Instantiate message of [packet type].

The above structure is defined in the messages.proto as the type **GameMessage** which is the generic
message-type encapsulating the type and the actual packet. It does not include the packet size

### Go support
Download the sources (I used the "ALL" version) and build them using make && make install

And we also need these:

        go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

The example contains the file messages.proto, compile it to GO using one of the below:

        protoc --go_out=. *.proto

### Python support
Googles Python generator is quite crappy, this one is much better.

https://github.com/danielgtaylor/python-betterproto 

So install that (in the correct virtual env, still "pyg") and have a brief look at that page for 
some basic instruction on usage.

        pip install "betterproto[compiler]"
        pip install betterproto    

To generate the Python message classes we do:

        protoc --python_betterproto_out=. messages.proto

### Updating message model
The neat thing now of course is, that each time we want to update a message, or add a new one we
just generate the message stubs for the client and server with a single command

Essentially combining the above calls to protoc like this, with the addition of generating
the respective stubs in the correct sub-directory:

        protoc --go_out=./server --python_betterproto_out=./client messages.proto

And why not add that command to a simple shell script [gen_stubs.py](./example1/gen_stubs.sh).

        hanseklund@Hans-MacBook-Pro example1 % ./gen_stubs.sh
        Writing __init__.py
        Writing messages.py

### Python client message serialization
Just a brief pointer on how we would use the stubs, in the [mousy_game_client.py](./example1/mousy_game_client.py):

        mouse_move_message = messages.MouseMove(x, y, left_click=False, right_click=False)        
        MessageTCPClient.send(bytes(mouse_move_message))

In those two lines, from a shared model (with the Go-server) we instantiate a Python message
object and just serialize and send it over the wire.


### Go server message de-serialization
And, on the server-side we do this:

    message := PlayerLogin{}
    proto.Unmarshal(buffer, &message)
    return message

### Thoughts on using typed Structs in Go
Above unmarshalling is simple enough, but the Go developer will need to handle the typing hassle, 
which is not  as straight forward as for the Python side. See the go-file 
[packet.go](./server/core/packet.go) for one approach and also how the caller to PacketToGameMessage() 
method needs to be aware of the type returned and act accordingly.

Felt this would need to be made smoother. Maybe doing away with the packet_type byte and/or integrate
it into a single generic protobuf message sent over the wire (ie. socket layer), passing it upwards and
letting the above layers decide how/when to unmarshal the data into richer types.

This relates to what data the Go channels are built to support. I think it is nice to just have the
channels going from the socketserver layer just work with the []byte type. The Login is sort of an
odd-ball, having to Unmarshal it on that very low layer. But I think it is OK, making life fairly
simple to just close the connection on failed login etc.

