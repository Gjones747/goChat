package socket

import (
	"crypto/sha1"
	"encoding/base64"
)

func handShakeKey(key string) string {
	const magicGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	shaWriter := sha1.New()

	shaWriter.Write([]byte(key + magicGUID))
	return base64.StdEncoding.EncodeToString(shaWriter.Sum(nil))
}
