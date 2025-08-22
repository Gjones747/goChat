package routes

import (
	"log"
	"net/http"

	"github.com/Gjones747/goChat/internal/server/models"
	socketcontroller "github.com/Gjones747/goChat/internal/server/socketController"
)

// This is where we define what happens when someone makes a room

func CreateRoom(responeWriter http.ResponseWriter, request *http.Request, roomHub *models.RoomHub, roomCode string) {
	newUser, err := socketcontroller.Upgrader(responeWriter, request)
	if err != nil {
		log.Println(err)
	}

	newRoom := models.InitRoom(newUser)
	roomHub.Rooms[roomCode] = newRoom
	newUser.JoinRoom(roomHub.Rooms[roomCode])

	// starts the room
	go roomHub.Rooms[roomCode].StartRoom()

	log.Println("here")

	// starts the new room

}
