package main

import (
	"log"
	"net/http"

	"github.com/Gjones747/goChat/internal/middleWare"
	"github.com/Gjones747/goChat/internal/socket"
)

func main() {
	mux := http.NewServeMux()
	log.Println("Running on :8080")

	mux.Handle("GET /ws", http.HandlerFunc(socket.Upgrader))
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", middleWare.LogMiddleWare(mux)))
}
