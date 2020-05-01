package traffic

import (
	"fmt"
	"github.com/flxxyz/ServerStatus/cmd"
	"github.com/flxxyz/ServerStatus/msg"
	"github.com/flxxyz/ServerStatus/utils"
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
		time.Sleep(time.Second)
		fmt.Printf("\033[2J\033[0;0HTotal Rx: \033[44;37m %s \033[0m\n", utils.FileSize(sys.TrafficRxTotal))
		fmt.Printf("Rx/s: %4s\033[42;37m %s/s \033[0m\n", "", utils.FileSize(sys.TrafficRx))
		fmt.Printf("Total Tx: \033[44;37m %s \033[0m \n", utils.FileSize(sys.TrafficTxTotal))
		fmt.Printf("Tx/s: %4s\033[42;37m %s/s \033[0m\n-------------------------\n监听实时流量中... [Ctrl+C 退出]", "", utils.FileSize(sys.TrafficTx))
	}
}
