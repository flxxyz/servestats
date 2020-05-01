package system

import (
	"fmt"
	"github.com/flxxyz/ServerStatus/cmd"
	"github.com/flxxyz/ServerStatus/msg"
	"time"
)

func Run(p *cmd.Cmd) {
	_ = p

	sys := msg.NewSystemInfo(true)
	go sys.GetTraffic()
	for {
		sys.Update()
		time.Sleep(time.Second)
		data, _ := sys.Json()
		fmt.Printf("\u001B[2J\u001B[0;0H%s\n", string(data[:]))
	}
}
