# Game server

The WTF game server is written in Go. Follow instructions over att https://go.dev to get your basic Go environment
installed

MAC:

    brew install golang

Ubuntu:
    
    apt-get install golang

Add this to your whatever is your shell init file

    export GOPATH=$HOME/go
    export PATH=$PATH:$GOPATH/bin

And source it, or re-launch your shell. Head over back here and do:

    go mod download

This will also install the protobuf/capnproto libraries needed for the stub generation of the common directory.

To update all modules, do
    
    go get -u
    go mod tidy


## Running the server

You can build it and then run the built binary:

    go build .

Or run it using go run

    go run .

For development purposes you should be able to do:

    go run main.go

## Build more

Include git version (server only for now)

    ./server/ $ go build -ldflags="-X main.Commit=$(git rev-parse HEAD)"

Build for windows:

    GOOS=windows GOARCH=amd64 go build

Better is to make a github actions workflow for windows work.

## Update packages
Do this to update all packages in the GOPATH

    go get -u all
