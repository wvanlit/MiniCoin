/**
 * @Author: Wessel van Lit
 * @Project: minicoin
 * @Date: 25-May-2021
 */

package main

import (
	"os"

	"github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"

	"github.com/wvanlit/minicoin/minicoin/connection"
	"github.com/wvanlit/minicoin/minicoin/node"
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
	log.Info("Started MiniCoin")

	hub := connection.NewConnectionHub(":8080")

	myNode := node.New(hub)
	myNode.Start()
}
