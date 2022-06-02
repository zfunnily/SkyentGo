package components

import (
	"container/list"
	"sync"
)

const (
	DefaultSlotSize = 4
	HANDLE_MASK     = 0xffffff
)

type HandleName struct {
	name   string
	handle uint32
}

type HandleStorage struct {
	mutex       sync.RWMutex
	harbor      uint32
	handleIndex uint32
	slotSize    uint32
	slot        []*Context
	Name        []*HandleName
	retire      *list.List
}

var H *HandleStorage

func HSInst() *HandleStorage {
	return H
}

func HSInit(harbor int) *HandleStorage {
	if H == nil {
		H = &HandleStorage{
			mutex:       sync.RWMutex{},
			harbor:      uint32(harbor),
			handleIndex: 0,
			slotSize:    DefaultSlotSize,
			slot:        make([]*Context, DefaultSlotSize),
			Name:        make([]*HandleName, 2),
			retire:      list.New(),
		}
	}
	return H
}

func (h *HandleStorage) Register(ctx *Context) uint32 {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	handle := h.handleIndex + 1
	hashV := h.retire.Front()
	if hashV != nil {
		handle = h.retire.Remove(hashV).(uint32)
		h.slot[handle] = ctx
	} else {
		if handle <= uint32(len(h.slot)-1) {
			h.slot[handle] = ctx
		} else {
			h.slot = append(h.slot, ctx)
		}
		h.handleIndex = handle
	}
	return handle
}

func (h *HandleStorage) Retire(handle uint32) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	ctx := h.slot[handle]
	if ctx != nil {
		h.slot[handle] = nil
		h.retire.PushBack(handle)
	}
}

func (h *HandleStorage) Grab(handle uint32) *Context {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	return h.slot[handle]
}

func (h *HandleStorage) RetireAll() {
}

func (h *HandleStorage) FindName(name string) {
}

func (h *HandleStorage) NameHandle(handle uint32, name string) {
}
