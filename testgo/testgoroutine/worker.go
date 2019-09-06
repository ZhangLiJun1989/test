package main

import "fmt"

type Job1 interface {
	Do()
}

type Worker1 struct {
	JobQueue chan Job1
}

func NewWorker1() *Worker1 {
	return &Worker1{JobQueue: make(chan Job1)}
}

func (w Worker1) Run(wq chan chan Job1) {
	go func() {
		for {
			wq <- w.JobQueue
			select {
			case job := <-w.JobQueue:
				job.Do()
			}
		}
	}()
}

type WorkerPool1 struct {
	workerlen   int
	JobQueue    chan Job1
	WorkerQueue chan chan Job1
}

func NewWorkerPool1(num int) *WorkerPool1 {
	return &WorkerPool1{
		workerlen:   num,
		JobQueue:    make(chan Job1),
		WorkerQueue: make(chan chan Job1, num),
	}
}

func (wp *WorkerPool1) Run() {
	for i := 0; i < wp.workerlen; i++ {
		worker := NewWorker1()
		worker.Run(wp.WorkerQueue)
	}
	go func() {
		for {
			select {
			case job := <-wp.JobQueue:
				worker := <-wp.WorkerQueue
				worker <- job
			}
		}
	}()
}

func test() {
	num := 100 * 100 * 20
	p := NewWorkerPool1(num)
	p.Run()
	dataNum := 100 * 100 * 100 * 1000
	go func() {
		for i := 0; i < dataNum; i++ {
			sc := &Score1{Name: i}
			p.JobQueue <- sc
		}
	}()
}

type Score1 struct {
	Name int
}

func (s *Score1) Do() {
	fmt.Println("num: %d", s.Name)
}
