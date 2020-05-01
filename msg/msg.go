package msg

import (
	"encoding/xml"
	"sort"
)

type TypeMessage int

const (
	AuthorizeMessage byte = iota + 30
	SuccessAuthorizeMessage
	NotExistFailMessage
	NotEnableFailMessage
	ReceiveMessage
	CloseMessage
	HeartbeatMessage
)

type Node struct {
	Id       string `json:"-" xml:"-"`
	Name     string `json:"name" xml:"name"`
	Location string `json:"location" xml:"location"`
	Enable   bool   `json:"-" xml:"-"`
	Region   string `json:"region" xml:"region"`
}

type RichNode struct {
	*Node
	*SystemInfo
	Online bool `json:"online" xml:"online"`
}

func NewRichNode(node *Node, sys *SystemInfo, online bool) *RichNode {
	return &RichNode{node, sys, online}
}

type Response struct {
	Message    string      `json:"message" xml:"message"`
	Servers    []*RichNode `json:"server" xml:"server"`
	UpdateChan chan string `json:"-" xml:"-"`
	Data       []byte      `json:"-" xml:"-"`
}

func (rsp *Response) Json() (data []byte, err error) {
	data, err = json.Marshal(rsp)

	return
}

func (rsp *Response) JsonFormat(prefix, indent string) (data []byte, err error) {
	data, err = json.MarshalIndent(rsp, prefix, indent)

	return
}

func (rsp *Response) XML() (data []byte, err error) {
	data, err = xml.Marshal(rsp)

	return
}

func (rsp *Response) XMLFormat(prefix, indent string) (data []byte, err error) {
	data, err = xml.MarshalIndent(rsp, prefix, indent)

	return
}

func (rsp *Response) Update(richNodeList map[string]*RichNode) {
	keys := make([]string, 0)
	for key, _ := range richNodeList {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	rsp.Servers = rsp.Servers[0:0]
	for i, _ := range keys {
		id := keys[i]
		rsp.Servers = append(rsp.Servers, richNodeList[id])
	}

	return
}

func NewResponse(message string, richNodeList map[string]*RichNode) (r *Response) {
	r = &Response{
		Message:    message,
		Servers:    make([]*RichNode, 0),
		UpdateChan: make(chan string, 1024),
	}

	r.Update(richNodeList)

	return
}
