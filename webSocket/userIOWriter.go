package socketcontroller

import (
	"encoding/binary"
	"encoding/json"
	"log"
	"net"

	"github.com/Gjones747/goChat/api"
)

// this function essentially sends messages from the room back to the user
func SendToUser(connection net.Conn, rawMessage api.Envelope) {
	
	message, err := json.Marshal(rawMessage)
	if err != nil {
		log.Println(err)
	}

	var frame []byte
	var length = len(message)

	// denotes that is a text frame for a message
	frame = append(frame, 0b10000001)

	if length < 126 {
		frame = append(frame, byte(length))
	} else if length <= 0xFFFF {
		frame = append(frame, 126)
		extraLength := make([]byte, 2)
		binary.BigEndian.PutUint16(extraLength, uint16(length))

		frame = append(frame, extraLength...)
	} else {
		frame = append(frame, 127)
		extraLength := make([]byte, 8)
		binary.BigEndian.PutUint64(extraLength, uint64(length))

		frame = append(frame, extraLength...)
	}

	frame = append(frame, message...)
	_, err = connection.Write(frame)
	if err != nil {
		log.Println(err)
	}
}
