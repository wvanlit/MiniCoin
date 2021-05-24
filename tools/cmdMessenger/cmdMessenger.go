package main

import (
	"bufio"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
	"github.com/wvanlit/mini-btc/messages"
	"net/url"
	"os"
	"strings"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		DisableQuote:    false,
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
	})
	log.SetOutput(colorable.NewColorableStdout())
	log.SetOutput(os.Stdout)

	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("Recv: '%s'", message)
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	log.Println("WebSocket Messenger")
	for {
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)

		byteText := checkShorthands(text)
		log.Infof("Sending:%s", string(byteText))
		err := c.WriteMessage(websocket.TextMessage, byteText)
		if err != nil {
			log.Println("write:", err)
			return
		}
	}
}

func checkShorthands(text string) []byte {
	if !strings.HasPrefix(text, ":") {
		return []byte(text)
	}

	fullCommand := text[1:]
	commandParts := strings.Split(fullCommand, "|")
	command := commandParts[0]

	switch command {
	case "auth":
		msg := messages.CreateAuthenticationRequestMessage(commandParts[1])
		bytesMsg, err := json.Marshal(msg)
		if err != nil {
			log.Fatal(err)
		}
		return bytesMsg
	case "data":
		msg := messages.CreateDataMessage(commandParts[1])
		bytesMsg, err := json.Marshal(msg)
		if err != nil {
			log.Fatal(err)
		}
		return bytesMsg
	default:
		return []byte(text)
	}
}
