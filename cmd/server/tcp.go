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
	echo         *echoServer
	conf         *config.Config
	authorizes   = make(map[string]string, 0)
	richNodeList = make(map[string]*msg.RichNode, 0)
	checkNodes   = make(map[string]bool, 0)
	locker       = &sync.RWMutex{}
)

type echoServer struct {
	*gnet.EventServer
	sockets sync.Map
}

func (es *echoServer) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	log.Printf("[OPEN ] socket with address: %s\n", c.RemoteAddr().String())

	authorizes[c.RemoteAddr().String()] = ""
	out = msg.Write(msg.AuthorizeMessage)

	return
}

func (es *echoServer) OnClosed(c gnet.Conn, _ error) (action gnet.Action) {
	log.Printf("[CLOSE] socket with address: %s\n", c.RemoteAddr().String())

	if id, ok := authorizes[c.RemoteAddr().String()]; ok {
		if id != "" {
			if _, ok := checkNodes[id]; ok {
				richNodeList[id].Reset()
				richNodeList[id].Online = false
			}
			log.Println("[CLOSE] 关闭链接")
		}
		es.sockets.Delete(id)
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
						es.sockets.Store(strId, c)
						authorizes[c.RemoteAddr().String()] = strId
						richNodeList[strId].Online = true
						out = msg.Write(msg.SuccessAuthorizeMessage)
						closer = false
					} else {
						out = msg.Write(msg.NotEnableFailMessage)
						closer = false
					}
				} else {
					out = msg.Write(msg.NotExistFailMessage)
					closer = false
				}
			}
		case msg.ReceiveMessage:
			////取出id
			if id, err := utils.TrimLine(buf); err == nil {
				strId := string(id[:])
				if _, ok := conf.Get(strId); ok {
					//取出数据
					if sys, err := utils.TrimLine(buf); err == nil {
						richNodeList[strId].SystemInfo.Set(sys)
						closer = false
					}
				}
			}
		case msg.HeartbeatMessage:
			if id, err := utils.TrimLine(buf); err == nil {
				strId := string(id[:])
				if _, ok := conf.Get(strId); ok {
					out = msg.Write(msg.HeartbeatMessage, "pong")
					closer = false
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
	locker.Lock()
	data := conf.GetData()
	for i, _ := range data {
		m := data[i].(map[string]interface{})
		id := m["id"].(string)

		if _, ok := richNodeList[id]; !ok {
			node := &msg.Node{
				Id:       id,
				Name:     m["name"].(string),
				Location: m["location"].(string),
				Enable:   m["enable"].(bool),
				Region:   m["region"].(string),
			}

			richNodeList[id] = msg.NewRichNode(node, &msg.SystemInfo{}, false)
		}

		checkNodes[id] = true
	}

	//清除配置中不存在的节点
	for id, _ := range checkNodes {
		if _, ok := conf.Get(id); !ok {
			log.Println("[Not Exist Node]", ok)
			delete(richNodeList, id)
			delete(checkNodes, id)

			if c, ok := echo.sockets.Load(id); ok {
				_ = c.(gnet.Conn).Close()
			}
		}
	}

	r := msg.NewResponse(m, richNodeList)
	rsp, _ = r.Json()
	locker.Unlock()
	return
}

func Run(p *cmd.Cmd) {
	conf = config.NewConfig(p.Filename, make([]interface{}, 0))
	echo = &echoServer{}

	_ = response("init")

	go func() {
		log.Fatal(gnet.Serve(echo,
			fmt.Sprintf("tcp://%s:%d", p.Host, p.Port),
			gnet.WithMulticore(p.Multicore)))
	}()

	go func() {
		for {
			select {
			case <-conf.C:
				_ = response("reload")
			}
		}
	}()

	httpRun(p)
}
