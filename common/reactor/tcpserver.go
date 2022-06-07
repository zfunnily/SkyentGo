package reactor

import (
	"fmt"
)

type ConnectionMap map[string]*TcpConnection

type TcpServer struct {
	Loop                  *EventLoop
	acceptor              *Acceptor
	connections           ConnectionMap
	connectionCallback    ConnectionCallback
	messageCallback       MessageCallback
	writeCompleteCallback WriteCompleteCallback
	nextConnId            int
	name                  string
	LoopPool              *EventLoopPool
}

func NewTcpServer(loop *EventLoop, port int, name string) (*TcpServer, error) {
	acceptor, err := NewAcceptor(loop, port)
	if err != nil {
		return nil, err
	}
	s := &TcpServer{
		Loop:        loop,
		acceptor:    acceptor,
		connections: make(ConnectionMap),
		nextConnId:  0,
		name:        name,
		LoopPool:    NewEventLoopPool(loop, name),
	}
	s.connectionCallback = s.defaultConnectionCallback
	s.messageCallback = s.defaultMessageCallback
	s.LoopPool.SetNum(10)
	acceptor.SetNewConnectionCallback(s.NewConnection)
	return s, nil
}

func (s *TcpServer) SetConnectionCallback(cbk ConnectionCallback) {
	s.connectionCallback = cbk
}

func (s *TcpServer) SetMessageCallback(cbk MessageCallback) {
	s.messageCallback = cbk
}

func (s *TcpServer) SetWriteCompleteCallback(cbk WriteCompleteCallback) {
	s.writeCompleteCallback = cbk
}

func (s *TcpServer) NewConnection(sockfd int) {
	s.nextConnId++
	name := fmt.Sprintf("%s:%d", s.name, s.nextConnId)

	loop := s.LoopPool.GetNextLoop()
	conn := NewTcpConnection(loop, sockfd, name)
	s.connections[name] = conn

	conn.SetConnectionCallback(s.connectionCallback)
	conn.SetMessageCallback(s.messageCallback)
	conn.SetWriteCompleteCallback(s.writeCompleteCallback)
	conn.SetCloseCallback(s.removeConnection)

	loop.RunInLoop(func() {
		conn.connectEstablished()
	})
}

func (s *TcpServer) removeConnection(conn *TcpConnection) {
	s.removeConnectionInLoop(conn)
}

func (s *TcpServer) removeConnectionInLoop(conn *TcpConnection) {
	delete(s.connections, conn.name)
	conn.connectDestoryed()
}

func (s *TcpServer) defaultConnectionCallback(conn *TcpConnection) {
}

func (s *TcpServer) defaultMessageCallback(conn *TcpConnection, buffer *Buffer) {
}

func (s *TcpServer) Start() error {
	s.LoopPool.Start()
	return s.acceptor.listen()
}
