package main

import (
	"log"
	"net/http"

	"github.com/Gjones747/goChat/internal/server/models"
	"github.com/Gjones747/goChat/internal/server/routes"
)

func main() {
	mux := http.NewServeMux()
	roomHub := models.RoomHub{}

	createRoomHandler := func(w http.ResponseWriter, r *http.Request){
		routes.CreateRoom(w, r, roomHub)
	}

	mux.Handle("GET /ws", http.HandlerFunc(createRoomHandler))

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", mux))
}
