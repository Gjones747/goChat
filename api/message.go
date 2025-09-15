package api

import "time"

// since message data structure is going to be shared between the server and the client it will go here in api
// netID is going to be a more unique identifier for the users session
type Message struct {
	Type string // can either by SYS or USER 
	SessionID []byte
	SenderName string
	DateTime   time.Time
	Contents   []byte
}
