package routes

import (
	"errors"
	"net/http"

	"github.com/Gjones747/goChat/server/models"
	socketcontroller "github.com/Gjones747/goChat/webSocket"
)

// This is where we define what happens when someone makes a room

func CreateRoom(responeWriter http.ResponseWriter, request *http.Request, roomHub *models.RoomHub, roomCode string) error {
	userName := request.URL.Query().Get("user_name")
	if userName == "" {
		return errors.New("did not send a user_name query param")
	}

	connection, err := socketcontroller.Upgrader(responeWriter, request)
	if err != nil {
		return err
	}

	newUser := models.NewUser(connection, userName)

	newRoom := models.InitRoom()
	roomHub.Rooms[roomCode] = newRoom
	newUser.JoinRoom(roomHub.Rooms[roomCode])

	// starts the room
	go roomHub.Rooms[roomCode].StartRoom()

	roomHub.Rooms[roomCode].AddUser(newUser)

	return nil

}
