package models

import (
	"log"
)

type Room struct {
	users map [*User]bool
	messages chan []byte
	register chan *User
	deregester chan *User
}


func InitRoom() *Room {
	return &Room{
		users: make(map [*User]bool),
		messages: make(chan []byte),
	}
}

func (room *Room) StartRoon() {
	for {
		select {
		case user := <-room.register:
			room.users[user] = true
			log.Println("New user added")
		}
	}
}
