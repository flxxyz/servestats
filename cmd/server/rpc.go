package server

import (
	"github.com/flxxyz/ServerStatus/msg"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type MessageService int

func (s *MessageService) Report(args *msg.RpcClientMessage, reply *msg.RpcServerReply) error {
	if index, ok, _ := response.FindServer(args.Id); ok {
		s := response.Servers[index]
		if s.Id == args.Id {
			s.SystemInfo = args.SystemInfo
		}
		log.Printf("[%s] ID: %s\n", msg.ServerReceiveType, args.Id)
	}
	*reply = msg.RpcServerReply(msg.ServerReceiveMessage)
	return nil
}

func RPCServer(addr string) {
	_ = rpc.Register(new(MessageService))
	rpc.HandleHTTP()

	go response.Update()

	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}

	go func() {
		log.Println("[RPCServer]", l.Addr())
		log.Fatal(http.Serve(l, nil))
	}()
	select {}
}
