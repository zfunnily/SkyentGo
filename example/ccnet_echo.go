package main

import (
	"fmt"
	"pro2d/common/reactor"
)

type EchoServer struct {
	Server *reactor.TcpServer
}

func NewEchoServer(loop *reactor.EventLoop, port int) *EchoServer {
	s, err := reactor.NewTcpServer(loop, port, "echo")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	e := &EchoServer{Server: s}
	e.Server.SetConnectionCallback(e.OnConnect)
	e.Server.SetMessageCallback(e.OnMessage)
	return e
}

func (s *EchoServer) OnConnect(conn *reactor.TcpConnection) {
	fmt.Println("a new conn")
}

func (s *EchoServer) OnMessage(conn *reactor.TcpConnection, buffer *reactor.Buffer) {
	fmt.Printf("recv msg: %s\n", buffer.Peek()[:buffer.ReadableBytes()])
	conn.Send(buffer)
	buffer.RetrieveAll()
}

func (s *EchoServer) Start() error {
	return s.Server.Start()
}
func main() {
	loop := reactor.NewEventLoop()
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
