package reactor

import (
	"fmt"
	"syscall"
)

const (
	kInitialHeadSize = 8
	kInitialSize     = 1024
)

type Buffer struct {
	buffer     []byte
	writeIndex int64
	readIndex  int64
}

func NewBuffer() *Buffer {
	return &Buffer{
		buffer:     make([]byte, kInitialSize+kInitialHeadSize),
		writeIndex: kInitialHeadSize,
		readIndex:  kInitialHeadSize,
	}
}

func (b *Buffer) WritableBytes() int64 {
	return int64(len(b.buffer)) - b.writeIndex
}

func (b *Buffer) ReadableBytes() int64 {
	return b.writeIndex - b.readIndex
}
func (b *Buffer) Peek() []byte {
	return b.Begin(b.readIndex)
}

func (b *Buffer) RetrieveAll() {
	b.readIndex = kInitialHeadSize
	b.writeIndex = kInitialHeadSize
}

func (b *Buffer) Retrieve(len int64) {
	if len < b.ReadableBytes() {
		b.readIndex += len
	} else {
		b.RetrieveAll()
	}
}

func (b *Buffer) Begin(idx int64) []byte {
	return b.buffer[idx:]
}

func (b *Buffer) BeginWrite() []byte {
	return b.buffer[b.writeIndex:]
}

func (b *Buffer) ReadFd(fd int) int64 {
	var n, l int
	var err error
	for {
		l, err = syscall.Read(fd, b.BeginWrite())
		if err != nil {
			fmt.Errorf("read error: %s\n", err.Error())
			break
		}
		if l <= 0 {
			break
		}

		if int64(l) <= b.WritableBytes() {
			b.writeIndex = b.writeIndex + int64(l)
		} else {
			b.writeIndex = int64(len(b.buffer))
		}

		n = n + l
	}

	return int64(n)
}

func (b *Buffer) Append(buf []byte, l int64) {
	if l > b.WritableBytes() {
		copy(b.buffer[b.writeIndex:], buf[:b.WritableBytes()])
		b.buffer = append(b.buffer, buf[b.WritableBytes():l]...)
	} else {
		copy(b.buffer[b.writeIndex:], buf[:l])
	}
	b.writeIndex = b.writeIndex + l
}
