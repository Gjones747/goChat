package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Gjones747/goChat/internal/server/models"
	"github.com/Gjones747/goChat/internal/server/socketController"
)

// This is where we define what happens when someone makes a room

func CreateRoom(responeWriter http.ResponseWriter, request *http.Request, roomHub models.RoomHub) {
	newUser, err := socketcontroller.Upgrader(responeWriter, request)
	if err != nil {
		log.Println(err)
	}

	fmt.Printf("new user %s is making a new room!", newUser)
}
