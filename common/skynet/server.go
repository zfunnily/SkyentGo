package skynet

type CtxCallback func(ctx *Context, cbud interface{}, data interface{}, typ int)

type Context struct {
	handle uint32
	queue  *MessageQueue
	cbud   interface{}
	cb     CtxCallback
}

func NewContext() *Context {
	ctx := &Context{}
	handle := HSInst().Register(ctx)
	ctx.handle = handle

	q := NewMessageQueue(handle)
	GQInst().Push(q)
	ctx.queue = q
	return ctx
}

func (c *Context) Handle() uint32 {
	return c.handle
}

func (c *Context) Callback(cbud interface{}, cb CtxCallback) {
	c.cbud = cbud
	c.cb = cb
}

func (c *Context) ContextSend(source uint32, typ int, session int, data []byte, sz uint64) {
	msg := &Message{
		Source:  source,
		Session: session,
		Data:    data,
		Typ:     typ,
	}
	c.queue.Push(msg)
}

// Send TODO
func (c *Context) Send(source uint32, destination uint32, typ int, session int, data []byte) int {
	if source == 0 {
		source = c.handle
	}

	msg := &Message{
		Source:  source,
		Session: session,
		Data:    data,
		Typ:     typ,
	}
	if ContextPush(destination, msg) != 0 {
		return -1
	}
	return session
}

func (c *Context) DispatchMessage(message *Message) {
	if c.cb != nil {
		c.cb(c, c.cbud, message.Data, message.Typ)
	}
}

func (c *Context) DispatchMessageAll() {
	for msg := c.queue.Pop(); msg != nil; {
		c.DispatchMessage(msg)
	}
}

func ContextPush(handle uint32, message *Message) int {
	ctx := HSInst().Grab(handle)
	if ctx == nil {
		return -1
	}
	ctx.queue.Push(message)
	return 0
}

func ContextMessageDispatch(q *MessageQueue) *MessageQueue {
	if q == nil {
		q = GQInst().Pop()
		if q == nil {
			return nil
		}
	}
	handle := q.Handle()
	ctx := HSInst().Grab(handle)
	if ctx == nil {
		return GQInst().Pop()
	}
	n := 1
	for i := 0; i < n; i++ {
		msg := ctx.queue.Pop()
		if msg != nil {
			ctx.DispatchMessage(msg)
		}
	}

	nq := GQInst().Pop()
	if nq != nil {
		GQInst().Push(q)
		q = nq
	}
	return q
}
