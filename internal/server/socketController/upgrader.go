package socketcontroller

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
)

func Upgrader(responseWriter http.ResponseWriter, request *http.Request) {

	fmt.Println(request.Header)

	if request.Header.Get("Sec-WebSocket-Version") != "13" || request.Header.Get("Upgrade") != "websocket" {
		http.Error(responseWriter, "Did not send a websocket upgrade request", 400)
		return
	}

	var key string

	if possibleKey := request.Header.Get("Sec-Websocket-Key"); possibleKey != "" {
		key = possibleKey
	} else {
		http.Error(responseWriter, "Did not include an upgrade key", 404)
		return
	}

	fmt.Println(key)

	hijackedCon, ok := responseWriter.(http.Hijacker)
	if !ok {
		http.Error(responseWriter, "failed to hijack connection", 500)
		return
	}

	_, readWriteBuffer, err := hijackedCon.Hijack()
	if err != nil {
		http.Error(responseWriter, "failed to grab hijacked connection", 500)
		return
	}

	returnKey := fmt.Sprintf("%s258EAFA5-E914-47DA-95CA-C5AB0DC85B11", key)
	fmt.Println(returnKey)
	hasher := sha1.New()
	hasher.Write([]byte(returnKey))

	hashedReturnKey := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

	fmt.Fprintf(readWriteBuffer, "HTTP/1.1 101 Switching Protocols\r\n")
	fmt.Fprintf(readWriteBuffer, "Upgrade: websocket\r\n")
	fmt.Fprintf(readWriteBuffer, "Connection: Upgrade\r\n")
	fmt.Fprintf(readWriteBuffer, "Sec-Websocket-Accept: %s\r\n", string(hashedReturnKey))
	fmt.Fprintf(readWriteBuffer, "\r\n")

	readWriteBuffer.Flush()

	fmt.Println("Websocket upgrade request sent")

}
