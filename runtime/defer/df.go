package main

import (
	"fmt"
)

// Worker 表示一个工作协程
type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
}

// NewWorker 创建一个 Worker
func NewWorker(workerPool chan chan Job) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool)}
}

// Start 启动 Worker
func (w Worker) Start() {
	go func() {
		for {
			// 将当前的 Worker 注册到 Worker 队列中
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				// 接收到一个工作请求
				if err := job.Do(); err != nil {
					fmt.Println("Error when handling job:", err.Error())
				}

			case <-w.quit:
				// 收到停止信号
				return
			}
		}
	}()
}

// Stop 停止 Worker
func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

// Job 表示一个工作任务
type Job struct {
	Payload Payload
}

// Do 执行 Job
func (j *Job) Do() error {
	// 执行具体的任务
	fmt.Println("execute job:", j.Payload)
	return nil
}

// Payload 表示 Job 的负载
type Payload struct {
	Content string
}

// Dispatcher 表示调度器
type Dispatcher struct {
	WorkerPool chan chan Job
	MaxWorkers int
}

// NewDispatcher 创建一个 Dispatcher
func NewDispatcher(maxWorkers int) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{WorkerPool: pool, MaxWorkers: maxWorkers}
}

// Run 启动 Dispatcher
func (d *Dispatcher) Run() {
	// 初始化指定数量的 Worker
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker(d.WorkerPool)
		worker.Start()
	}

	go d.dispatch()
}

// dispatch 负责协调、分发任务
func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-JobQueue:
			// 接收到一个工作请求
			go func(job Job) {
				// 尝试获取一个可用的 worker
				jobChannel := <-d.WorkerPool

				// 将工作分发给该 worker
				jobChannel <- job
			}(job)
		}
	}
}

// JobQueue 表示 Job 队列
var JobQueue chan Job

func main() {
	// 设置最大的 Worker 数量
	dispatcher := NewDispatcher(100)
	dispatcher.Run()

	// 将任务添加到 JobQueue 中
	JobQueue = make(chan Job, 10)
	job := Job{Payload: Payload{Content: "Hello, World!"}}
	JobQueue <- job
}
