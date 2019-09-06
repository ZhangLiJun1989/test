package main

import "time"

//如何对goroutine进行并发控制
const (
	MaxWorkerPoolSize = 100 //并发池最大长度
	MaxJobPoolSize    = 100 //任务队列缓冲器
)

type Job2 interface {
	ID() string
	Do() error
}

type JobChan chan Job2

type WorkerChan chan JobChan

type Worker2 struct {
	ID         int
	JobChannel JobChan
	quit       chan bool
}

func NewWorker2(id int) *Worker2 {
	return &Worker2{
		ID:         id,
		JobChannel: make(chan Job2),
		quit:       make(chan bool),
	}
}

func (w *Worker2) Stop() {
	w.quit <- true
}

func (w *Worker2) Start(wq WorkerChan) {
	go func() {
		for {
			wq <- w.JobChannel
			select {
			case job := <-w.JobChannel:
				job.Do()
			case <-w.quit:
				return
			}
			time.Sleep(time.Second * 1)
		}
	}()
}

type Dispather struct {
	Workers    []*Worker2
	JobQueue   JobChan
	WorkerPool WorkerChan
	quit       chan bool
}

func (d *Dispather) Stop() {
	for _, worker := range d.Workers {
		worker.Stop()
	}
	d.quit <- true
}

func (d *Dispather) Start() {
	d.WorkerPool = make(WorkerChan, MaxWorkerPoolSize)
	d.JobQueue = make(JobChan, MaxJobPoolSize)
	d.quit = make(chan bool)

	for i := 0; i < MaxWorkerPoolSize; i++ {
		worker := NewWorker2(i)
		d.Workers = append(d.Workers, worker)
		worker.Start(d.WorkerPool)
	}

	for {
		select {
		case job := <-d.JobQueue:
			go func(job Job2) {
				jobchan := <-d.WorkerPool
				jobchan <- job
			}(job)
		case <-d.quit:
			return
		}
	}
}

func testgggg() {

}
