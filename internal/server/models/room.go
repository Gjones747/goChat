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


// each new room has to start with a user 
func InitRoom(user *User) *Room {
	newRoom :=  &Room{
		users: make(map [*User]bool),
		messages: make(chan []byte),
	}

	newRoom.users[user] = true
	return newRoom
}

// returns false so it can be removed from TBDDDD how this works 
func (room *Room) StartRoon() bool {
	for {
		select {
		case user := <-room.register:
			room.users[user] = true
			log.Println("New user added")
		case user := <- room.deregester:
			room.users[user] = false
			log.Println("User left the room")
			if !room.checkRoom() {
				return false
			}
		case message := <- room.messages:
			log.Println(message)
			// this is where we need to send the message to everyone in the room essentially
		}
	}
}

// this function essentially checks if there are any users left in the room and returns false if there aren't any
func (room *Room) checkRoom() bool {

	for _, val := range room.users {
		if val != true {
			return false
		}
	}
	return true
}
