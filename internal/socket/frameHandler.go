package socket

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

func handleFrames(connection net.Conn) {
	defer connection.Close()
	for {
		header := make([]byte, 2)
		_, err := io.ReadFull(connection, header)
		if err != nil {
			log.Print("could not read")
			break
		}

		fin := header[0]&0x80 != 0
		opcode := header[0] & 0x0F
		masked := header[1]&0x80 != 0
		payloadLen := int(header[1] & 0x7F)

		if opcode == 0x8 {
			log.Print("client closed connection")
			break
		}

		if payloadLen == 126 {
			headerCopy := make([]byte, 2)
			io.ReadFull(connection, headerCopy)
			payloadLen = int(headerCopy[0])>>8 | int(headerCopy[1])
		} else if payloadLen == 127 {
			ext := make([]byte, 8)
			if _, err := io.ReadFull(connection, ext); err != nil {
				log.Println("Error reading 64-bit payload length:", err)
				break
			}
			payloadLen64 := binary.BigEndian.Uint64(ext)
			if payloadLen64 > (1 << 31) { // limit to 2GB for sanity
				log.Println("Payload too large, skipping")
				break
			}
			payloadLen = int(payloadLen64)
		}
		var maskKey []byte

		if masked {
			maskKey = make([]byte, 4)
			if _, err := io.ReadFull(connection, maskKey); err != nil {
				log.Println("Error reading maskKey:", err)
				return
			}
		}

		payload := make([]byte, payloadLen)

		if _, err := io.ReadFull(connection, payload); err != nil {
			log.Println("Error reading payload:", err)
			return
		}

		if masked && len(maskKey) == 4 {
			for i := 0; i < len(payload); i++ {
				payload[i] ^= maskKey[i%4]
			}
		}

		msg := string(payload)

		if fin {
			fmt.Println("Message: " + msg)
			SendTextFrame(connection, msg)
		}
	}
}
