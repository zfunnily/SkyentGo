package main

import (
	"fmt"
	ccnet2 "pro2d/common/ccnet"
)

type EchoServer struct {
	Server *ccnet2.TcpServer
}

func NewEchoServer(loop *ccnet2.EventLoop, port int) *EchoServer {
	s, err := ccnet2.NewTcpServer(loop, port, "echo")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	e := &EchoServer{Server: s}
	e.Server.SetConnectionCallback(e.OnConnect)
	e.Server.SetMessageCallback(e.OnMessage)
	return e
}

func (s *EchoServer) OnConnect(conn *ccnet2.TcpConnection) {
	fmt.Println("a new conn")
}

func (s *EchoServer) OnMessage(conn *ccnet2.TcpConnection, buffer *ccnet2.Buffer) {
	fmt.Printf("recv msg: %s\n", buffer.Peek()[:buffer.ReadableBytes()])
	conn.Send(buffer)
	buffer.RetrieveAll()
}

func (s *EchoServer) Start() error {
	return s.Server.Start()
}
func main() {
	loop := ccnet2.NewEventLoop()
	e := NewEchoServer(loop, 80)
	if e == nil {
		return
	}
	if err := e.Start(); err != nil {
		fmt.Errorf(err.Error())
		return
	}
	loop.Loop()
}
