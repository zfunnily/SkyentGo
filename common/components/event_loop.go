package components

import "fmt"

type LoopFunc func()

type EventLoop struct {
	Quit   bool
	Poller *SocketPoll
}

func NewEventLoop() *EventLoop {
	loop := &EventLoop{
		Quit: true,
	}
	loop.Poller = NewPoll(loop)
	return loop
}

func (e *EventLoop) Loop() {
	e.Quit = false
	for !e.Quit {
		err := e.Poller.Poll(func(event *Event) {
			if event != nil {
				event.HandleEvents()
			}
		})
		if err != nil {
			break
		}
	}
	fmt.Println("loop exit")
}

func (e *EventLoop) RunInLoop(cb LoopFunc) {
	cb()
}

func (e *EventLoop) UpdateEvent(event *Event) {
	e.Poller.UpdateEvent(event)
}

func (e *EventLoop) RemoveEvent(event *Event) {
	e.Poller.RemoveEvent(event)
}
