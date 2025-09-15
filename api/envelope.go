package api

import (
	"encoding/json"
	"time"
)

type Envelope struct {
	DataType string
	Data     json.RawMessage
}

func NewMessage(messageType string, senderUserName string, sessionID []byte, messageContents []byte) Envelope {
	currentTime := time.Now()
	message := Message{

		Type: messageType,
		SenderName: senderUserName,
		SessionID: sessionID,
		DateTime:   currentTime,
		Contents:   messageContents,
	}

	messageJSON, _ := json.Marshal(message)

	return Envelope{
		DataType: "message",
		Data:     messageJSON,
	}

}

//userlist just contains a list of usernames
func NewUserList(users [][]string) (Envelope, error) {
	userList := UserList {
		Users: users,
	}


	NewUserListJson, err := json.Marshal(userList) 
	if err != nil {
		return Envelope{}, err
	}

	return Envelope{
		DataType: "userList",
		Data: NewUserListJson,
	}, nil
}
