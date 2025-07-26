package socket

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func Upgrader(response http.ResponseWriter, request *http.Request) {

	if request.Header.Get("Connection") != "Upgrade" || request.Header.Get("Upgrade") != "websocket" {
		http.Error(response, "Did not include websocket request headers", 400)

	}

	key := request.Header.Get("Sec-WebSocket-Key")
	acceptKey := handShakeKey(key)

	fmt.Println(acceptKey)

	// essentially next bit takes the underlying http request and hijacks it to perform low level stuff?

	hijack, ok := response.(http.Hijacker)
	if !ok {
		http.Error(response, "500 error hijacking", 500)
	}
	connection, readWriteBuffer, err := hijack.Hijack()
	if err != nil {
		log.Print(err)
	}

	log.Printf("WebSocket Connection established")

	fmt.Fprintf(readWriteBuffer, "HTTP/1.1 101 Switching Protocals\r\n")
	fmt.Fprintf(readWriteBuffer, "Upgrade: websocket\r\n")
	fmt.Fprintf(readWriteBuffer, "Connection: Upgrade\r\n")
	fmt.Fprintf(readWriteBuffer, "Sec-WebSocket-Accept: %s\r\n", acceptKey)
	fmt.Fprint(readWriteBuffer, "\r\n")

	readWriteBuffer.Flush()
	time.Sleep(5 * time.Second)

	defer connection.Close()

}
