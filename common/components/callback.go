package components

type (
	ConnectionCallback    func(*TcpConnection)
	CloseCallback         func(*TcpConnection)
	MessageCallback       func(*TcpConnection, *Buffer)
	TimerCallback         func(IConnection)
	WriteCompleteCallback func(*TcpConnection)
)
