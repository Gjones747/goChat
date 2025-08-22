package models

import (
	"net"
)

type User struct {
	room       *Room
	Connection *net.Conn
	send       chan []byte
}

func NewUser(connection *net.Conn) *User {
	return &User{
		room:       nil,
		Connection: connection,
		send:       make(chan []byte),
	}                     
}

func (user *User) JoinRoom(room *Room) {
	user.room = room
}
