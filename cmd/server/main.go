package main

import (
	"log"
	"net/http"

	"github.com/Gjones747/goChat/internal/server/models"
	"github.com/Gjones747/goChat/internal/server/routes"
)

func main() {
	mux := http.NewServeMux()
	roomHub := models.MakeRoomHub()

	createRoomHandler := func(w http.ResponseWriter, r *http.Request) {

		roomCode := r.URL.Query().Get("room_code")
		if roomCode == "" {
			log.Println("bro you gotta send a room code when creating a room")
			return
		}

		routes.CreateRoom(w, r, roomHub, roomCode)
	}

	joinRoomHandler := func(w http.ResponseWriter, r *http.Request) {
		roomCode := r.URL.Query().Get("room_code")
		if roomCode == "" {
			log.Println("bro you gotta send a room code when creating a room")
			return
		}

		routes.JoinRoom(w, r, roomHub, roomCode)
	}

	// Endpoint is essentially ..../makeRoom?room_code={whatever}&user_name={username}
	// spawns a new room and automatically adds the user to the room
	mux.Handle("GET /makeRoom", http.HandlerFunc(createRoomHandler))
	mux.Handle("GET /joinRoom", http.HandlerFunc(joinRoomHandler))

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", mux))
}
