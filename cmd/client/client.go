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
	params   *cmd.Cmd
	sys      *msg.SystemInfo
	emptyBuf []byte
	pool     []*Client
	index    = 0
)

func init() {
	emptyBuf = make([]byte, 4096)
	pool = make([]*Client, 0)
}

type Client struct {
	index int
	conn  net.Conn
	buf   *bytes.Buffer
	data  []byte
}

func (c *Client) write(msg []byte) {
	_, err := c.conn.Write(msg)
	if err != nil {
		_ = c.conn.Close()
	}
}

//发送心跳
func (c *Client) heartbeat(interval time.Duration) {
	//第一次心跳
	c.write(msg.Write(msg.HeartbeatMessage, params.Id))

	t := time.NewTicker(interval)
	for range t.C {
		c.write(msg.Write(msg.HeartbeatMessage, params.Id))
	}
}

//发送系统信息
func (c *Client) sent(interval time.Duration) {
	t := time.NewTicker(interval)
	for range t.C {
		sys.Update()
		packet, _ := sys.Json()
		c.write(msg.Write(msg.ReceiveMessage, params.Id, packet))
	}
}

//发送关闭链接消息
func (c *Client) closeSent() {
	c.write(msg.Write(msg.CloseMessage))
}

//发送关闭链接消息并且退出程序
func (c *Client) exitSent() {
	c.closeSent()
	os.Exit(1)
}

//重置缓存的空间
func (c *Client) reset() {
	copy(c.data, emptyBuf)
	c.buf.Reset()
}

func (c *Client) getMessageType() (byte, error) {
	lineBtyes, err := utils.TrimLine(c.buf)
	if err != nil {
		return 0, err
	}
	return lineBtyes[0], nil
}

func (c *Client) Task() {
	defer c.conn.Close()
	c.conn.SetReadDeadline(time.Now().Add(config.TimeoutReadDeadline))

	for {
		if _, err := c.conn.Read(c.data); err != nil {
			break
		}
		c.buf.Write(c.data)

		//取出消息类型
		messageType, err := c.getMessageType()
		if err != nil {
			break
		}

		switch messageType {
		case msg.AuthorizeMessage:
			c.write(msg.Write(msg.AuthorizeMessage, params.Id))
		case msg.SuccessAuthorizeMessage:
			log.Println("[AUTHORIZE]", "success")
			go c.sent(params.Interval)
			go c.heartbeat(config.IntervalHeartbeat)
		case msg.HeartbeatMessage:
			pong, _ := utils.TrimLine(c.buf)
			log.Println("[HEARTBEAT]", string(pong))
		case msg.NotExistFailMessage:
			log.Println("[FAIL]", "This Node is not exist")
			c.exitSent()
		case msg.NotEnableFailMessage:
			log.Println("[FAIL]", "This Node is not enable")
			c.exitSent()
		case msg.CloseMessage:
			log.Println("[CLOSE]")
			c.exitSent()
		}

		c.reset()
	}
}

func doConn() {
	for {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", params.Host, params.Port))
		if err != nil {
			log.Println("[CONNECT]", "fail", config.IntervalReconnect.String()+"后尝试重新连接")
			index++
			time.Sleep(config.IntervalReconnect)
		} else {
			log.Println("[CONNECT]", "ok")
			c := &Client{
				index: index,
				conn:  conn,
				buf:   bytes.NewBuffer(make([]byte, 0)),
				data:  make([]byte, config.MessageBufferSize),
			}
			c.Task()
		}
	}
}

func Run(p *cmd.Cmd) {
	params = p
	if p.Interval <= 0 {
		p.Interval = config.IntervalSent //限制发送间隔要大于等于一秒
	}

	sys = msg.NewSystemInfo(p.OutputText)
	go sys.CheckIPvNSupport()
	go sys.GetTraffic()
	go doConn()

	c := make(chan bool)
	<-c
}
