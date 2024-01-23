package main

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// Task 是一个需要执行的任务函数定义
type Task func() error

// Pool 是协程池的结构体
type Pool struct {
	lock        sync.Mutex     // 用于锁定协程池的状态
	tasks       []Task         // 任务队列
	concurrency int            // 协程池的并发数
	tasksChan   chan Task      // 任务通道
	wg          sync.WaitGroup // 等待所有工作协程结束
	errors      chan error     // 错误通道，用于传递任务中产生的错误
	done        chan struct{}  // 用于通知关闭协程池的通道
}

// NewPool 创建一个新的协程池
func NewPool(concurrency, maxQueueSize int) *Pool {
	p := &Pool{
		concurrency: concurrency,
		tasksChan:   make(chan Task, maxQueueSize),
		errors:      make(chan error, maxQueueSize),
		done:        make(chan struct{}),
	}
	p.start()
	return p
}

// start 开始协程池的运行，启动固定数量的工作协程等待处理任务
func (p *Pool) start() {
	for i := 0; i < p.concurrency; i++ {
		p.wg.Add(1)
		go p.worker()
	}
}

// Stop 停止协程池的工作
func (p *Pool) Stop() {
	close(p.done)
	p.wg.Wait() // 等待所有协程结束
	close(p.errors)
}

// Exec 提交一个任务到协程池
func (p *Pool) Exec(task Task) {
	select {
	case p.tasksChan <- task:
	case <-p.done:
		p.errors <- errors.New("协程池已经停止接受新的任务")
	}
}

// worker 是协程池中的工作协程
func (p *Pool) worker() {
	defer p.wg.Done()
	for {
		select {
		case <-p.done: // 接收到停止信号
			return
		case task, ok := <-p.tasksChan: // 从任务通道获取任务
			if !ok {
				return
			}
			// 错误处理交由外部调用者
			if err := task(); err != nil {
				p.errors <- err
			}
		}
	}
}

// ErrorChan 提供错误通道供外部监听并处理错误
func (p *Pool) ErrorChan() <-chan error {
	return p.errors
}

func main() {
	// 创建并启动协程池
	pool := NewPool(5, 10) // 最大并发数为5，任务队列大小为10

	// 提交任务到协程池
	for i := 0; i < 30; i++ {
		count := i
		pool.Exec(func() error {
			fmt.Printf("当前协程数量%d\n", runtime.NumGoroutine())
			time.Sleep(time.Second) // 模拟耗时任务
			if count%7 == 0 {
				return fmt.Errorf("任务 %d 执行出错", count)
			}
			fmt.Printf("任务 %d 执行完毕\n", count)
			return nil
		})
	}

	// 监听并处理错误
	go func() {
		for err := range pool.ErrorChan() {
			fmt.Println("捕获到错误:", err)
		}
	}()

	// 等待一段时间后停止协程池
	time.Sleep(10 * time.Second)
	pool.Stop()
	fmt.Println("协程池已经停止")
}
