package server

import (
	"ServerStatus/cmd"
	"ServerStatus/config"
	"ServerStatus/msg"
	"ServerStatus/utils"
	"bytes"
	"fmt"
	"github.com/panjf2000/gnet"
	"log"
	"sync"
)

var (
	conf         *config.Config
	authorizes   = make(map[string]string, 0)
	richNodeList = make(map[string]*msg.RichNode, 0)
)

type echoServer struct {
	*gnet.EventServer
	sockets sync.Map
}

func (es *echoServer) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	log.Printf("[OPEN ] socket with address: %s\n", c.RemoteAddr().String())

	es.sockets.Store(c.RemoteAddr().String(), c)
	authorizes[c.RemoteAddr().String()] = ""
	out = msg.Write(msg.AuthorizeMessage)

	return
}

func (es *echoServer) OnClosed(c gnet.Conn, _ error) (action gnet.Action) {
	log.Printf("[CLOSE] socket with address: %s\n", c.RemoteAddr().String())

	es.sockets.Delete(c.RemoteAddr().String())

	if id, ok := authorizes[c.RemoteAddr().String()]; ok {
		if id != "" {
			richNodeList[id].Reset()
			richNodeList[id].Online = false
			log.Println("[CLOSE] 关闭链接", richNodeList[id].Online)
		}
	}

	return
}

func (es *echoServer) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Printf("TCP server is listening on %s (multi-cores: %t, loops: %d)\n",
		srv.Addr.String(), srv.Multicore, srv.NumEventLoop)
	return
}

func (es *echoServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	buf := bytes.NewBuffer(frame)
	t, err := utils.TrimLine(buf)
	if err != nil {
		_ = c.Close()
		return
	}

	if len(t) > 0 {
		switch t[0] {
		case msg.AuthorizeMessage:
			id, err := utils.TrimLine(buf)
			if err != nil {
				_ = c.Close()
				return
			}

			strId := string(id[:])
			if node, ok := config.GetConf(strId); ok {
				m := node.(map[string]interface{})
				enable := m["enable"].(bool)
				if enable {
					authorizes[c.RemoteAddr().String()] = strId
					richNodeList[strId].Online = true
					out = msg.Write(msg.SuccessAuthorizeMessage)
				} else {
					out = msg.Write(msg.NotEnableFailMessage)
				}
			} else {
				out = msg.Write(msg.NotExistFailMessage)
			}
		case msg.ReceiveMessage:
			id, err := utils.TrimLine(buf)
			if err != nil {
				_ = c.Close()
				return
			}

			localId, ok := authorizes[c.RemoteAddr().String()]
			if !ok {
				_ = c.Close()
				return
			}

			if localId == "" {
				_ = c.Close()
				return
			}

			strId := string(id[:])
			if localId != strId {
				_ = c.Close()
				return
			}

			sys, err := utils.TrimLine(buf)
			if err != nil {
				_ = c.Close()
				return
			}

			richNodeList[strId].SystemInfo.Set(sys)
		case msg.HeartbeatMessage:
			out = msg.Write(msg.HeartbeatMessage, "pong")
		case msg.CloseMessage:
			//主动关闭
			_ = c.Close()
		default:
			//不明来历链接全关咯
			_ = c.Close()
		}
	}

	return
}

func response(m string) (data []byte) {
	for i, _ := range conf.Data {
		m := conf.Data[i].(map[string]interface{})
		id := m["id"].(string)

		if _, ok := richNodeList[id]; !ok {
			node := &msg.Node{
				Id:       m["id"].(string),
				Name:     m["name"].(string),
				Location: m["location"].(string),
				Enable:   m["enable"].(bool),
				Region:   m["region"].(string),
			}

			richNode := msg.NewRichNode(node, &msg.SystemInfo{}, false)
			richNodeList[node.Id] = richNode
		}
	}

	r := msg.NewResponse(m, richNodeList)
	data, _ = r.Json()
	return
}

func Run(p *cmd.Cmd) {
	conf = config.NewConfig(p.Filename, make([]interface{}, 0))
	echo := &echoServer{}

	go func() {
		log.Fatal(gnet.Serve(echo,
			fmt.Sprintf("tcp://%s:%d", p.Host, p.Port),
			gnet.WithMulticore(p.Multicore)))
	}()

	httpRun(p)
}
