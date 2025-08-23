package routes

import (
	"net/http"

	"github.com/Gjones747/goChat/internal/server/models"
	socketcontroller "github.com/Gjones747/goChat/internal/server/socketController"
)

// this is the code for someone when they join a room



func JoinRoom(responseWriter http.ResponseWriter, request *http.Request, roomHub *models.RoomHub, roomCode string) error {
	
	user, err := socketcontroller.Upgrader(responseWriter, request) 
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	roomHub.Rooms[roomCode].AddUser(user)
	user.JoinRoom(roomHub.Rooms[roomCode])

	socketcontroller.UserIO(user, roomHub)

	return nil
}
