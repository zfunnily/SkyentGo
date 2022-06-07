package reactor

type (
	ConnectionCallback    func(*TcpConnection)
	CloseCallback         func(*TcpConnection)
	MessageCallback       func(*TcpConnection, *Buffer)
	TimerCallback         func(*TcpConnection)
	WriteCompleteCallback func(*TcpConnection)
)
