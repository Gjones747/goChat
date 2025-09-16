package models

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Gjones747/goChat/api"
)

type Room struct {
	roomCode   string
	users      map[*User]bool
	messages   chan api.Envelope
	register   chan *User
	deregester chan *User

	hub *RoomHub
}

// each new room has to start with a user
func InitRoom(hub *RoomHub, code string) *Room {
	newRoom := &Room{
		roomCode:   code,
		users:      make(map[*User]bool),
		messages:   make(chan api.Envelope, 10),
		register:   make(chan *User),
		deregester: make(chan *User),

		hub: hub,
	}

	return newRoom
}

func (room *Room) sendServerPing() {
	pingMessage := api.NewMessage("SYS", "", nil, []byte("System ping"))
	defer log.Println("Closing ping cycle for " + room.roomCode)
	for {
		time.Sleep(30 * time.Second)

		for user, isInRoom := range room.users {
			if isInRoom {
				user.Recieve(pingMessage)
				log.Println("Pinged users in " + room.roomCode)
			}
		}
	}
}

func (room *Room) StartRoom() bool {
	defer log.Println("Closing room " + room.roomCode)
	log.Println("starting room " + room.roomCode)
	go room.sendServerPing()
	for {
		select {
		case user := <-room.register:
			room.users[user] = true
			go user.UserIOReader()
			go user.userIOWriter()
			fmt.Println("New user added")
			room.sendUserJoinMessage(user.Name)
			room.sendUserList()
		case user := <-room.deregester:
			name := user.Name
			user.Connection.Close()
			close(user.send)
			room.users[user] = false
			delete(room.users, user)
			log.Printf("User: %s left the room\n", user.Name)
			if !room.checkRoom() {
				delete(room.hub.Rooms, room.roomCode)
				if _, ok := room.hub.Rooms[room.roomCode]; ok {
					log.Printf("delete failed")
				}
				return false
			}
			room.sendUserLeaveMessage(name)
			room.sendUserList()
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
func (room *Room) sendUserLeaveMessage(userName string) {
	message := api.NewMessage("JOIN_LEAVE", "", nil, []byte(fmt.Sprintf("%s has left the room!", userName)))

	for user, isInRoom := range room.users {
		if isInRoom {
			user.Recieve(message)
		}
	}
}

func (room *Room) sendUserJoinMessage(userName string) {
	message := api.NewMessage("JOIN_LEAVE", "", nil, []byte(fmt.Sprintf("%s has joined the room!", userName)))

	for user, isInRoom := range room.users {
		if isInRoom {
			user.Recieve(message)
		}
	}
}

func (room *Room) sendUserList() {
	userList := [][]string{}

	for user, inRoom := range room.users {
		if inRoom {
			userData := [][]string{{user.Name, string(user.SessionID)}}
			userList = append(userList, userData...)
		}
	}

	userListEnvelope, err := api.NewUserList(userList)

	if err != nil {
		log.Println(err)
	}

	room.sendMessage(userListEnvelope)
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

func (room *Room) AddMessage(message []byte) {
	var messageJson api.Envelope
	err := json.Unmarshal(message, &messageJson)

	if err != nil {
		log.Println(err)
	}

	room.messages <- messageJson
}
