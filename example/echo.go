package main

import (
	"fmt"
	"pro2d/common/components"
)

type EchoServer struct {
	Server *components.TcpServer
}

func NewEchoServer(loop *components.EventLoop, port int) *EchoServer {
	s, err := components.NewTcpServer(loop, port, "echo")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	e := &EchoServer{Server: s}
	e.Server.SetConnectionCallback(e.OnConnect)
	e.Server.SetMessageCallback(e.OnMessage)
	return e
}

func (s *EchoServer) OnConnect(conn *components.TcpConnection) {
	fmt.Println("a new conn")
}

func (s *EchoServer) OnMessage(conn *components.TcpConnection, buffer *components.Buffer) {
	fmt.Printf("recv msg: %s\n", buffer.Peek()[:buffer.ReadableBytes()])
	conn.Send(buffer)
	buffer.RetrieveAll()
}

func (s *EchoServer) Start() error {
	return s.Server.Start()
}
func main() {
	loop := components.NewEventLoop()
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
