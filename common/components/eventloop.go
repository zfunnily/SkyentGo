package components

type EventLoop struct {
	Looping bool
	Quit    bool
	Poller  *SocketPoll
}

func (e *EventLoop) Loop() {
	e.Quit = false
	e.Looping = true
	for e.Quit {
		err := e.Poller.Poll(func(event *Event) {
			event.HandleEvents()
		})
		if err != nil {
			break
		}
	}
}

func (e *EventLoop) RunInLoop() {
}

func (e *EventLoop) UpdateEvent(event *Event) {
	e.Poller.UpdateEvent(event)
}

func (e *EventLoop) RemoveEvent(event *Event) {
	e.Poller.RemoveEvent(event)
}
