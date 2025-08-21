package models

import (
	"net"
)

type User struct {
	Room       *Room
	Connection net.Conn
	send       chan []byte
}

func NewUser(connection net.Conn, room *Room) *User {
	return &User{
		Room:       room,
		Connection: connection,
		send:       make(chan []byte),
	}
}
