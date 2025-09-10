package models

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Gjones747/goChat/api"
)

type Room struct {
	users      map[*User]bool
	messages   chan api.Envelope
	register   chan *User
	deregester chan *User
}

// each new room has to start with a user
func InitRoom() *Room {
	newRoom := &Room{
		users:      make(map[*User]bool),
		messages:   make(chan api.Envelope, 10),
		register:   make(chan *User),
		deregester: make(chan *User),
	}

	return newRoom
}

func (room *Room) StartRoom() bool {
	log.Println("starting room")
	for {
		select {
		case user := <-room.register:
			room.users[user] = true
			go user.UserIOReader()
			go user.userIOWriter()
			fmt.Println("New user added")
		case user := <-room.deregester:
			user.Connection.Close()
			close(user.send)
			room.users[user] = false
			log.Printf("User: %s left the room\n", user.Name)
			if !room.checkRoom() {
				return false
			}
		case message := <-room.messages:
			log.Printf("messege: %s", message)
			// this is where we need to send the message to everyone in the room essentially
			room.sendMessage(message)
		}
	}
}

// this functions sends the users message to each person in the room with them
func (room *Room) sendMessage(message api.Envelope) {
	for key, val := range room.users {
		if val {
			key.Recieve(message)
		}
	}
}

// this function essentially checks if there are any users left in the room and returns false if there aren't any
func (room *Room) checkRoom() bool {
	for _, val := range room.users {
		if val {
			return true
		}
	}
	return false
}

// this function will add someone to a room using the register channel
func (room *Room) AddUser(user *User) {
	room.register <- user
}

func (room *Room) RemoveUser(user *User) {
	room.deregester <- user
}

func (room *Room) AddMessage(message []byte)  {
	var messageJson api.Envelope
	err := json.Unmarshal(message, &messageJson)
    
	if err != nil {
		log.Println(err)
	}

	room.messages <- messageJson
}
