package components

const (
	DATA   = 0
	CLOSE  = 1
	ACCEPT = 2
	ERROR  = 4
	WARING = 5
	OPEN   = 6
)

type ResultMessage struct {
	id     uint32
	opaque uint32 // handle
	data   []byte
}

type SocketMessage struct {
	id     uint32
	typ    int
	opaque uint32 // handle
	data   []byte
}

type ServerOption func(*SocketServer)

func WithPlugin(iPlugin IPlugin) ServerOption {
	return func(server *SocketServer) {
		server.plugins = iPlugin
	}
}

func WithSplitter(splitter ISplitter) ServerOption {
	return func(server *SocketServer) {
		server.splitter = splitter
	}
}

func WithConnCbk(cb ConnectionCallback) ServerOption {
	return func(server *SocketServer) {
		server.connectionCallback = cb
	}
}

func WithMsgCbk(cb MessageCallback) ServerOption {
	return func(server *SocketServer) {
		server.messageCallback = cb
	}
}

func WithCloseCbk(cb CloseCallback) ServerOption {
	return func(server *SocketServer) {
		server.closeCallback = cb
	}
}

func WithTimerCbk(cb TimerCallback) ServerOption {
	return func(server *SocketServer) {
		server.timerCallback = cb
	}
}

func NewServer(port int, options ...ServerOption) IServer {
	s := &SocketServer{
		port:       port,
		connManage: NewConnManage(),
		Ctx:        NewContext(),
	}
	for _, option := range options {
		option(s)
	}

	return s
}

func (s *SocketServer) ForwardMessage(typ int, result *ResultMessage) {
	sm := &SocketMessage{
		id:     result.id,
		typ:    typ,
		opaque: result.opaque,
		data:   result.data,
	}

	msg := &Message{
		Source:  0,
		Session: 0,
		Data:    sm,
		Typ:     PTYPE_SOCKET,
	}

	ContextPush(result.opaque, msg)
}

func (s *SocketServer) Poll() {
	for {
		select {
		case accept := <-s.accept:
			result := &ResultMessage{
				id:     accept.GetID(),
				opaque: accept.GetCtx().Handle(),
				data:   nil,
			}
			s.ForwardMessage(ACCEPT, result)

		case clo := <-s.close:
			result := &ResultMessage{
				id:     clo.GetID(),
				opaque: clo.GetCtx().Handle(),
				data:   nil,
			}
			s.ForwardMessage(CLOSE, result)
		case data := <-s.data:
			conn := s.connManage.GetConn(data.GetSID())
			result := &ResultMessage{
				id:     conn.GetID(),
				opaque: conn.GetCtx().Handle(),
				data:   data.GetData(),
			}
			s.ForwardMessage(DATA, result)
		}
	}
}

func (s *SocketServer) AddEvent() {
}

func (s *SocketServer) DelEvent() {
}

func (s *SocketServer) newConnection(conn IConnection) {
	conn.SetConnectionCallback(s.connectionCallback)
	conn.SetCloseCallback(s.removeConnection)
	conn.SetMessageCallback(s.messageCallback)
	conn.SetTimerCallback(s.timerCallback)

	conn.Start()
}

func (s *SocketServer) removeConnection(conn IConnection) {
	s.closeCallback(conn)
}

func (s *SocketServer) Start() error {
	port := fmt.Sprintf(":%d", s.port)
	l, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	//监听端口
	logger.Debug("listen on %s\n", port)
	id := 0
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}

		id++
		client := NewConn(id, conn, s.splitter)
		s.accept <- client

		//s.newConnection(client)
	}
}
