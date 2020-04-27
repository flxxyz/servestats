package msg

const (
    AuthorizeMessage = iota + 30
    SuccessAuthorizeMessage
    FailAuthorizeMessage
    ReceiveMessage
    CloseMessage
    PingMessage
    PongMessage
)

type Node struct {
    Id       string `json:"id"`
    Name     string `json:"name"`
    Location string `json:"location"`
    Enable   bool   `json:"enable"`
    Region   string `json:"region"`
}

type RichNode struct {
    *Node

}

type Response struct {
    Message string `json:"message"`
    Data    []Node `json:"data"`
}

func NewResponse(message string, nodes []Node) *Response {
    return &Response{
        Message: message,
        Data:    nodes,
    }
}

type TestMsg struct {
    A int         `json:"a"`
    B string      `json:"b"`
    C interface{} `json:"c"`
}

func NewTestMsg() *TestMsg {
    return &TestMsg{
        A: 666,
        B: "测试",
        C: struct {
            Num0 interface{} `json:"0"`
            One  interface{} `json:"one"`
        }{
            Num0: "one",
            One:  "two",
        },
    }
}
