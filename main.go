package main

import (
	"fmt"
	"github.com/wvanlit/mini-iota/connection"
	"github.com/wvanlit/mini-iota/node"
)

func main() {
	fmt.Println("Hello World!")

	hub := connection.NewConnectionHub(":8080")

	myNode := node.New(hub)
	myNode.Start()
}
