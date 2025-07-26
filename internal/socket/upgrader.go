package socket

import "net/http"

func Upgrader(response http.ResponseWriter, request *http.Request) {
	if request.Header.Get("Connection") != "Upgrade" || request.Header.Get("Upgrade") != "websocket" {
		http.Error(response, "Did not include websocket request headers", 400)
	}
}
