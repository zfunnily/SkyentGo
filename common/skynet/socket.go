package skynet

import (
	"net"
	"pro2d/common/components"
	"pro2d/common/reactor"
)

type SocketServer struct {
	efd      int
	fd       int
	Listener net.Listener

	PluginPath string
	plugins    components.IPlugin
	splitter   components.ISplitter

	connectionCallback reactor.ConnectionCallback
	messageCallback    reactor.MessageCallback
	closeCallback      reactor.CloseCallback
	timerCallback      reactor.TimerCallback

	port       int
	connManage components.IConnManage
	Actions    map[interface{}]interface{}

	accept chan components.IConnection
	data   chan components.IMessage
	close  chan components.IConnection

	Ctx *Context
}

var SS *SocketServer

func SSInst() *SocketServer {
	if SS == nil {
		return nil
	}
	return SS
}

func (s *SocketServer) Stop() {
	StopTimer()
	s.connManage.StopAllConns()
}

func (s *SocketServer) GetSplitter() components.ISplitter {
	return s.splitter
}

func (s *SocketServer) GetPlugin() components.IPlugin {
	return s.plugins
}

func (s *SocketServer) GetAction(cmd uint32) interface{} {
	if s.plugins != nil {
		f := s.plugins.GetAction(cmd)
		if f != nil {
			return f
		}
	}

	return s.Actions[cmd]
}

func (s *SocketServer) SetActions(mi map[interface{}]interface{}) {
	s.Actions = mi
}

func (s *SocketServer) GetConnManage() components.IConnManage {
	return s.connManage
}

func (s *SocketServer) SetConnectionCallback(cb reactor.ConnectionCallback) {
	s.connectionCallback = cb
}

func (s *SocketServer) SetMessageCallback(cb reactor.MessageCallback) {
	s.messageCallback = cb
}

func (s *SocketServer) SetCloseCallback(cb reactor.CloseCallback) {
	s.closeCallback = cb
}

func (s *SocketServer) SetTimerCallback(cb reactor.TimerCallback) {
	s.timerCallback = cb
}
