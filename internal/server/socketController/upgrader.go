package socketcontroller

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Gjones747/goChat/internal/server/models"
)

func Upgrader(responseWriter http.ResponseWriter, request *http.Request) (*models.User, error) {

	if request.Header.Get("Sec-WebSocket-Version") != "13" || request.Header.Get("Upgrade") != "websocket" {
		http.Error(responseWriter, "Did not send a websocket upgrade request", 400)
		return &models.User{}, errors.New("did not send a socket request header")
	}

	var key string

	userName := request.URL.Query().Get("user_name")
	if userName == ""{
		return &models.User{}, errors.New("did not send a user_name query param")
	}

	if possibleKey := request.Header.Get("Sec-Websocket-Key"); possibleKey != "" {
		key = possibleKey
	} else {
		http.Error(responseWriter, "Did not include an upgrade key", 404)
		return &models.User{}, errors.New("did not send a socket request header")
	}

	hijackedCon, ok := responseWriter.(http.Hijacker)
	if !ok {
		http.Error(responseWriter, "failed to hijack connection", 500)
		return &models.User{}, errors.New("did not send a socket request header")
	}

	connection, readWriteBuffer, err := hijackedCon.Hijack()
	if err != nil {
		http.Error(responseWriter, "failed to grab hijacked connection", 500)
		return &models.User{}, errors.New("did not send a socket request header")
	}

	newUser := models.NewUser(connection, userName)

	returnKey := fmt.Sprintf("%s258EAFA5-E914-47DA-95CA-C5AB0DC85B11", key)
	hasher := sha1.New()
	hasher.Write([]byte(returnKey))

	hashedReturnKey := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

	fmt.Fprintf(readWriteBuffer, "HTTP/1.1 101 Switching Protocols\r\n")
	fmt.Fprintf(readWriteBuffer, "Upgrade: websocket\r\n")
	fmt.Fprintf(readWriteBuffer, "Connection: Upgrade\r\n")
	fmt.Fprintf(readWriteBuffer, "Sec-Websocket-Accept: %s\r\n", string(hashedReturnKey))
	fmt.Fprintf(readWriteBuffer, "\r\n")

	readWriteBuffer.Flush()

	log.Println("Websocket upgrade request sent")
	return newUser, nil
}
