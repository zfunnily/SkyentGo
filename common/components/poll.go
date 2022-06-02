package components

import "syscall"

type SocketPoll struct {
	sockFd   int
	epFd     int
	Events   []syscall.EpollEvent
	EventMap map[int32]*Event
}
