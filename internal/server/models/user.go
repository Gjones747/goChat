package models

import (
	"net"
)

type User struct {
	room       *Room
	Connection net.Conn
	Name string
	send       chan []byte
}

func NewUser(connection net.Conn, userName string) *User {
	return &User{
		room:       nil,
		Connection: connection,
		Name: userName,
		send:       make(chan []byte),
	}                     
}

func (user *User) JoinRoom(room *Room) {
	user.room = room
}

