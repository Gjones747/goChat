package socket

import (
	"fmt"
	"net/http"
)

func Upgrader(response http.ResponseWriter, request *http.Request) {
	httpLogger(request)

	if request.Header.Get("Connection") != "Upgrade" || request.Header.Get("Upgrade") != "websocket" {
		http.Error(response, "Did not include websocket request headers", 400)
	}

	key := request.Header.Get("Sec-WebSocket-Key")
	acceptKey := handShakeKey(key)

	fmt.Println(acceptKey)
}
