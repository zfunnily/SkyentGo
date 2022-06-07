package reactor

type EventLoopPool struct {
	BaseLoop *EventLoop
	name     string
	Started  bool
	Num      int
	next     int
	Loops    []*EventLoop
}

func NewEventLoopPool(loop *EventLoop, name string) *EventLoopPool {
	return &EventLoopPool{
		BaseLoop: loop,
		name:     name,
		Started:  false,
		Num:      0,
		next:     0,
	}
}

func (lp *EventLoopPool) SetNum(num int) {
	lp.Num = num
}

func (lp *EventLoopPool) Start() {
	if lp.Num < 1 {
		lp.Num = 1
	} else {
		lp.Loops = make([]*EventLoop, lp.Num)
	}

	for i := 0; i < lp.Num; i++ {
		lp.Loops[i] = NewEventLoop()
		go lp.Loops[i].Loop()
	}
}

func (lp *EventLoopPool) GetNextLoop() *EventLoop {
	loop := lp.BaseLoop
	if len(lp.Loops) != 0 {
		lp.next++
		loop = lp.Loops[lp.next]
		if lp.next >= len(lp.Loops) {
			lp.next = 0
		}
	}
	return loop
}

func (lp *EventLoopPool) GetHashLoop() {
}

func (lp *EventLoopPool) GetStarted() bool {
	return lp.Started
}

func (lp *EventLoopPool) GetAllLoops() []*EventLoop {
	return lp.Loops
}

func (lp *EventLoopPool) Name() string {
	return lp.name
}
