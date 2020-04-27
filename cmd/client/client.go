package client

import (
    "ServerStatus/cmd"
    "ServerStatus/cmd/system"
    "ServerStatus/config"
    "ServerStatus/msg"
    "bytes"
    "fmt"
    "log"
    "net"
    "os"
    "time"
)

var (
    params   *cmd.Cmd
    buffer   = bytes.NewBuffer(make([]byte, 0))
    emptyBuf = make([]byte, 4096)
    buf      = make([]byte, 4096)
)

//发送验证消息
func auth(c net.Conn, id string) {
    authMsg := bytes.NewBuffer(make([]byte, 0))
    authMsg.WriteByte(msg.AuthorizeMessage)
    authMsg.WriteString(id)

    _, err := c.Write(authMsg.Bytes())
    if err != nil {
        _ = c.Close()
    }
}

func ticker(callback func(), interval time.Duration) {
    ticker := time.NewTicker(interval)
    //defer ticker.Stop()

    for range ticker.C{
        callback()
    }
}

//发送心跳
func heartbeat(c net.Conn, interval time.Duration) {
    ticker(func() {
        if _, err := c.Write(msg.Write(msg.PingMessage, params.Id)); err != nil {
            _ = c.Close()
        }
    }, interval)
}

//发送通过验证的消息
func sent(c net.Conn, interval time.Duration) {
    sys := system.NewSystemInfo(params.Convert)

    ticker(func() {
        sys.Update()
        packet, _ := sys.Json()
        if _, err := c.Write(msg.Write(msg.ReceiveMessage, params.Id, packet)); err != nil {
            _ = c.Close()
        } else {
            //log.Printf("[SENT] %s\n", packet)
        }
    }, interval)
}

//发送关闭链接消息
func closeSent(c net.Conn) {
    if _, err := c.Write(msg.Write(msg.CloseMessage)); err != nil {
        _ = c.Close()
    }
}

func Run(p *cmd.Cmd) {
    params = p
    if p.Interval <= 0 {
        p.Interval = config.IntervalSent //限制发送间隔要大于等于一秒
    }

    conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", p.Host, p.Port))
    if err != nil {
        panic(err)
    }
    //defer conn.Close()

    for {
        if _, err := conn.Read(buf); err != nil {
            closeSent(conn)
            os.Exit(0)
        } else {
            buffer.Write(buf)
        }

        typeMessage, _ := buffer.ReadByte()
        switch typeMessage {
        case msg.AuthorizeMessage:
            auth(conn, p.Id)
        case msg.SuccessAuthorizeMessage:
            log.Println("[AUTHORIZE] success")
            go sent(conn, p.Interval)
            go heartbeat(conn, time.Second*config.IntervalHeartbeat)
        case msg.FailAuthorizeMessage:
            log.Println("[AUTHORIZE] fail")
            closeSent(conn)
        case msg.CloseMessage:
            log.Println("[CLOSE]")
            closeSent(conn)
        case msg.PongMessage:
            log.Println("[PONG]")
        }

        copy(buf, emptyBuf)
        buffer.Reset()
    }
}
