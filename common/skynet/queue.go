package skynet

import (
	"container/list"
	"sync"
)

const (
	SIZE_MAX           = 0xffffffffffffffff64
	MESSAGE_TYPE_MASK  = (SIZE_MAX >> 8)
	MESSAGE_TYPE_SHIFT = 56
)

type Message struct {
	Source  uint32
	Session int
	Data    interface{}
	Typ     int
}

type MessageQueue struct {
	mutex sync.RWMutex
	*list.List
	handle uint32
}

func NewMessageQueue(handle uint32) *MessageQueue {
	return &MessageQueue{
		mutex:  sync.RWMutex{},
		List:   list.New(),
		handle: handle,
	}
}

func (m *MessageQueue) Push(msg *Message) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.PushBack(msg)
}

func (m *MessageQueue) Pop() *Message {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	msg := m.Front()
	if msg == nil {
		return nil
	}
	return m.Remove(msg).(*Message)
}

func (m *MessageQueue) Handle() uint32 {
	return m.handle
}

type GlobalQueue struct {
	*list.List
	mutex sync.RWMutex
}

var Q *GlobalQueue

func GQInst() *GlobalQueue {
	if Q == nil {
		Q = &GlobalQueue{
			mutex: sync.RWMutex{},
			List:  list.New(),
		}
	}
	return Q
}

func (g *GlobalQueue) Push(queue *MessageQueue) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.PushBack(queue)
}

func (g *GlobalQueue) Pop() *MessageQueue {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	f := g.Front()
	if f == nil {
		return nil
	}
	return g.Remove(f).(*MessageQueue)
}
