package msg

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
	Id       string `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Enable   bool   `json:"-"`
	Region   string `json:"region"`
}

type RichNode struct {
	*Node
	*SystemInfo
	Online bool `json:"online"`
}

func NewRichNode(node *Node, sys *SystemInfo, online bool) *RichNode {
	return &RichNode{node, sys, online}
}

type Response struct {
	Message string      `json:"message"`
	Servers []*RichNode `json:"servers"`
}

func (rsp *Response) Json() (data []byte, err error) {
	data, err = json.Marshal(rsp)
	return
}

func (rsp *Response) JsonFormat(prefix, indent string) (data []byte, err error) {
	data, err = json.MarshalIndent(rsp, prefix, indent)

	return
}

func NewResponse(message string, richNodes map[string]*RichNode) (r *Response) {
	r = &Response{
		Message: message,
		Servers: make([]*RichNode, 0),
	}

	for key, _ := range richNodes {
		r.Servers = append(r.Servers, richNodes[key])
	}

	return
}
