package server

import (
    "ServerStatus/cmd"
    "ServerStatus/config"
    "ServerStatus/msg"
    "fmt"
    jsoniter "github.com/json-iterator/go"
    "github.com/panjf2000/gnet"
    "log"
    "time"
)

var (
    json = jsoniter.ConfigCompatibleWithStandardLibrary
    data = make([]msg.Node, 0)

    authorizeNodes = make(map[string]bool, 0)
    richNodes      = make(map[string][]byte, 0)
)

type echoServer struct {
    *gnet.EventServer
    tick time.Duration
}

func (es *echoServer) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
    log.Printf("[OPEN ] socket with address: %s\n", c.RemoteAddr().String())
    authorizeNodes[c.RemoteAddr().String()] = false
    out = msg.Write(msg.AuthorizeMessage)
    return
}

func (es *echoServer) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
    log.Printf("[CLOSE] socket with address: %s\n", c.RemoteAddr().String())
    if authorizeNodes[c.RemoteAddr().String()] {
        delete(richNodes, c.RemoteAddr().String())
    }
    delete(authorizeNodes, c.RemoteAddr().String())
    return
}

func (es *echoServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
    switch frame[0] {
    case msg.AuthorizeMessage:
        id := string(frame[1:37])
        if _, ok := config.GetConf(id); ok {
            authorizeNodes[c.RemoteAddr().String()] = true
            out = msg.Write(msg.SuccessAuthorizeMessage, id)
        } else {
            out = msg.Write(msg.FailAuthorizeMessage)
        }
    case msg.ReceiveMessage:
        richNodes[c.RemoteAddr().String()] = frame[37:]
    case msg.PingMessage:
        out = msg.Write(msg.PongMessage, "pong")
    case msg.CloseMessage:
        //主动关闭
        _ = c.Close()
    }

    return
}

func Run(p *cmd.Cmd) {
    _ = config.NewConfig(p.Filename, data)
    echo := &echoServer{tick: p.Interval}

    log.Printf("server listening tcp://%s:%d\n", p.Host, p.Port)
    log.Fatal(gnet.Serve(echo,
        fmt.Sprintf("tcp://%s:%d", p.Host, p.Port),
        gnet.WithMulticore(p.Multicore)))
}
