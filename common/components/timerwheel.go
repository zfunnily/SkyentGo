package components

import (
	"container/list"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

//skynet的时间轮 + 协程池
const (
	TimeNearShift  = 8
	TimeNear       = 1 << TimeNearShift
	TimeLevelShift = 6
	TimeLevel      = 1 << TimeLevelShift
	TimeNearMask   = TimeNear - 1
	TimeLevelMask  = TimeLevel - 1

	//协程池 大小
	WorkerPoolSize   = 10
	MaxTaskPerWorker = 20
)

type bucket struct {
	timers *list.List
	mu     sync.Mutex
}

func newBucket() *bucket {
	return &bucket{
		timers: list.New(),
		mu:     sync.Mutex{},
	}
}

func (b *bucket) Add(t *TimerEvent) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.timers.PushBack(t)
}

func (b *bucket) Flush(reinsert func(t *TimerEvent)) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for e := b.timers.Front(); e != nil; {
		next := e.Next()
		reinsert(e.Value.(*TimerEvent))

		b.timers.Remove(e)
		e = next
	}
}

type TimerEvent struct {
	expiration uint32
	handle     uint32
	session    int
}

var TimingWheel *TimeWheel

func TWInst() *TimeWheel {
	if TimingWheel == nil {
		TimingWheel = NewTimeWheel()
	}
	return TimingWheel
}

type TimeWheel struct {
	tick   time.Duration
	ticker *time.Ticker
	near   [TimeNear]*bucket
	t      [4][TimeLevel]*bucket
	time   uint32 // cs 厘s

	WorkPool *WorkPool
	exit     chan struct{}
	exitFlag uint32

	startTime    int64 // s
	current      int   // cs
	currentPoint int64 // cs
}

func NewTimeWheel() *TimeWheel {
	tw := &TimeWheel{
		tick:     10 * time.Millisecond,
		time:     0,
		WorkPool: NewWorkPool(WorkerPoolSize, MaxTaskPerWorker),
		exit:     make(chan struct{}),
		exitFlag: 0,
	}
	for i := 0; i < TimeNear; i++ {
		tw.near[i] = newBucket()
	}

	for i := 0; i < 4; i++ {
		for j := 0; j < TimeLevel; j++ {
			tw.t[i][j] = newBucket()
		}
	}

	now := time.Now()
	tw.startTime = now.Unix()
	tw.currentPoint = now.UnixMilli() / 10
	return tw
}

func (tw *TimeWheel) add(t *TimerEvent) bool {
	time := t.expiration
	currentTime := atomic.LoadUint32(&tw.time)
	if time <= currentTime {
		return false
	}

	if (time | TimeNearMask) == (currentTime | TimeNearMask) {
		tw.near[time&TimeNearMask].Add(t)
	} else {
		i := 0
		mask := TimeNear << TimeNearShift
		for i = 0; i < 3; i++ {
			if (time | uint32(mask-1)) == (currentTime | uint32(mask-1)) {
				break
			}
			mask <<= TimeLevelShift
		}

		tw.t[i][((time >> (TimeNearShift + i*TimeLevelShift)) & TimeLevelMask)].Add(t)
	}
	return true
}

func (tw *TimeWheel) addOrRun(t *TimerEvent) {
	if !tw.add(t) {
		msg := &Message{
			Source:  0,
			Session: t.session,
			Data:    nil,
			Typ:     PTYPE_RESPONSE,
		}
		fmt.Println("add or run run run ...")
		ContextPush(t.handle, msg)
	}
}

func (tw *TimeWheel) moveList(level, idx int) {
	current := tw.t[level][idx]
	current.Flush(tw.addOrRun)
}

func (tw *TimeWheel) shift() {
	mask := TimeNear
	ct := atomic.AddUint32(&tw.time, 1)
	if ct == 0 {
		tw.moveList(3, 0)
	} else {
		time := ct >> TimeNearShift

		i := 0
		for (ct & uint32(mask-1)) == 0 {
			idx := time & TimeLevelMask
			if idx != 0 {
				tw.moveList(i, int(idx))
				break
			}

			mask <<= TimeLevelShift
			time >>= TimeLevelShift
			i++
		}
	}
}

func (tw *TimeWheel) execute() {
	idx := tw.time & TimeNearMask
	tw.near[idx].Flush(tw.addOrRun)
}

func (tw *TimeWheel) update() {
	tw.execute()
	tw.shift()
	tw.execute()
}

func (tw *TimeWheel) UpdateTime() {
	cp := time.Now().UnixMilli() / 10
	if cp < tw.currentPoint {
		tw.currentPoint = cp
	} else {
		diff := cp - tw.currentPoint
		tw.currentPoint = cp
		tw.current = tw.current + int(diff)

		for i := 0; i < int(diff); i++ {
			tw.update()
		}
	}
}

func (tw *TimeWheel) Stop() {
	flag := atomic.LoadUint32(&tw.exitFlag)
	if flag != 0 {
		return
	}

	atomic.StoreUint32(&tw.exitFlag, 1)
	close(tw.exit)
}

func (tw *TimeWheel) TimeOut(handle uint32, time int, session int) int {
	if time <= 0 {
		msg := &Message{
			Source:  0,
			Session: session,
			Data:    nil,
			Typ:     PTYPE_RESPONSE,
		}
		if ContextPush(handle, msg) != 0 {
			return -1
		}
	} else {
		atotime := atomic.LoadUint32(&tw.time)
		tw.addOrRun(&TimerEvent{
			expiration: atotime + uint32(time),
			handle:     handle,
			session:    session,
		})
	}
	return session
}

func StopTimer() {
	TimingWheel.Stop()
}
