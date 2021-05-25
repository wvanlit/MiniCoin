/**
 * @Author: Wessel van Lit
 * @Project: minicoin
 * @Date: 25-May-2021
 */

package connection

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	id       string
	conn     *websocket.Conn
	isActive bool
}

func createClient(id string, conn *websocket.Conn) Client {
	return Client{
		id:   id,
		conn: conn,
	}
}

func (c *Client) handleMessages() {
	outgoing := make(chan []byte, 8)
	incoming := make(chan []byte, 8)
	c.isActive = true

	go func() {
		c.ReadIncomingMessages(incoming)
	}()

	go func() {
		c.WriteOutgoingMessages(outgoing)
	}()

	for {
		select {
		case msg := <-incoming:
			outgoing <- msg
		}
	}
}

func (c *Client) WriteOutgoingMessages(outgoing chan []byte) {
	for c.isActive {
		message := <-outgoing
		err := c.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
	log.Infof("Shutting down Read on %s", c.conn.RemoteAddr())
}

func (c *Client) ReadIncomingMessages(incoming chan []byte) {
	for c.isActive {
		mt, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		if mt > 0 {
			incoming <- message
		}
	}
	log.Infof("Shutting down Write on %s", c.conn.RemoteAddr())
}
