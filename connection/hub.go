package connection

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/wvanlit/mini-btc/messages"
	"net"
	"net/http"
	"time"
)

type Hub struct {
	addr      string
	upgrader  *websocket.Upgrader
	clients   map[*Client]Client
	broadcast chan []byte
	shutdown  chan bool
	register  chan *Client
}

func NewConnectionHub(addr string) *Hub {
	upgrader := &websocket.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	return &Hub{
		addr:      addr,
		upgrader:  upgrader,
		clients:   make(map[*Client]Client),
		broadcast: make(chan []byte),
		shutdown:  make(chan bool),
		register:  make(chan *Client, 4),
	}
}

func (h *Hub) Run() {
	h.SetupListener()
	isRunning := true

	log.Println("Entered Main Loop")
	for isRunning {
		select {
		case c := <-h.register:
			log.Println("Registering Client", c)
			h.clients[c] = *c
		case <-h.shutdown:
			log.Println("Shutting down")
			isRunning = false
			break
		case <-time.After(500 * time.Millisecond):
		}
	}
}

func (h *Hub) SetupListener() {
	go func() { log.Fatal(http.ListenAndServe(h.addr, nil)) }()

	http.HandleFunc("/ws", h.connect)
}

func (h *Hub) connect(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("Error on upgrading to websocket: %s", err)
		return
	}

	log.Infof("Started connection with %s", conn.RemoteAddr())

	// Authenticate
	auth := h.HandleAuthentication(conn)
	// Create Client
	if auth == nil {
		log.Errorf("Closing connection with %s because of errors", conn.RemoteAddr())
		conn.Close()
		return
	}
	log.Infof("Authentication on id: %s by %s", auth.Id, conn.RemoteAddr())

	client := createClient(auth.Id, conn)
	h.clients[&client] = client

	msg := messages.CreateAuthenticationResultMessage(true, "")
	msgJson, err := json.Marshal(msg) // err is unnecessary, as we created it using a pre-defined interface

	err = conn.WriteMessage(websocket.TextMessage, msgJson)
	if err != nil {
		log.Warnf("Error on writing message to %s, err: %s", conn.RemoteAddr(), err)
	}

	client.handleMessages()
}

func (h *Hub) HandleAuthentication(c *websocket.Conn) *messages.AuthenticationRequest {
	var auth *messages.AuthenticationRequest = nil
	for auth == nil {
		msg, err := ReceiveMessage(c)
		if err != nil {
			switch err.(type) {
			case *net.OpError: // This prevents crashing upon unexpected close
				return nil
			default:
				continue
			}
		}

		if msg.MsgType != messages.AUTHENTICATION_REQUEST {
			msg := messages.CreateErrorMessage(messages.NOT_AUTHENTICATED, "You are not authenticated yet")
			MarshalAndSendMessage(msg, c)
			continue
		}

		payload, err := messages.UnmarshalAuthenticationRequest(msg.Payload.(map[string]interface{}))
		if err != nil {
			log.Warnln("Cannot unmarshal message to Auth request:", msg.Payload)
			continue
		}

		found := false
		for client := range h.clients {
			if client.id == auth.Id {
				found = true
				break
			}
		}

		if !found {
			auth = payload
		} else {
			msg := messages.CreateAuthenticationResultMessage(false, "ID already in use")
			MarshalAndSendMessage(msg, c)
		}
	}
	return auth
}

func MarshalAndSendMessage(msg messages.Message, c *websocket.Conn) {
	msgJson, _ := json.Marshal(msg) // err is unnecessary, as we created it using a pre-defined interface

	err := c.WriteMessage(websocket.TextMessage, msgJson)
	if err != nil {
		log.Warnf("Error on writing message to %s, err: %s", c.RemoteAddr(), err)
	}
}

func ReceiveMessage(c *websocket.Conn) (*messages.Message, error) {
	mt, message, err := c.ReadMessage()

	if err != nil || mt <= 0 {
		log.Warnf("Error on read message from %s, err: %s", c.RemoteAddr(), err)

		msg := messages.CreateErrorMessage(messages.INVALID_MESSAGE, "Could not read message")
		msgJson, _ := json.Marshal(msg) // err is unnecessary, as we created it using a pre-defined interface

		err2 := c.WriteMessage(websocket.TextMessage, msgJson)
		if err2 != nil {
			log.Errorf("Error on writing message to %s, err: %s", c.RemoteAddr(), err)
			return nil, err2
		}

		return nil, err
	}

	msg, err := messages.UnmarshalMessageJSON(message)

	if err != nil {
		log.Warnf("Error on unmarshal message from %s, err: %s", c.RemoteAddr(), err)

		msg := messages.CreateErrorMessage(messages.INVALID_FORMAT, "Could not parse message")
		MarshalAndSendMessage(msg, c)

		return nil, err
	}
	return msg, nil
}
