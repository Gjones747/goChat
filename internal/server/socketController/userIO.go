package socketcontroller

import (
	"github.com/Gjones747/goChat/internal/server/models"

	"io"
	"log"
)

// this is the file that will contain all the handling for the users io message stream

func UserIO(user *models.User, roomHub *models.RoomHub) {

	go ioReader(user)

	// this function will launch two go routines for the user a reader and a writer

}

// reads stuff coming in from user and adds to room messages queue
func ioReader(user *models.User) {
	defer user.Connection.Close()
	for {
		header := make([]byte, 2)
		_, err := io.ReadFull(user.Connection, header)
		if err != nil {
			log.Println(err)
			break
		}
	}

}
