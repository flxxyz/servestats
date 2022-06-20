package main

import (
	"github.com/flxxyz/servestats/cmd"
	"github.com/flxxyz/servestats/cmd/client"
	"github.com/flxxyz/servestats/cmd/server"
	"github.com/flxxyz/servestats/cmd/system"
	"github.com/flxxyz/servestats/cmd/traffic"
	"github.com/flxxyz/servestats/cmd/uuid"
)

func main() {
	c := cmd.Run()

	switch c.T {
	case "server":
		server.Run(c)
	case "client":
		client.Run(c)
	case "uuid":
		uuid.Run(c)
	case "system":
		system.Run(c)
	case "traffic":
		traffic.Run(c)
	}
}
