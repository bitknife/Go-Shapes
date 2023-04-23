package main

import (
	"fmt"
	"net/http"
)

func main() {
	addr := "localhost:8080"
	fmt.Print("Serving at http://", addr)
	fmt.Print("\n")
	http.ListenAndServe(addr, http.FileServer(http.Dir(".")))
}
