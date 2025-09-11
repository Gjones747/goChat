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
        if r.Method != "GET" {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }
        log.Println("Request Sent to /makeRoom")
        roomCode := r.URL.Query().Get("room_code")
        if roomCode == "" {
            log.Println("bro you gotta send a room code when creating a room")
            return
        }
        routes.CreateRoom(w, r, roomHub, roomCode)
    }

    joinRoomHandler := func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "GET" {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }
        roomCode := r.URL.Query().Get("room_code")
        if roomCode == "" {
            log.Println("bro you gotta send a room code when creating a room")
            return
        }
        if _, ok := roomHub.Rooms[roomCode]; !ok {
            routes.CreateRoom(w, r, roomHub, roomCode)
            log.Println("Room doesnt exist - making a new one anyway")
            return
        }
        routes.JoinRoom(w, r, roomHub, roomCode)
    }

    mux.HandleFunc("/makeRoom", createRoomHandler)
    mux.HandleFunc("/joinRoom", joinRoomHandler)

    mux.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "html/index.html")
    })
	mux.Handle("/chat/", http.StripPrefix("/chat/", http.FileServer(http.Dir("html"))))
    log.Println("Server starting on port 8081")
    log.Fatal(http.ListenAndServe("0.0.0.0:8081", mux))
}
