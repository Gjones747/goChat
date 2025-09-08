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
