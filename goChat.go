package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Gjones747/goChat/internal/socket"
)

func main() {
	fmt.Println("Here")
	http.HandleFunc("/ws", socket.Upgrader)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
