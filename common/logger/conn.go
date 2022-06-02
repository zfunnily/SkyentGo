package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

type connLogger struct {
	sync.Mutex
	innerWriter    io.WriteCloser
	Net            string `json:"net"`
	Addr           string `json:"addr"`
	Level          string `json:"level"`
	LogLevel       int
	illNetFlag     bool //网络异常标记
}

func (c *connLogger) Init(jsonConfig string, appName string) error {
	if len(jsonConfig) == 0 {
		return nil
	}
	//fmt.Printf("consoleWriter Init:%s\n", jsonConfig)
	err := json.Unmarshal([]byte(jsonConfig), c)
	if err != nil {
		return err
	}
	if l, ok := LevelMap[c.Level]; ok {
		c.LogLevel = l
	}
	if c.innerWriter != nil {
		c.innerWriter.Close()
		c.innerWriter = nil
	}

	go func() {
		for {
			c.connect()
			time.Sleep(10*time.Millisecond)
		}
	}()

	return nil
}

func (c *connLogger) LogWrite(when time.Time, msgText interface{}, level int) (err error) {
	if level > c.LogLevel {
		return nil
	}

	msg, ok := msgText.(*loginfo)
	if !ok {
		return
	}

	if c.innerWriter != nil {
		err = c.println(when, msg)
		//网络异常，通知处理网络的go程自动重连
		if err != nil {
			c.innerWriter.Close()
			c.innerWriter = nil
		}
	}

	return
}

func (c *connLogger) Destroy() {
	if c.innerWriter != nil {
		c.innerWriter.Close()
	}
}

func (c *connLogger) connect() error {
	if c.innerWriter != nil {
		return nil
	}
	addrs := strings.Split(c.Addr, ";")
	for _, addr := range addrs {
		conn, err := net.DialTimeout(c.Net, addr, 1 * time.Second)
		if err != nil {
			fmt.Printf("net.Dial error:%v\n", err)
			//continue
			return err
		}

		if tcpConn, ok := conn.(*net.TCPConn); ok {
			tcpConn.SetKeepAlive(true)
		}
		c.innerWriter = conn
		return nil
	}
	return fmt.Errorf("hava no valid logs service addr:%v", c.Addr)
}

func (c *connLogger) println(when time.Time, msg *loginfo) error {
	c.Lock()
	defer c.Unlock()
	ss, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = c.innerWriter.Write(append(ss, '\n'))

	//返回err，解决日志系统网络异常后的自动重连
	return err
}

func init() {
	Register(AdapterConn, &connLogger{LogLevel: LevelTrace})
}
