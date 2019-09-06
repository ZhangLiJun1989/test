package main

import (
	"fmt"
	"runtime"
	"time"
)

type Score struct {
	Num int
}

func (s *Score) Do() {
	fmt.Println("num:", s.Num)
	time.Sleep(1 * 1 * time.Second)
}

// Job 数据接口，所有数据都要实现该接口，才能被传递进来
// 实现Job接口的一个数据实例，需要实现一个Do()方法，对数据的处理就在这个Do()方法中
type Job interface {
	Do()
}

// 设置两个Job通道
// 1. WorkPool 的 Job 通道，调用者把具体的数据写入这里，WorkPool 读取
// 2. Worker 的 Job 通道，当 WorkPool 读取到 Job，并拿到可用的 Worker 时，会将 Job 实例写入该 Worker 的 Job 通道，用来执行Do()方法

// Worker Worker
// 1. 每一个被初始化的 worker 都会在后期单独占用一个协程
// 2. 初始化的时候会先把自己的 JobQueue 传递到 Worker 通道中
// 3. 然后阻塞读取自己的 JobQueue，读到一个 Job 就执行 Job 对象的Do()方法
type Worker struct {
	JobQueue chan Job
}

// NewWorker -> Worker
func NewWorker() Worker {
	return Worker{JobQueue: make(chan Job)}
}

// Run Worker
func (w Worker) Run(wq chan chan Job) {
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

// WorkerPool 工作池
// 1. 初始化时会按照传入的num，启动num个后台协程，然后循环读取 Job 通道里面的数据
// 2. 读取到一个数据时，获取一个可用的 Worker ，并将 Job 对象传递到该 Worker 的 Job 通道
type WorkerPool struct {
	workerlen   int      //同时存在的 Worker 个数
	JobQueue    chan Job //WorkPool 的 Job 通道
	WorkerQueue chan chan Job
}

// NewWorkerPool NewWorkerPool
func NewWorkerPool(workerlen int) *WorkerPool {
	return &WorkerPool{
		workerlen:   workerlen,
		JobQueue:    make(chan Job),
		WorkerQueue: make(chan chan Job, workerlen),
	}
}

// Run WorkerPool
func (wp *WorkerPool) Run() {
	fmt.Println("初始化worker")
	//初始化worker
	for i := 0; i < wp.workerlen; i++ {
		worker := NewWorker()
		worker.Run(wp.WorkerQueue)
	}
	// 循环获取可用的worker,往worker中写job
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

func main() {
	// dataChan := make(chan int, 100)
	// go func() {
	// 	for {
	// 		select {
	// 		case data := <-dataChan:
	// 			fmt.Println("data:", data)
	// 			time.Sleep(1 * time.Second)
	// 		}
	// 	}
	// }()

	// //填充数据
	// for i := 0; i < 100; i++ {
	// 	dataChan <- i
	// }

	// for {
	// 	fmt.Println("runtime.NumGoroutine():", runtime.NumGoroutine())
	// 	time.Sleep(2 * time.Second)
	// }

	// 示例二
	num := 100 * 100 * 20
	// debug.SetMaxThreads(num + 1000) //设置最大线程数
	// 注册工作池，传入任务
	// 参数1 worker并发个数
	p := NewWorkerPool(num)
	p.Run()
	datanum := 100 * 100 * 100 * 100
	go func() {
		for i := 1; i <= datanum; i++ {
			sc := &Score{Num: i}
			p.JobQueue <- sc
		}
	}()

	for {
		fmt.Println("runtime.NumGoroutine() :", runtime.NumGoroutine())
		time.Sleep(2 * time.Second)
	}

}
