# Game server

The WTF game server is written in Go. Follow instructions over att https://go.dev to get your basic Go environment
installed

    brew install golang

Add this to your whatever is your shell init file

    export GOPATH=$HOME/go
    export PATH=$PATH:$GOPATH/bin

And source it, or re-launch your shell. Head over back here and do:

    go mod download

This will also install the protobuf/capnproto libraries needed for the stub generation of the common directory.

## Running the server
For development purposes you should be able to do:

    go run main.go

