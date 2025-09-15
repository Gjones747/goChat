package socketcontroller

import (
	"bufio"
	"crypto/rand"
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
)

func SendUpgrade(url url.URL) (net.Conn, error) {

	conn, err := tls.Dial("tcp", "kaolun.site:443", &tls.Config{
		ServerName: "kaolun.site",
	})

	if err != nil {
		log.Printf("Failed to connect to: %s", url.Host)
		log.Fatal(err)
		return nil, err
	}

	keyBytes := make([]byte, 16)
	if _, err = rand.Read(keyBytes); err != nil {
		return nil, err
	}
	clientKey := base64.StdEncoding.EncodeToString(keyBytes)

	requestPathWithQuery, err := querySetter(url)
	if err != nil {
		return nil, err
	}

	upgradeRequest := fmt.Sprintf("GET %s HTTP/1.1\r\n", requestPathWithQuery) +
		fmt.Sprintf("Host: %s\r\n", url.Host) +
		"Upgrade: websocket\r\n" +
		"Connection: Upgrade\r\n" +
		fmt.Sprintf("Sec-WebSocket-Key: %s\r\n", clientKey) +
		"Sec-WebSocket-Version: 13\r\n" +
		"Origin: http://localHost\r\n" +
		"\r\n"

	_, err = conn.Write([]byte(upgradeRequest))

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	reader := bufio.NewReader(conn)

	resp, err := http.ReadResponse(reader, &http.Request{Method: "GET"})
	if err != nil {
		log.Println("Failed to read server respone")
		log.Fatal(err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusSwitchingProtocols {
		log.Fatal(resp.StatusCode)
		return nil, errors.New("Did not Recieve a switching protical response")
	}

	serverKey := resp.Header.Get("Sec-WebSocket-Accept")

	if !keyCheck(clientKey, serverKey) {
		log.Fatal("keycheck")
		return nil, errors.New("Did not recieve the proper accept token")
	}

	// at this point the websocket connection has been established and the connection can be pased on
	return conn, nil
}

func keyCheck(clientKey string, serverKey string) bool {
	clientKey = fmt.Sprintf("%s258EAFA5-E914-47DA-95CA-C5AB0DC85B11", clientKey)

	hasher := sha1.New()
	hasher.Write([]byte(clientKey))

	signedClient := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

	if signedClient == serverKey {
		return true
	}

	return false
}

func querySetter(url url.URL) (string, error) {
	if url.RawQuery != "" {
		return url.Path + "?" + url.RawQuery, nil
	}

	return "", errors.New("Couldn't join url")
}
