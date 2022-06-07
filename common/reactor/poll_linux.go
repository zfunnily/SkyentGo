package reactor

import (
	"fmt"
	"syscall"
)

func NewPoll(Loop *EventLoop) *SocketPoll {
	efd, err := syscall.EpollCreate1(syscall.EPOLL_CLOEXEC)
	if err != nil {
		panic(err)
	}
	return &SocketPoll{epFd: efd, Loop: Loop, EventMap: make(map[int]*Event)}
}

func (p *SocketPoll) Poll(cb func(event *Event)) error {
	p.Events = make([]syscall.EpollEvent, 64)
	for {
		n, err := syscall.EpollWait(p.epFd, p.Events, -1)
		if err != nil && err != syscall.EINTR {
			return err
		}
		if n != 0 {
			fmt.Println("wait n:", n)
		}
		for i := 0; i < n; i++ {
			fd := int(p.Events[i].Fd)
			eventM := p.EventMap[fd]
			if eventM == nil {
				eventM = NewEvent(fd, p.Loop)
			}
			eventM.SetEvents(int(p.Events[i].Events))
			p.EventMap[fd] = eventM
			cb(eventM)
		}
	}
}

func (p *SocketPoll) Update(operation int, event *Event) {
	fd := event.Fd()
	if err := syscall.EpollCtl(p.epFd, operation, int(fd),
		&syscall.EpollEvent{Fd: int32(fd),
			Events: uint32(event.Events()),
		},
	); err != nil {
		fmt.Errorf(err.Error())
	}
}

func (p *SocketPoll) UpdateEvent(event *Event) {
	fd := event.Fd()
	_, ok := p.EventMap[fd]
	if !ok {
		p.EventMap[fd] = event
		p.Update(syscall.EPOLL_CTL_ADD, event)
		fmt.Println("event ADD: ", event.Events())
	} else {
		if event.IsNoEvent() {
			fmt.Println("event DEL: ", event.Events())
			p.Update(syscall.EPOLL_CTL_DEL, event)
		} else {
			fmt.Println("event MOD: ", event.Events())
			p.Update(syscall.EPOLL_CTL_MOD, event)
		}
	}
}

func (p *SocketPoll) RemoveEvent(event *Event) {
	fd := event.Fd()
	p.Update(syscall.EPOLL_CTL_DEL, event)
	delete(p.EventMap, fd)
}
