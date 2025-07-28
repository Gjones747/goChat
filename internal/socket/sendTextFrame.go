package socket

import "net"

func SendTextFrame(connection net.Conn, msg string) {
	payload := []byte(msg)
	header := []byte{0x81}

	if len(payload) < 126 {
		header = append(header, byte(len(payload)))
	} else if len(payload) <= 65535 {
		header = append(header, 126, byte(len(payload)>>8), byte(len(payload)))
	}

	connection.Write(header)
	connection.Write(payload)
}
