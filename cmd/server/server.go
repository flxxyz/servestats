package server

import (
	"fmt"
	"github.com/flxxyz/servestats/msg"
	"runtime"

	"github.com/flxxyz/servestats/cmd"
	"github.com/flxxyz/servestats/config"
)

var (
	response *msg.Response
)

func Run(p *cmd.Cmd) {
	config.NewConfig(p.Filename)

	if p.Multicore {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	go RPCServer(fmt.Sprintf("%s:%d", p.Host, p.Port))
	HTTPServer(fmt.Sprintf("%s:%d", p.Host, p.HTTPPort))
}

func init() {
	response = &msg.Response{
		Message: "init",
		Servers: make([]*msg.ResponseNode, 0),
		Nodes:   make(map[string]*msg.ResponseNode, 0),
	}
}
