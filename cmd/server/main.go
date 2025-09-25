package main

import (
	"log"
	"net/http"

	"github.com/Gjones747/goChat/server/models"
	"github.com/Gjones747/goChat/server/routes"
)

func main() {
	mux := http.NewServeMux()
	roomHub := models.MakeRoomHub()

	createRoomHandler := func(w http.ResponseWriter, r *http.Request) {

		log.Println("Request Sent to /makeRoom")
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

		// checks to make sure the room exists if it doesnt it will call the createroom func

		if _, ok := roomHub.Rooms[roomCode]; !ok {
			routes.CreateRoom(w, r, roomHub, roomCode)
			log.Println("Room doesnt exist - making a new one anyway")
			return
		}
		routes.JoinRoom(w, r, roomHub, roomCode)
	}

	// Endpoint is essentially ..../makeRoom?room_code={whatever}&user_name={username}
	// spawns a new room and automatically adds the user to the room
	mux.HandleFunc("/makeRoom", createRoomHandler)
	mux.HandleFunc("/joinRoom", joinRoomHandler)

	log.Println("Server starting on port 8081")
	log.Fatal(http.ListenAndServe("0.0.0.0:8081", mux))
}
