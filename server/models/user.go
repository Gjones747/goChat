package models

import (
	"log"
	"net"

	socketcontroller "github.com/Gjones747/goChat/webSocket"
)

type User struct {
	room       *Room
	Connection net.Conn
	Name       string
	send       chan []byte
}

func NewUser(connection net.Conn, userName string) *User {
	return &User{
		room:       nil,
		Connection: connection,
		Name:       userName,
		send:       make(chan []byte),
	}
}

func (user *User) JoinRoom(room *Room) {
	user.room = room
}

func (user *User) userIOWriter() {
	for msg := range user.send {
		socketcontroller.SendToUser(user.Connection, msg)
	}
}

func (user *User) SendMessage(message []byte) {
	user.room.messages <- message
}

func (user *User) Recieve(message []byte) {
	user.send <- message
}

// reads stuff coming in from user and adds to room messages queue
func (user *User) UserIOReader() {
	// all logic for closing down a user goes in this call
	defer user.room.RemoveUser(user)

	for {
		message, err := socketcontroller.ReadFromUser(user.Connection)

		if err != nil {
			log.Println("user Reader failed to run")
			return
		}

		user.SendMessage(message)
	}

}
