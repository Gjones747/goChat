package api

import (
	"crypto/rand"
	"encoding/binary"
	"time"
)

func MakeSessionID() ([]byte, error) {
	timeStamp := time.Now().UnixNano()
	timeStampBytes := make([]byte, 8)

	binary.BigEndian.PutUint64(timeStampBytes, uint64(timeStamp))

	randomBytes := make([]byte, 8)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	sessionID := append(timeStampBytes, randomBytes...)
	return sessionID, nil
}
