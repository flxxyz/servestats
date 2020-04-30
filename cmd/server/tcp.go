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
	"strings"
	"sync"
)

var (
	echo          *echoServer
	conf          *config.Config
	authorizes    map[string]string
	richNodeList  map[string]*msg.RichNode
	checkNodeList map[string]bool
	locker        *sync.RWMutex
	r             *msg.Response
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
			if _, ok := checkNodeList[id]; ok {
				richNodeList[id].Reset()
				richNodeList[id].Online = false
			}
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
					if m["enable"].(bool) {
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
						r.UpdateChan <- "receive:" + strId
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

func init() {
	echo = &echoServer{}
	authorizes = make(map[string]string, 0)
	richNodeList = make(map[string]*msg.RichNode, 0)
	checkNodeList = make(map[string]bool, 0)
	locker = &sync.RWMutex{}
	r = msg.NewResponse("init", make(map[string]*msg.RichNode, 0))
}

func Run(p *cmd.Cmd) {
	conf = config.NewConfig(p.Filename, make([]interface{}, 0))

	go func() {
		for {
			select {
			case <-conf.C:
				r.UpdateChan <- "reload"
				log.Println("[Reload]", "config.json")
			case message := <-r.UpdateChan:
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

					checkNodeList[id] = true
				}

				//清除配置中不存在的节点
				for id, _ := range checkNodeList {
					if _, ok := conf.Get(id); !ok {
						delete(richNodeList, id)
						delete(checkNodeList, id)

						if c, ok := echo.sockets.Load(id); ok {
							_ = c.(gnet.Conn).Close()
						}
					}
				}

				r.Message = strings.Split(message, ":")[0]
				r.Update(richNodeList)
				locker.Unlock()
				log.Println("[Update]", "message:"+message)
			}
		}
	}()

	r.UpdateChan <- "init"

	go func() {
		log.Fatal(gnet.Serve(echo,
			fmt.Sprintf("tcp://%s:%d", p.Host, p.Port),
			gnet.WithMulticore(p.Multicore)))
	}()
	httpRun(p)
}
