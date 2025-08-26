package socketcontroller

import (

	"github.com/Gjones747/goChat/internal/server/models"
)

// this is the file that will contain all the handling for the users io message stream

func UserIO(user *models.User, roomHub *models.RoomHub) {

	go UserIOReader(user)

	// this function will launch two go routines for the user a reader and a writer

}












