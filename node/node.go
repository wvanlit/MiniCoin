package node

import (
	log "github.com/sirupsen/logrus"
	"github.com/wvanlit/mini-btc/connection"
)

type Node struct {
	hub *connection.Hub
}

func New(hub *connection.Hub) *Node {
	return &Node{
		hub: hub,
	}
}

func (n Node) Start() {
	log.Info("Starting Node")
	n.hub.Run()
}
