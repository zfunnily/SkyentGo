package components

import (
	"fmt"
	"syscall"
)

// CreatePoll EPoll
func (p *SocketPoll) CreatePoll() {
	fd, err := syscall.EpollCreate1(syscall.EPOLL_CLOEXEC)
	if err != nil {
		panic(err)
	}
	p.epFd = fd
}

// AddReadWrite ...
func (p *SocketPoll) AddReadWrite(fd int) {
	if err := syscall.EpollCtl(p.epFd, syscall.EPOLL_CTL_ADD, fd,
		&syscall.EpollEvent{Fd: int32(fd),
			Events: syscall.EPOLLIN | syscall.EPOLLOUT,
		},
	); err != nil {
		panic(err)
	}
}

// AddRead ...
func (p *SocketPoll) AddRead(fd int) {
	if err := syscall.EpollCtl(p.epFd, syscall.EPOLL_CTL_ADD, fd,
		&syscall.EpollEvent{Fd: int32(fd),
			Events: syscall.EPOLLIN,
		},
	); err != nil {
		panic(err)
	}
}

// ModRead ...
func (p *SocketPoll) ModRead(fd int) {
	if err := syscall.EpollCtl(p.epFd, syscall.EPOLL_CTL_MOD, fd,
		&syscall.EpollEvent{Fd: int32(fd),
			Events: syscall.EPOLLIN,
		},
	); err != nil {
		panic(err)
	}
}

// ModReadWrite ...
func (p *SocketPoll) ModReadWrite(fd int) {
	if err := syscall.EpollCtl(p.epFd, syscall.EPOLL_CTL_MOD, fd,
		&syscall.EpollEvent{Fd: int32(fd),
			Events: syscall.EPOLLIN | syscall.EPOLLOUT,
		},
	); err != nil {
		panic(err)
	}
}

// ModDetach ...
func (p *SocketPoll) ModDetach(fd int) {
	if err := syscall.EpollCtl(p.epFd, syscall.EPOLL_CTL_DEL, fd,
		&syscall.EpollEvent{Fd: int32(fd),
			Events: syscall.EPOLLIN | syscall.EPOLLOUT,
		},
	); err != nil {
		panic(err)
	}
}

func (p *SocketPoll) Poll(cb func(event *Event)) error {
	p.Events = make([]syscall.EpollEvent, 64)
	for {
		n, err := syscall.EpollWait(p.epFd, p.Events, 100)
		if err != nil && err != syscall.EINTR {
			return err
		}
		for i := 0; i < n; i++ {
			fd := p.Events[i].Fd
			eventM := p.EventMap[fd]
			if eventM == nil {
				eventM = NewEvent(fd)
				eventM.SetEvents(p.Events[i].Events)
			}
			cb(eventM)
		}
	}
}

func (p *SocketPoll) Update(operation int, event *Event) {
	fd := event.Fd()
	if err := syscall.EpollCtl(p.epFd, operation, int(fd),
		&syscall.EpollEvent{Fd: fd,
			Events: event.Events(),
		},
	); err != nil {
		fmt.Println(err)
	}
}

func (p *SocketPoll) UpdateEvent(event *Event) {
	fd := event.Fd()
	_, ok := p.EventMap[fd]
	if !ok {
		p.EventMap[fd] = event
		p.Update(syscall.EPOLL_CTL_ADD, event)
	} else {
		if event.IsNoEvent() {
			p.Update(syscall.EPOLL_CTL_DEL, event)
		} else {
			p.Update(syscall.EPOLL_CTL_MOD, event)
		}
	}
}

func (p *SocketPoll) RemoveEvent(event *Event) {
	fd := event.Fd()
	p.Update(syscall.EPOLL_CTL_DEL, event)
	delete(p.EventMap, fd)
}

//func (p *SocketPoll) Start() error {
//	port := fmt.Sprintf(":%d", s.port)
//
//	l, err := net.Listen("tcp", port)
//	if err != nil {
//		return err
//	}
//
//	s.Listener = l
//
//	var f *os.File
//	switch netln := l.(type) {
//	case nil:
//	case *net.TCPListener:
//		f, err = netln.File()
//	}
//	if err != nil {
//		return err
//	}
//
//	p.fd = int(f.Fd())
//	return syscall.SetNonblock(s.fd, true)
//}
