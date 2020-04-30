package main

import (
	"ServerStatus/cmd"
	"ServerStatus/cmd/client"
	"ServerStatus/cmd/server"
	"ServerStatus/cmd/system"
	"ServerStatus/cmd/traffic"
	"ServerStatus/cmd/uuid"
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
