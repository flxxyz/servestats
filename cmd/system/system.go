package system

import (
	"ServerStatus/cmd"
	"ServerStatus/msg"
	"fmt"
	"time"
)

var sys *msg.SystemInfo

func init() {
	sys = &msg.SystemInfo{}
}

func Run(p *cmd.Cmd) {
	_ = p

	go sys.GetTraffic()
	for {
		sys.Update()
		time.Sleep(time.Second)
		data, _ := sys.Json()
		fmt.Printf("\u001B[2J\u001B[0;0H%s\n", string(data[:]))
	}
}
