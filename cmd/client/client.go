package client

import (
	"fmt"
	"log"
	"net/rpc"
	"time"

	"github.com/flxxyz/ServerStatus/cmd"
	"github.com/flxxyz/ServerStatus/config"
	"github.com/flxxyz/ServerStatus/msg"
)

type RpcClient struct {
	sentChan  chan bool
	closeChan chan bool
	conn      *rpc.Client
	err       error
	cmd       *cmd.Cmd
	msg       *msg.RpcClientMessage
	reply     msg.RpcServerReply
}

func (rc *RpcClient) call() {
	defer func() {
		_ = rc.conn.Close()
		rc.sentChan <- true
	}()

	rc.err = rc.conn.Call("MessageService.Report", rc.msg, &rc.reply)
	if rc.err != nil {
		log.Printf("[%s] %s: %s\n", msg.ClientReportType, msg.ClientErrorMessage, rc.err.Error())
		return
	}
	log.Printf("[%s] %s: %s\n", msg.ClientReportType, msg.ClientReplyMessage, string(rc.reply))
}

func (rc *RpcClient) report() {
	rc.conn, rc.err = rpc.DialHTTP("tcp", fmt.Sprintf("%s:%d", rc.cmd.Host, rc.cmd.Port))
	if rc.err != nil {
		log.Printf("[%s] %s: %s后尝试重新连接", msg.ClientConnectType, msg.ClientFailMessage, config.IvReconnect.String())
		rc.closeChan <- true
	} else {
		rc.msg.Update()
		rc.call()
	}
}

func doConn(rc *RpcClient) {
	rc.sentChan <- true

	for {
		select {
		case <-rc.sentChan:
			time.Sleep(rc.cmd.Interval)
			rc.report()
		case <-rc.closeChan:
			time.Sleep(config.IvReconnect)
			rc.report()
		}
	}
}

func Run(p *cmd.Cmd) {
	if p.Interval <= 0 {
		p.Interval = config.IvSent //限制发送间隔要大于等于一秒
	}

	rpcClient := &RpcClient{
		sentChan:  make(chan bool, 1),
		closeChan: make(chan bool, 1),
		cmd:       p,
		msg: &msg.RpcClientMessage{
			Id:         p.Id,
			SystemInfo: msg.NewSystemInfo(),
		},
		reply: make([]byte, 0),
	}
	go rpcClient.msg.CheckIPvNSupport()
	go rpcClient.msg.GetTraffic()
	go doConn(rpcClient)

	c := make(chan bool)
	<-c
}
