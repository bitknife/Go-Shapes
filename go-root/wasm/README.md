# WTF Client

## Build and Run
Simply do:

    ./build.sh
    go run runserver.go

And surf to http://localhost:8080

## WASM Build details

https://www.bradcypert.com/an-introduction-to-targeting-web-assembly-with-golang/

Build

    GOOS=js GOARCH=wasm go build -o main.wasm

Serve using any HTTP server, for example using goexec 

    go get -u github.com/shurcooL/goexec
    goexec 'http.ListenAndServe(:8080, http.FileServer(http.Dir(.)))'
