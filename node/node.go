package node

import (
	"fmt"
	"github.com/wvanlit/mini-iota/connection"
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
	fmt.Println("Starting Node")
	n.hub.Run()
}
