package components

import "syscall"

type SocketPoll struct {
	sockFd   int
	epFd     int
	Events   []syscall.EpollEvent
	EventMap map[int]*Event

	Loop *EventLoop
}
