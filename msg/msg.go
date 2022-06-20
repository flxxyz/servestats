package msg

import (
	"encoding/xml"
	jsoniter "github.com/json-iterator/go"
	"reflect"

	"github.com/flxxyz/servestats/config"
)

const (
	ServerReceiveType    = "RECEIVE"
	ServerReceiveMessage = "ok"
	ClientConnectType    = "CONNECT"
	ClientReportType     = "REPORT"
	ClientReplyMessage   = "Reply"
	ClientErrorMessage   = "Error"
	ClientFailMessage    = "Fail"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

// RpcServerReply 服务端回复消息
type RpcServerReply []byte

// RpcClientMessage 客户端响应消息
type RpcClientMessage struct {
	Id string
	*SystemInfo
}

// ResponseNode 响应消息中单个节点
type ResponseNode struct {
	Id       string `json:"-" xml:"-"`
	Name     string `json:"name" xml:"name"`
	Location string `json:"location" xml:"location"`
	Region   string `json:"region" xml:"region"`
	*SystemInfo
}

// Response 响应消息
type Response struct {
	Message string                   `json:"message" xml:"message"`
	Servers []*ResponseNode          `json:"server" xml:"server"`
	Nodes   map[string]*ResponseNode `json:"-" xml:"-"`
}

// Json 输出json格式数据
func (rsp *Response) Json() (data []byte, err error) {
	data, err = json.Marshal(rsp)

	return
}

// JsonFormat 输出格式化的json格式数据
func (rsp *Response) JsonFormat(prefix, indent string) (data []byte, err error) {
	data, err = json.MarshalIndent(rsp, prefix, indent)

	return
}

// XML 输出xml格式数据
func (rsp *Response) XML() (data []byte, err error) {
	data, err = xml.Marshal(rsp)

	return
}

// XMLFormat 输出格式化的xml格式数据
func (rsp *Response) XMLFormat(prefix, indent string) (data []byte, err error) {
	data, err = xml.MarshalIndent(rsp, prefix, indent)

	return
}

// Update 更新消息
func (rsp *Response) Update() {
	for {
		select {
		case <-config.Conf.C:
			rsp.Message = "update"
			rsp.updateNodes()
			rsp.updateServers()
		}
	}
}

func (rsp *Response) updateNodes() {
	for id, node := range config.Conf.All() {
		if _, ok := rsp.Nodes[id]; ok {
			//更新节点
			if node.Enable {
				rsp.Nodes[id].Id = id
				rsp.Nodes[id].Name = node.Name
				rsp.Nodes[id].Location = node.Location
				rsp.Nodes[id].Region = node.Region
			}
		} else {
			//添加节点
			if node.Enable {
				rsp.Nodes[id] = &ResponseNode{
					Id:       id,
					Name:     node.Name,
					Location: node.Location,
					Region:   node.Region,
					SystemInfo: &SystemInfo{
						OS: NewUnknownOS(),
					},
				}
			}
		}
	}
}

// FindServer 查找服务器
func (rsp *Response) FindServer(id string) (int, bool, error) {
	sVal := reflect.ValueOf(rsp.Servers)
	kind := sVal.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		for i := 0; i < sVal.Len(); i++ {
			n := sVal.Index(i).Interface().(*ResponseNode)
			if n.Id == id {
				return i, true, nil
			}
		}

		return 0, false, nil
	}

	return -1, false, nil
}

func (rsp *Response) updateServers() {
	for id, _ := range rsp.Nodes {
		if _, ok, _ := rsp.FindServer(id); !ok {
			rsp.Servers = append(rsp.Servers, rsp.Nodes[id])
		}
	}
}
