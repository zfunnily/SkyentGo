package components

import "sync"

type ConnManage struct {
	mu    sync.RWMutex
	conns map[uint32]IConnection

	r2cRW sync.RWMutex
	r2c   map[string]uint32 // role to connID

	u2cRW sync.RWMutex
	u2c   map[string]uint32 // uid to connID
}

func NewConnManage() *ConnManage {
	return &ConnManage{
		mu:    sync.RWMutex{},
		conns: make(map[uint32]IConnection),

		r2cRW: sync.RWMutex{},
		r2c:   make(map[string]uint32),

		u2cRW: sync.RWMutex{},
		u2c:   make(map[string]uint32),
	}
}

func (c *ConnManage) AddConn(id uint32, connection IConnection) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.conns[id] = connection
}

func (c *ConnManage) GetConn(id uint32) IConnection {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.conns[id]
}

func (c *ConnManage) DelConn(id uint32) IConnection {
	c.mu.Lock()
	defer c.mu.Unlock()
	conn := c.conns[id]
	delete(c.conns, id)
	return conn
}

func (c *ConnManage) Range(f func(key interface{}, value interface{}) bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, v := range c.conns {
		if ok := f(k, v); !ok {
			return
		}
	}
}

func (c *ConnManage) StopAllConns() {
	c.Range(func(key interface{}, value interface{}) bool {
		conn := value.(IConnection)
		conn.Stop()
		return true
	})
}

func (c *ConnManage) AddRID(rid string, id uint32) {
	c.r2cRW.Lock()
	defer c.r2cRW.Unlock()
	c.r2c[rid] = id
}

func (c *ConnManage) DelRID(rid string) {
	c.r2cRW.Lock()
	defer c.r2cRW.Unlock()
	delete(c.r2c, rid)
}

func (c *ConnManage) GetConnByRID(rid string) IConnection {
	c.r2cRW.RLock()
	defer c.r2cRW.RUnlock()
	cid := c.r2c[rid]
	return c.GetConn(cid)
}

func (c *ConnManage) AddUID(uid string, id uint32) {
	c.u2cRW.Lock()
	defer c.u2cRW.Unlock()
	c.u2c[uid] = id
}

func (c *ConnManage) DelUID(uid string) {
	c.r2cRW.Lock()
	defer c.r2cRW.Unlock()
	delete(c.r2c, uid)
}

func (c *ConnManage) GetConnByUID(uid string) IConnection {
	c.u2cRW.RLock()
	defer c.u2cRW.RUnlock()
	cid := c.u2c[uid]
	return c.GetConn(cid)
}
