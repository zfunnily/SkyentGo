package components

import "syscall"

type EventCbk func()
type ReadCbk func()

type Event struct {
	fd      int32
	events  uint32
	revents uint32

	ReadCbk  ReadCbk
	WriteCbk EventCbk
	CloseCbk EventCbk
	ErrorCbk EventCbk

	Loop EventLoop
}

func NewEvent(fd int32) *Event {
	return &Event{fd: fd}
}

func (e *Event) Update() {
	e.Loop.UpdateEvent(e)
}

func (e *Event) Remove() {
}

func (e *Event) Fd() int32 {
	return e.fd
}

func (e *Event) IsNoEvent() bool {
	return e.events == 0
}

func (e *Event) Events() uint32 {
	return e.events
}

func (e *Event) SetEvents(event uint32) {
	e.revents = event
}

func (e *Event) IsWriting() uint32 {
	return e.events & (syscall.EPOLLOUT)
}

func (e *Event) IsReading() uint32 {
	return e.events & (syscall.EPOLLIN | syscall.EPOLLPRI)
}

func (e *Event) EnableConnecting() {
	e.events |= syscall.EPOLLIN
	e.Update()
}

func (e *Event) EnableReading() {
	event := syscall.EPOLLIN | syscall.EPOLLPRI | syscall.EPOLLET
	e.events |= uint32(event)
	e.Update()
}

func (e *Event) DisableReading() {
	event := syscall.EPOLLIN | syscall.EPOLLPRI | syscall.EPOLLET
	e.events &= uint32(event)
	e.Update()
}

func (e *Event) EnableWriting() {
	e.events |= syscall.EPOLLOUT
	e.Update()
}

func (e *Event) DisableWriting() {
	e.events &= syscall.EPOLLOUT
	e.Update()
}

func (e *Event) HandleEvents() {
	if (e.revents&syscall.EPOLLHUP) > 0 && (e.revents&syscall.EPOLLIN) <= 0 {
		// close
	}

	if e.revents&(syscall.EPOLLIN|syscall.EPOLLPRI|syscall.EPOLLRDHUP) > 0 {
		// read
	}

	if e.revents&syscall.EPOLLOUT > 0 {
		// write
	}

	if e.revents&syscall.EPOLLERR > 0 {
		// error
	}
}
