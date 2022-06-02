package components

import (
	"net"
)

type SocketServer struct {
	efd      int
	fd       int
	Listener net.Listener

	PluginPath string
	plugins    IPlugin
	splitter   ISplitter

	connectionCallback ConnectionCallback
	messageCallback    MessageCallback
	closeCallback      CloseCallback
	timerCallback      TimerCallback

	port       int
	connManage IConnManage
	Actions    map[interface{}]interface{}

	accept chan IConnection
	data   chan IMessage
	close  chan IConnection

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

func (s *SocketServer) GetSplitter() ISplitter {
	return s.splitter
}

func (s *SocketServer) GetPlugin() IPlugin {
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

func (s *SocketServer) GetConnManage() IConnManage {
	return s.connManage
}

func (s *SocketServer) SetConnectionCallback(cb ConnectionCallback) {
	s.connectionCallback = cb
}

func (s *SocketServer) SetMessageCallback(cb MessageCallback) {
	s.messageCallback = cb
}

func (s *SocketServer) SetCloseCallback(cb CloseCallback) {
	s.closeCallback = cb
}

func (s *SocketServer) SetTimerCallback(cb TimerCallback) {
	s.timerCallback = cb
}
