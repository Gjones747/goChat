package api

import "time"

// since message data structure is going to be shared between the server and the client it will go here in api


type Message struct {
	SenderName string
	DateTime time.Time
	Contents []byte
}


func NewMessage(SenderUserName string, messageContents []byte) *Message {
	currentTime := time.Now()
	return &Message{
		SenderName: SenderUserName,
		DateTime: currentTime,
		Contents: messageContents,
	}

}
