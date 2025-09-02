package routes

import (
	"errors"
	"log"
	"net/http"

	"github.com/Gjones747/goChat/server/models"
	socketcontroller "github.com/Gjones747/goChat/webSocket"
)

// this is the code for someone when they join a room

func JoinRoom(responseWriter http.ResponseWriter, request *http.Request, roomHub *models.RoomHub, roomCode string) error {
	userName := request.URL.Query().Get("user_name")
	if userName == "" {
		return errors.New("did not send a user_name query param")
	}

	connection, err := socketcontroller.Upgrader(responseWriter, request)
	user := models.NewUser(connection, userName)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	log.Println(roomCode)
	roomHub.Rooms[roomCode].AddUser(user)
	user.JoinRoom(roomHub.Rooms[roomCode])

	return nil
}
