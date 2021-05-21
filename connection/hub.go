package connection

import (
	"github.com/gorilla/websocket"
	"log"
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
	c, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	// Authenticate
	isAuthenticated := false
	for !isAuthenticated {

	}
	// Create Client

	h.handleMessages(c)
}

func (h *Hub) handleMessages(conn *websocket.Conn) {
	outgoing := make(chan []byte, 8)
	incoming := make(chan []byte, 8)
	isActive := true

	go func() {
		for isActive {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			if mt > 0 {
				incoming <- message
			}
		}
		log.Println("Shutting down Write")
	}()

	go func() {
		for isActive {
			message := <-outgoing
			err := conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
		log.Println("Shutting down Read")
	}()

	for {
		select {
		case msg := <-incoming:
			outgoing <- msg
		}
	}
}
