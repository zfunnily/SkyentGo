package components

import (
	"bufio"
	"fmt"
	"net"
	"pro2d/common"
	"pro2d/common/logger"
	"sync"
	"sync/atomic"
	"time"
)

type Connection struct {
	IConnection
	net.Conn
	splitter ISplitter
	Id       uint32
	fd       int

	scanner       *bufio.Scanner
	writer        *bufio.Writer
	WBuffer       chan []byte
	Quit          chan *Connection
	readFunc      chan func()
	timerFunc     chan func()
	customizeFunc chan func()

	messageCallback    MessageCallback
	connectionCallback ConnectionCallback
	closeCallback      CloseCallback
	timerCallback      TimerCallback

	Status uint32

	Ctx *Context
}

var connectionPool = &sync.Pool{
	New: func() interface{} { return new(Connection) },
}

func NewConnFd(fd int) IConnection {
	return &Connection{fd: fd}
}

func NewConn(id int, conn net.Conn, splitter ISplitter) IConnection {
	c := new(Connection)
	status := atomic.LoadUint32(&c.Status)
	if status != 0 {
		connectionPool.Put(c)
		c = new(Connection)
	}
	c.Ctx = NewContext()

	atomic.StoreUint32(&c.Id, uint32(id))
	c.Conn = conn
	c.splitter = splitter

	c.scanner = bufio.NewScanner(conn)
	c.writer = bufio.NewWriter(conn)

	c.reset()

	return c
}

func (c *Connection) reset() {
	c.WBuffer = make(chan []byte, common.MaxMsgChan)
	c.Quit = make(chan *Connection)

	if c.readFunc == nil {
		c.readFunc = make(chan func(), 10)
	}
	if c.timerFunc == nil {
		c.timerFunc = make(chan func(), 10)
	}
	if c.customizeFunc == nil {
		c.customizeFunc = make(chan func(), 10)
	}

	//c.connectionCallback	= c.defaultConnectionCallback
	//c.messageCallback 		= c.defaultMessageCallback
	//c.closeCallback 		= c.defaultCloseCallback
	//c.timerCallback 		= c.defaultTimerCallback
}

func (c *Connection) GetID() uint32 {
	return atomic.LoadUint32(&c.Id)
}

func (c *Connection) SetConnectionCallback(cb ConnectionCallback) {
	c.connectionCallback = cb
}

func (c *Connection) SetMessageCallback(cb MessageCallback) {
	c.messageCallback = cb
}

func (c *Connection) SetCloseCallback(cb CloseCallback) {
	c.closeCallback = cb
}

func (c *Connection) SetTimerCallback(cb TimerCallback) {
	c.timerCallback = cb
}

func (c *Connection) Start() {
	go c.write()
	go c.read()
	go c.listen()

	atomic.StoreUint32(&c.Status, 1)
	c.connectionCallback(c)
	c.handleTimeOut()
}

func (c *Connection) Stop() {
	if atomic.LoadUint32(&c.Status) == 0 {
		return
	}

	sendTimeout := time.NewTimer(5 * time.Millisecond)
	defer sendTimeout.Stop()
	// 发送超时
	select {
	case <-sendTimeout.C:
		return
	case c.Quit <- c:
		return
	}
}

func (c *Connection) Send(errCode int32, cmd uint32, data []byte) error {
	buf, err := c.splitter.Pack(cmd, data, errCode, 0)
	if err != nil {
		return err
	}

	sendTimeout := time.NewTimer(5 * time.Millisecond)
	defer sendTimeout.Stop()
	// 发送超时
	select {
	case <-sendTimeout.C:
		return fmt.Errorf("send buff msg timeout")
	case c.WBuffer <- buf:
		return nil
	}
}

func (c *Connection) SendSuccess(cmd uint32, data []byte) error {
	buf, err := c.splitter.Pack(cmd, data, 0, 0)
	if err != nil {
		return err
	}

	sendTimeout := time.NewTimer(5 * time.Millisecond)
	defer sendTimeout.Stop()
	// 发送超时
	select {
	case <-sendTimeout.C:
		return fmt.Errorf("send buff msg timeout")
	case c.WBuffer <- buf:
		return nil
	}
}

func (c *Connection) CustomChan() chan<- func() {
	return c.customizeFunc
}

func (c *Connection) defaultConnectionCallback(conn IConnection) {
}

func (c *Connection) defaultMessageCallback(msg IMessage) {
}

func (c *Connection) defaultCloseCallback(conn IConnection) {
}

func (c *Connection) defaultTimerCallback(conn IConnection) {
}

func (c *Connection) write() {
	defer func() {
		//logger.Debug("write close")
		c.Stop()
	}()

	for msg := range c.WBuffer {
		n, err := c.writer.Write(msg)
		if err != nil {
			logger.Error("write fail err: "+err.Error(), "n: ", n)
			return
		}
		if err := c.writer.Flush(); err != nil {
			logger.Error("write Flush fail err: " + err.Error())
			return
		}
		logger.Debug("write n: %d", n)
	}
}

func (c *Connection) read() {
	defer func() {
		c.Stop()
	}()

	c.scanner.Split(c.splitter.ParseMsg)

	for c.scanner.Scan() {
		req, err := c.splitter.UnPack(c.scanner.Bytes())
		if err != nil {
			return
		}

		req.SetSID(c.GetID())
		c.readFunc <- func() {
			c.messageCallback(req)
		}
	}

	if err := c.scanner.Err(); err != nil {
		logger.Error("scanner.err: %s", err.Error())
		return
	}
}

//此设计目的是为了让网络数据与定时器处理都在一条协程里处理。不想加锁。。。
func (c *Connection) listen() {
	defer func() {
		//logger.Debug("listen close")
		c.quitting()
	}()

	for {
		select {
		case timerFunc := <-c.timerFunc:
			timerFunc()
		case readFunc := <-c.readFunc:
			readFunc()
		case customizeFunc := <-c.customizeFunc:
			customizeFunc()
		case <-c.Quit:
			return
		}
	}
}

func (c *Connection) handleTimeOut() {
	if atomic.LoadUint32(&c.Status) == 0 {
		return
	}

	c.timerFunc <- func() {
		c.timerCallback(c)
	}
	//TimeOut(1*time.Second, c.handleTimeOut)
}

func (c *Connection) quitting() {
	if atomic.LoadUint32(&c.Status) == 0 {
		return
	}
	atomic.StoreUint32(&c.Status, 0)

	close(c.WBuffer)
	close(c.Quit)

	c.Conn.Close()
	c.closeCallback(c)

	//放回到对象池
	connectionPool.Put(c)
}
