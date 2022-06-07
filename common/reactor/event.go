package reactor

import (
	"fmt"
	"syscall"
)

type EventCbk func()
type ReadEventCbk func()

type Event struct {
	fd      int
	events  int
	revents int

	ReadCbk  ReadEventCbk
	WriteCbk EventCbk
	CloseCbk EventCbk
	ErrorCbk EventCbk

	Loop *EventLoop
}

func NewEvent(fd int, loop *EventLoop) *Event {
	return &Event{fd: fd, Loop: loop, events: 0, revents: 0}
}

func (e *Event) SetReadCbk(cbk ReadEventCbk) {
	e.ReadCbk = cbk
}

func (e *Event) SetWriteCbk(cbk EventCbk) {
	e.WriteCbk = cbk
}

func (e *Event) SetCloseCbk(cbk EventCbk) {
	e.CloseCbk = cbk
}

func (e *Event) SetErrorCbk(cbk EventCbk) {
	e.ErrorCbk = cbk
}

func (e *Event) Update() {
	e.Loop.UpdateEvent(e)
}

func (e *Event) Remove() {
	e.Loop.RemoveEvent(e)
}

func (e *Event) Fd() int {
	return e.fd
}

func (e *Event) IsNoEvent() bool {
	return e.events == 0
}

func (e *Event) Events() int {
	return e.events
}

func (e *Event) SetEvents(event int) {
	e.revents = event
}

func (e *Event) EnableConnecting() {
	e.events |= syscall.EPOLLIN
	e.Update()
}

func (e *Event) EnableReading() {
	e.events |= syscall.EPOLLIN | syscall.EPOLLPRI | syscall.EPOLLHUP
	e.Update()
}

func (e *Event) DisableReading() {
	e.events &= ^(syscall.EPOLLIN | syscall.EPOLLPRI | syscall.EPOLLHUP)
	e.Update()
}

func (e *Event) EnableWriting() {
	e.events |= syscall.EPOLLOUT
	e.Update()
}

func (e *Event) DisableWriting() {
	e.events &= ^syscall.EPOLLOUT
	e.Update()
}

func (e *Event) DisableAll() {
	e.events = 0
	e.Update()
}

func (e *Event) IsWriting() int {
	return e.events & (syscall.EPOLLOUT)
}

func (e *Event) IsReading() int {
	return e.events & (syscall.EPOLLIN | syscall.EPOLLPRI)
}

func (e *Event) HandleEvents() {
	if (e.revents&syscall.EPOLLHUP) > 0 && (e.revents&syscall.EPOLLIN) <= 0 {
		// close
		fmt.Println("handle events close: ", e.Fd())
		if e.CloseCbk != nil {
			e.CloseCbk()
		}
	}

	if e.revents&(syscall.EPOLLIN|syscall.EPOLLPRI|syscall.EPOLLRDHUP) > 0 {
		// read
		fmt.Println("handle events read: ", e.Fd())
		if e.ReadCbk != nil {
			e.ReadCbk()
		}
	}

	if e.revents&syscall.EPOLLOUT > 0 {
		// write
		fmt.Println("handle events write: ", e.Fd())
		if e.WriteCbk != nil {
			e.WriteCbk()
		}
	}

	if e.revents&syscall.EPOLLERR > 0 {
		// error
		if e.ErrorCbk != nil {
			e.ErrorCbk()
		}
	}
}
