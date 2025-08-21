package main

import (
	"log"
	"net/http"

	socketcontroller "github.com/Gjones747/goChat/internal/server/socketController"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("GET /ws", http.HandlerFunc(socketcontroller.Upgrader))  

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", mux))
}
