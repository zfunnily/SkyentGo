package components

import (
	"fmt"
	"syscall"
)

type ConnState int

const (
	kDisconnected  ConnState = 0
	kConnecting    ConnState = 1
	kConnected     ConnState = 2
	kDisconnecting ConnState = 3
)

type TcpConnection struct {
	IConnection
	Loop      *EventLoop
	Event     *Event
	connState ConnState
	splitter  ISplitter
	name      string

	Id           uint32
	fd           int
	InputBuffer  *Buffer
	OutPutBuffer *Buffer

	// cbk
	connectionCallback    ConnectionCallback
	messageCallback       MessageCallback
	writeCompleteCallback WriteCompleteCallback
	timerCallback         TimerCallback
	closeCallback         CloseCallback
}

func NewTcpConnection(loop *EventLoop, fd int, name string) *TcpConnection {
	conn := &TcpConnection{
		Loop:         loop,
		fd:           fd,
		Event:        NewEvent(fd, loop),
		name:         name,
		InputBuffer:  NewBuffer(),
		OutPutBuffer: NewBuffer(),
	}

	conn.Event.SetReadCbk(conn.handleRead)
	conn.Event.SetWriteCbk(conn.handleWrite)
	conn.Event.SetCloseCbk(conn.handleClose)

	return conn
}

func (c *TcpConnection) connectEstablished() {
	c.setConnState(kConnected)
	c.Event.EnableReading()
	c.connectionCallback(c)
}

func (c *TcpConnection) connectDestoryed() {
	if c.connState == kConnected {
		c.setConnState(kDisconnected)
		c.Event.DisableAll()
		c.closeCallback(c)
	}
}

func (c *TcpConnection) Send(buf *Buffer) {
	if c.connState == kConnected {
		c.Loop.RunInLoop(func() {
			c.SendInLoop(buf.Peek(), buf.ReadableBytes())
		})
	}
}

func (c *TcpConnection) SendInLoop(buf []byte, len int64) {
	if c.connState == kConnected {
		c.OutPutBuffer.Append(buf, len)
		if c.Event.IsWriting() <= 0 {
			c.Event.EnableWriting()
		}
	}
}

func (c *TcpConnection) setConnState(state ConnState) {
	c.connState = state
}

func (c *TcpConnection) SetConnectionCallback(cbk ConnectionCallback) {
	c.connectionCallback = cbk
}

func (c *TcpConnection) SetMessageCallback(cbk MessageCallback) {
	c.messageCallback = cbk
}

func (c *TcpConnection) SetWriteCompleteCallback(cbk WriteCompleteCallback) {
	c.writeCompleteCallback = cbk
}

func (c *TcpConnection) SetCloseCallback(cbk CloseCallback) {
	c.closeCallback = cbk
}

func (c *TcpConnection) shutdown() {
	if c.connState == kConnected {
		c.setConnState(kDisconnected)
		c.Loop.RunInLoop(func() {
			if c.Event.IsWriting() == 0 {
				syscall.Shutdown(c.fd, syscall.SHUT_WR)
			}
		})
	}
}

func (c *TcpConnection) handleRead() {
	n := c.InputBuffer.ReadFd(c.fd)
	if n > 0 {
		c.messageCallback(c, c.InputBuffer)
	} else if n == 0 {
		c.handleClose()
	} else {
		c.handleError()
	}
}

func (c *TcpConnection) handleWrite() {
	if c.Event.IsWriting() > 0 {
		n, err := syscall.Write(c.fd, c.OutPutBuffer.Peek()[:c.OutPutBuffer.ReadableBytes()])
		if err != nil {
			fmt.Errorf(err.Error())
			return
		}
		if n > 0 {
			c.OutPutBuffer.Retrieve(int64(n))
			if c.OutPutBuffer.ReadableBytes() == 0 {
				c.Event.DisableWriting()
			}
		}
	}
}

func (c *TcpConnection) handleClose() {
	c.setConnState(kDisconnected)
	c.Event.DisableAll()
	c.closeCallback(c)
}

func (c *TcpConnection) handleError() {
}
