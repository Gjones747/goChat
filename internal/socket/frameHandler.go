package socket

import (
	"fmt"
	"net"
)

func handleFrames(connection net.Conn) {

	fmt.Println(connection)
}
