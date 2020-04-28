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
	//取出消息类型
	t, err := utils.TrimLine(buf)
	if err != nil {
		_ = c.Close()
		return
	}

	if len(t) > 0 {
		closer := true //控制关闭

		switch t[0] {
		case msg.AuthorizeMessage:
			if id, err := utils.TrimLine(buf); err == nil {
				strId := string(id[:])
				if node, ok := conf.Get(strId); ok {
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

				closer = false
			}
		case msg.ReceiveMessage:
			if localId, ok := authorizes[c.RemoteAddr().String()]; ok {
				if localId != "" {
					//取出id
					if id, err := utils.TrimLine(buf); err == nil {
						strId := string(id[:])
						if localId == strId {
							//取出数据
							if sys, err := utils.TrimLine(buf); err == nil {
								richNodeList[strId].SystemInfo.Set(sys)
								closer = false
							}
						}
					}
					closer = false
				}
			}
		case msg.HeartbeatMessage:
			localId, ok := authorizes[c.RemoteAddr().String()]
			if ok {
				if localId != "" {
					if id, err := utils.TrimLine(buf); err == nil {
						strId := string(id[:])
						if localId == strId {
							out = msg.Write(msg.HeartbeatMessage, "pong")
							closer = false
						}
					}
				}
			}
		case msg.CloseMessage:
			//主动关闭
			closer = true
		default:
			//不明来历链接全关咯
			closer = true
		}

		if closer {
			_ = c.Close()
			return
		}
	}

	return
}

func response(m string) (rsp []byte) {
	data := conf.GetData()
	for i, _ := range data {
		m := data[i].(map[string]interface{})
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
	rsp, _ = r.Json()
	return
}

func Run(p *cmd.Cmd) {
	conf = config.NewConfig(p.Filename, make([]interface{}, 0))
	echo := &echoServer{}

	_ = response("init")

	go func() {
		log.Fatal(gnet.Serve(echo,
			fmt.Sprintf("tcp://%s:%d", p.Host, p.Port),
			gnet.WithMulticore(p.Multicore)))
	}()

	httpRun(p)
}
