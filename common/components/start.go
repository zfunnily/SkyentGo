package components

import (
	"fmt"
	"time"
)

type MainApp struct {
	WP *WorkPool
}

var M *MainApp

func MAInst() *MainApp {
	if M == nil {
		HSInit(1)
		M = &MainApp{
			WP: NewWorkPool(4, 10),
		}
		//SS = NewServer(9081)
	}
	return M
}

func (m *MainApp) work() {
	var q *MessageQueue
	for {
		q = ContextMessageDispatch(q)
	}
}

func (m *MainApp) timer() {
	for {
		TWInst().UpdateTime()
		time.Sleep(time.Microsecond * 2500)
	}
}

func (m *MainApp) socket() {
	for {
	}
}

func (m *MainApp) Start() {
	M.WP.StartWorkerPool()
	fmt.Println("\n\nstart .... MainAPP")
	M.WP.TaskQueue[0] <- m.timer
	//M.WP.TaskQueue[1] <- m.socket
	M.WP.TaskQueue[2] <- m.work
	M.WP.TaskQueue[3] <- m.work
}
