package system

import (
	"fmt"
	"time"

	"github.com/flxxyz/ServerStatus/cmd"
	"github.com/flxxyz/ServerStatus/msg"
)

func Run(p *cmd.Cmd) {
	_ = p

	sysInfo := msg.NewSystemInfo()
	go sysInfo.CheckIPvNSupport()
	go sysInfo.GetTraffic()
	for {
		sysInfo.Update()
		time.Sleep(time.Second)
		data, _ := sysInfo.JsonFormat("", "    ")
		fmt.Printf("\u001B[2J\u001B[0;0H%s\n", string(data[:]))
	}
}
