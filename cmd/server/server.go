package server

import (
	"fmt"
	"github.com/flxxyz/ServerStatus/msg"
	"runtime"

	"github.com/flxxyz/ServerStatus/cmd"
	"github.com/flxxyz/ServerStatus/config"
)

var (
	response *msg.Response
)

func Run(p *cmd.Cmd) {
	config.NewConfig(p.Filename)

	if p.Multicore {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	addr := fmt.Sprintf("%s:%d", p.Host, p.Port)
	go RPCServer(addr)
	HTTPServer(addr)
}

func init() {
	response = &msg.Response{
		Message: "init",
		Servers: make([]*msg.ResponseNode, 0),
		Nodes:   make(map[string]*msg.ResponseNode, 0),
	}
}
