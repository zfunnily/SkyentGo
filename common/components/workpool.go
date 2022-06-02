package components

type Job func()

type WorkPool struct {
	WorkerPoolSize   int64
	MaxTaskPerWorker int64
	TaskQueue        []chan Job
}

func NewWorkPool(poolSize, maxTaskSize int64) *WorkPool {
	return &WorkPool{
		WorkerPoolSize:   poolSize,
		MaxTaskPerWorker: maxTaskSize,
		TaskQueue:        make([]chan Job, poolSize),
	}
}

//StartOneWorker 启动一个Worker工作流程
func (wp *WorkPool) StartOneWorker(workerID int, taskQueue chan Job) {
	//不断的等待队列中的消息
	for {
		select {
		//有消息则取出队列的Request，并执行绑定的业务方法
		case job := <-taskQueue:
			_ = workerID
			job()
		}
	}
}

func (wp *WorkPool) StartWorkerPool() {
	//遍历需要启动worker的数量，依此启动
	for i := 0; i < int(wp.WorkerPoolSize); i++ {
		//一个worker被启动
		//给当前worker对应的任务队列开辟空间
		wp.TaskQueue[i] = make(chan Job, wp.MaxTaskPerWorker)
		//启动当前Worker，阻塞的等待对应的任务队列是否有消息传递进来
		go wp.StartOneWorker(i, wp.TaskQueue[i])
	}
}
