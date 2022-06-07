package components

import (
	"fmt"
)

type ConnectionMap map[string]*TcpConnection

type TcpServer struct {
	Loop     *EventLoop
	acceptor *Acceptor

	connections ConnectionMap
	//cbk
	connectionCallback    ConnectionCallback
	messageCallback       MessageCallback
	writeCompleteCallback WriteCompleteCallback

	nextConnId int
	name       string
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
	}
	s.connectionCallback = s.defaultConnectionCallback
	s.messageCallback = s.defaultMessageCallback
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
	conn := NewTcpConnection(s.Loop, sockfd, name)
	s.connections[name] = conn

	conn.SetConnectionCallback(s.connectionCallback)
	conn.SetMessageCallback(s.messageCallback)
	conn.SetWriteCompleteCallback(s.writeCompleteCallback)
	conn.SetCloseCallback(s.removeConnection)

	s.Loop.RunInLoop(func() {
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
	return s.acceptor.listen()
}
