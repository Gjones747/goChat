package api

import (
	"encoding/json"
	"time"
)

type Envelope struct {
	DataType string
	Data     json.RawMessage
}

func NewMessage(SenderUserName string, messageContents []byte) Envelope {
	currentTime := time.Now()
	message := Message{
		SenderName: SenderUserName,
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
func NewUserList(users []string) (Envelope, error) {
	userList := UserList {
		users: users,
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
