package client

import (
	"ServerStatus/cmd"
	"ServerStatus/config"
	"ServerStatus/msg"
	"ServerStatus/timer"
	"bytes"
	"fmt"
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
	callback := func() {
		write(msg.Write(msg.HeartbeatMessage, params.Id))
	}

	timer.New(callback, interval)
}

//发送通过验证的消息
func sent(interval time.Duration) {
	callback := func() {
		sys.Update()
		packet, _ := sys.Json()
		write(msg.Write(msg.ReceiveMessage, params.Id, packet))
	}

	timer.New(callback, interval)
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

	sys = &msg.SystemInfo{HasConvStr: p.HasConvStr}
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

		t, _ := buffer.ReadByte()
		switch t {
		case msg.AuthorizeMessage:
			auth()
		case msg.SuccessAuthorizeMessage:
			log.Println("[AUTHORIZE]", "success")
			sent(p.Interval)
			heartbeat(config.IntervalHeartbeat)
		case msg.HeartbeatMessage:
			log.Println("[HEARTBEAT]")
		case msg.NotExistFailMessage:
			log.Println("[FAIL]", "This node is not exist")
			closeSent()
		case msg.NotEnableFailMessage:
			log.Println("[FAIL]", "This node is not enable")
			closeSent()
		case msg.CloseMessage:
			log.Println("[CLOSE]")
			closeSent()
		}

		copy(buf, emptyBuf)
		buffer.Reset()
	}
}
