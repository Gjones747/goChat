package socketcontroller

import (
	"encoding/binary"
	"fmt"

	"github.com/Gjones747/goChat/internal/server/models"

	"io"
	"log"
)

// reads stuff coming in from user and adds to room messages queue
func UserIOReader(user *models.User) {
	defer user.Connection.Close()

	for {
		header := make([]byte, 2)
		_, err := io.ReadFull(user.Connection, header)
		if err != nil {
			log.Println(err)
			return
		}

		// shifts the byte 7 over so [1xxxxxxx] = [00000001] then 0x01 gits the value of the first bit by masking out the other ones??
		// if fin = 1 it is the last data packet and needs to be added to the room channel
		fin := (header[0] >> 7)
		opcode := header[0] & 0b00001111
		mask := (header[1] >> 7) == 1
		payloadLength := []byte{header[1] & 0b01111111}
		maskKey := make([]byte, 4)
		payloadLengthAsInt := int(payloadLength[0])

		if payloadLength[0] == 126 {
			newPayloadLength := make([]byte, 2)
			_, err := io.ReadFull(user.Connection, newPayloadLength)
			if err != nil {
				log.Println("Error reading 2 byte payload length")
				return
			}
			payloadLength = newPayloadLength
			payloadLengthAsInt = int(binary.BigEndian.Uint16(payloadLength))
		} else if payloadLength[0] == 127 {
			newPayloadLength := make([]byte, 8)
			_, err := io.ReadFull(user.Connection, newPayloadLength)
			if err != nil {
				log.Println("Error reading 8 byte payload length")

				return
			}
			payloadLength = newPayloadLength
			payloadLengthAsInt = int(binary.BigEndian.Uint64(payloadLength))
		} else if payloadLength[0] > 127 {
			log.Println("Bro payload length is not supposed to be more than 127 you got some serious issues")
		}

		_, err = io.ReadFull(user.Connection, maskKey)
		if err != nil {
			log.Println(err)
			break
		}

		fmt.Println(payloadLengthAsInt)
		payLoad := make([]byte, payloadLengthAsInt)

		switch opcode {
		case 0x0:
			fmt.Println("this is a continuation frame")
		case 0x1:
			fmt.Println("this is a text frame")
			if !mask {
				log.Println("Bro you gotta make sure the message is masked")
				return
			}

			_, err = io.ReadFull(user.Connection, payLoad)
			if err != nil {
				log.Println("you got an error reading the payload bytes")
				return
			}

			decodedMessage := decodeMessage(payLoad, maskKey)
			user.Send(decodedMessage)
			fmt.Printf("the message sent was: %s \n", decodedMessage)
		case 0x2:
			fmt.Println("this is a binary frame")
		case 0x3:
			fmt.Println("this is a 'further non control' frame")
		case 0x4:
			fmt.Println("this is a 'further non control' frame")
		case 0x5:
			fmt.Println("this is a 'further non control' frame")
		case 0x6:
			fmt.Println("this is a 'further non control' frame")
		case 0x7:
			fmt.Println("this is a 'further non control' frame")
		case 0x8:
			fmt.Printf("Connection closed with user %s \n", user.Name)
		case 0x9:
			fmt.Println("this is a ping frame")
		case 0xA:
			fmt.Println("this is a pong frame")
		default:
			fmt.Println("what the sigma type frame is ts")
		}

		if fin == 1 {
			log.Println("this is the final part of the message ")
			// this is where the handling for posting the message to the room needs to go
		}
	}

}

func decodeMessage(message []byte, maskKey []byte) []byte {

	for index, _ := range message {
		message[index] = message[index] ^ maskKey[index%4]
	}
	return message
}
