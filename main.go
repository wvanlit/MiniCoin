package main

import (
	"github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
	"github.com/wvanlit/mini-iota/connection"
	"github.com/wvanlit/mini-iota/node"
	"os"
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
	log.Info("Started Mini-IOTA")

	hub := connection.NewConnectionHub(":8080")

	myNode := node.New(hub)
	myNode.Start()
}
