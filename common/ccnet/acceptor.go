package ccnet

import (
	"fmt"
	"net"
	"os"
	"syscall"
)

type NewConnectionCallback func(sofkfd int)

type Acceptor struct {
	Loop                  *EventLoop
	Event                 *Event
	newConnectionCallback NewConnectionCallback
	listening             bool
	Listener              net.Listener
	fd                    int
}

func NewAcceptor(loop *EventLoop, port int) (*Acceptor, error) {
	p := fmt.Sprintf(":%d", port)
	l, err := net.Listen("tcp", p)
	if err != nil {
		return nil, err
	}

	var f *os.File
	switch netln := l.(type) {
	case nil:
	case *net.TCPListener:
		f, err = netln.File()
		fmt.Println("tcp listener")
	}
	if err != nil {
		return nil, err
	}

	acceptor := &Acceptor{
		Listener:  l,
		listening: false,
		fd:        int(f.Fd()),
		Loop:      loop,
	}
	err = syscall.SetNonblock(acceptor.fd, true)
	if err != nil {
		return nil, err
	}

	acceptor.Event = NewEvent(acceptor.fd, loop)
	acceptor.Event.SetReadCbk(acceptor.handleAccept)

	fmt.Println("create acceptor successful, ", p)
	return acceptor, nil
}

func (a *Acceptor) Fd() int {
	return a.fd
}

func (a *Acceptor) SetNewConnectionCallback(cbk NewConnectionCallback) {
	a.newConnectionCallback = cbk
}

func (a *Acceptor) handleAccept() {
	nfd, _, err := syscall.Accept(a.fd)
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}
	if nfd > 0 {
		if err = syscall.SetNonblock(nfd, true); err != nil {
			fmt.Errorf(err.Error())
			return
		}
		if a.newConnectionCallback != nil {
			a.newConnectionCallback(nfd)
		} else {
			syscall.Close(a.fd)
			a.Listener.Close()
		}
	}
}

func (a *Acceptor) Close() {
	a.Event.DisableAll()
	a.Event.Remove()
}

func (a *Acceptor) listen() error {
	a.Event.EnableConnecting()
	fmt.Println("start listen fd,", a.fd, "...")
	return nil
}
