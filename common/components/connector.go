package components

import (
	"fmt"
	"net"
	"pro2d/common/logger"
	//"pro2d/pb"
)

type ConnectorOption func(*Connector)

func WithCtorSplitter(splitter ISplitter) ConnectorOption {
	return func(connector *Connector) {
		connector.splitter = splitter
	}
}

func WithCtorCount(count int) ConnectorOption {
	return func(connector *Connector) {
		connector.sum = count
	}
}

type Connector struct {
	IConnector
	IConnection
	IServer
	Id       int
	splitter ISplitter
	ip       string
	port     int
	sum      int
}

func NewConnector(ip string, port int, options ...ConnectorOption) IConnector {
	c := &Connector{
		ip:   ip,
		port: port,
	}
	for _, option := range options {
		option(c)
	}
	return c
}

func (c *Connector) GetConn() IConnection {
	return c.IConnection
}

func (c *Connector) Connect() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.ip, c.port))
	if err != nil {
		return err
	}
	cli := NewConn(c.Id, conn, c.splitter)
	cli.SetMessageCallback(c.OnMessage)
	cli.SetCloseCallback(c.OnClose)
	cli.SetTimerCallback(c.OnTimer)
	c.IConnection = cli

	return nil
}

func (c *Connector) DisConnect() {
	c.IConnection.Stop()
}

func (c *Connector) Send(cmd uint32, b []byte) error {
	logger.Debug("connector send cmd: %d, msg: %s", cmd, b)
	return c.IConnection.Send(0, cmd, b)
}

//func (c *Connector) SendPB(cmd pb.ProtoCode, b proto.Message) error {
//	if b == nil {
//		return c.Send(uint32(cmd), nil)
//	}
//
//	l, err := proto.Marshal(b)
//	if err != nil {
//		return err
//	}
//	return c.Send(uint32(cmd), l)
//}

func (c *Connector) GetSplitter() ISplitter {
	return c.splitter
}

func (c *Connector) OnMessage(msg IMessage) {
	logger.Debug("recv msg errorCode: %d cmd: %d, conn: %d data: %s", msg.GetHeader().GetErrCode(), msg.GetHeader().GetMsgID(), msg.GetSID(), msg.GetData())
}

func (c *Connector) OnClose(conn IConnection) {
	logger.Debug("onclose id: %d", conn.GetID())
}

func (c *Connector) OnTimer(conn IConnection) {
	//logger.Debug("ontimer id: %d", conn.GetID())
}
