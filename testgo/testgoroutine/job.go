package main

const (
	MaxTWorkerPoolSize = 100 //并发池最大长度
	MaxJobQueueSize    = 100 //任务队列最大长度
)

type Task interface {
	Do() error
}

type TaskChan chan Task

type TWorkerChan chan TaskChan

var (
	TaskQueue TaskChan
	TWorkPool TWorkerChan
)

type TWorker struct {
	TaskChannel TaskChan
	quit        chan bool
}

func NewTWorker() *TWorker {
	return &TWorker{
		TaskChannel: make(TaskChan),
		quit:        make(chan bool),
	}
}

func (w *TWorker) Stop() {
	w.quit <- true
}

func (w *TWorker) Start() {
	go func() {
		for {
			TWorkPool <- w.TaskChannel
			select {
			case task := <-w.TaskChannel:
				task.Do()
			case <-w.quit:
				return
			}
		}
	}()
}

type Dispather1 struct {
	TWorkers []*TWorker
	quit     chan bool
}

func (d *Dispather1) Stop() {
	for _, tw := range d.TWorkers {
		tw.Stop()
	}
	d.quit <- true
}

func (d *Dispather1) Start() {
	TWorkPool = make(TWorkerChan, MaxTWorkerPoolSize)
	TaskQueue = make(TaskChan, MaxJobQueueSize)
	d.quit = make(chan bool)

	for i := 0; i < MaxTWorkerPoolSize; i++ {
		worker := NewTWorker()
		d.TWorkers = append(d.TWorkers, worker)
		worker.Start()
	}

	for {
		select {
		case task := <-TaskQueue:
			go func(task Task) {
				taskchan := <-TWorkPool
				taskchan <- task
			}(task)
		case <-d.quit:
			return
		}
	}
}

func test1() {
	//往队列中添加任务
}
