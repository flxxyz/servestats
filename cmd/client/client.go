package client

import (
	"bytes"
	"fmt"
	"github.com/flxxyz/ServerStatus/cmd"
	"github.com/flxxyz/ServerStatus/config"
	"github.com/flxxyz/ServerStatus/msg"
	"github.com/flxxyz/ServerStatus/utils"
	"log"
	"net"
	"os"
	"time"
)

var (
	params *cmd.Cmd
	sys    *msg.SystemInfo
	conn   net.Conn
	err    error

	buffer   = bytes.NewBuffer(make([]byte, 0))
	emptyBuf = make([]byte, 4096)
	buf      = make([]byte, 4096)
)

func write(msg []byte) {
	_, err := conn.Write(msg)
	if err != nil {
		_ = conn.Close()
	}
}

//发送验证消息
func auth() {
	write(msg.Write(msg.AuthorizeMessage, params.Id))
}

//发送心跳
func heartbeat(interval time.Duration) {
	t := time.NewTicker(interval)
	for range t.C {
		write(msg.Write(msg.HeartbeatMessage, params.Id))
	}
}

//发送通过验证的消息
func sent(interval time.Duration) {
	t := time.NewTicker(interval)
	for range t.C {
		sys.Update()
		packet, _ := sys.Json()
		write(msg.Write(msg.ReceiveMessage, params.Id, packet))
	}
}

//发送关闭链接消息
func closeSent() {
	write(msg.Write(msg.CloseMessage))
}

func Run(p *cmd.Cmd) {
	params = p
	if p.Interval <= 0 {
		p.Interval = config.IntervalSent //限制发送间隔要大于等于一秒
	}

	sys = msg.NewSystemInfo(p.HasConvStr)
	go sys.CheckIPvNSupport()
	go sys.GetTraffic()

	conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", p.Host, p.Port))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	for {
		if _, err := conn.Read(buf); err != nil {
			closeSent()
			os.Exit(1)
		} else {
			buffer.Write(buf)
		}

		//取出消息类型
		t, err := utils.TrimLine(buffer)
		if err != nil {
			return
		}

		switch t[0] {
		case msg.AuthorizeMessage:
			auth()
		case msg.SuccessAuthorizeMessage:
			log.Println("[AUTHORIZE]", "success")
			go sent(p.Interval)
			go heartbeat(config.IntervalHeartbeat)
		case msg.HeartbeatMessage:
			pong, _ := utils.TrimLine(buffer)
			log.Println("[HEARTBEAT]", string(pong))
		case msg.NotExistFailMessage:
			log.Println("[FAIL]", "This Node is not exist")
			closeSent()
		case msg.NotEnableFailMessage:
			log.Println("[FAIL]", "This Node is not enable")
			closeSent()
		case msg.CloseMessage:
			log.Println("[CLOSE]")
			closeSent()
		}

		copy(buf, emptyBuf)
		buffer.Reset()
	}
}
