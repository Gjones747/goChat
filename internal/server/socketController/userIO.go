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

		log.Printf("message header %08b", header)

		// shifts the byte 7 over so [1xxxxxxx] = [00000001] then 0x01 gits the value of the first bit by masking out the other ones??
		// if fin = 1 it is the last data packet and needs to be added to the room channel 
		fin := (header[0] >> 7) & 0x01
		log.Println(fin)
	}

}
