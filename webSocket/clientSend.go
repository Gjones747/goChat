package socketcontroller

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"errors"
	"log"
	"net"

	"github.com/Gjones747/goChat/api"
)

// sends client data to server -- only going to have it be able to send one frame for now
func SendFramesToServer(connection net.Conn, message api.Envelope) error {

	// converts messageJson into byte[] JSON
	messageJson, err := json.Marshal(message)
	if err != nil {
		return err
	}

	var frame []byte
	var length = len(messageJson)

	// denotes that is a text frame for a messageJson
	frame = append(frame, 0b10000001)

	bitFlipperMask := byte(0b10000000)

	if length < 126 {
		var lengthByte byte = byte(length) | bitFlipperMask
		log.Println(lengthByte)
		frame = append(frame, lengthByte)
	} else if length <= 0xFFFF {
		var lengthByte byte = byte(126) | bitFlipperMask
		frame = append(frame, lengthByte)
		extraLength := make([]byte, 2)
		binary.BigEndian.PutUint16(extraLength, uint16(length))

		frame = append(frame, extraLength...)
	} else {
		var lengthByte byte = byte(127) | bitFlipperMask
		frame = append(frame, lengthByte)
		extraLength := make([]byte, 8)
		binary.BigEndian.PutUint64(extraLength, uint64(length))

		frame = append(frame, extraLength...)
	}

	mask := make([]byte, 4)
	n, err := rand.Read(mask)

	if n != 4 || err != nil {
		return errors.New("Failed to generate random 4 byte sequence ")
	}

	frame = append(frame, mask...)

	maskedPayload := payloadMasker(mask, messageJson)

	frame = append(frame, maskedPayload...)
	_, err = connection.Write(frame)
	if err != nil {
		log.Println(err)
	}

	return nil

}

// read the name you got this!
func payloadMasker(key []byte, payload []byte) []byte {

	for index, _ := range payload {
		payload[index] = payload[index] ^ key[index%4]
	}

	return payload
}
