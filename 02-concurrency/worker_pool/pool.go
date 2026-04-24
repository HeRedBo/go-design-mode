// Package workerpool 演示 Go 语言的工作池模式 (Worker Pool Pattern)
//
// 工作池模式解决的问题：
// 1. 控制并发数量，避免创建过多 goroutine 耗尽资源
// 2. 任务排队处理，实现负载均衡
// 3. 复用 goroutine，减少创建/销毁开销
//
// 使用场景：
// - 批量处理任务（数据处理、文件处理、网络请求等）
// - 限制并发数保护下游服务
// - CPU 密集型任务并行处理
//
// 核心组件：
// - jobs channel: 任务队列
// - results channel: 结果队列
// - workers: 固定数量的 goroutine 消费任务
// - WaitGroup: 等待所有任务完成
//
// 工作流程：
// 1. 创建 jobs 和 results channel
// 2. 启动 N 个 worker goroutine
// 3. 生产者发送任务到 jobs channel
// 4. worker 从 jobs 获取任务，处理后发送到 results
// 5. 关闭 jobs channel，等待所有任务完成
// 6. 关闭 results channel，收集所有结果
package workerpool

import (
	"fmt"
	"sync"
)

// Task 代表一个任务
type Task struct {
	ID      int                                    // 任务 ID
	Data    interface{}                            // 任务数据
	Handler func(interface{}) (interface{}, error) // 任务处理函数
}

// Result 代表任务处理结果
type Result struct {
	TaskID int         // 任务 ID
	Output interface{} // 处理结果
	Err    error       // 错误信息
}

// WorkerPool 工作池
type WorkerPool struct {
	numWorkers int            // worker 数量
	tasks      chan Task      // 任务队列
	results    chan Result    // 结果队列
	wg         sync.WaitGroup // 等待组
}

// New 创建新的工作池
// numWorkers: worker 数量
// bufferSize: channel 缓冲区大小（可选，设为 0 表示无缓冲）
func New(numWorkers int, bufferSize int) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		tasks:      make(chan Task, bufferSize),
		results:    make(chan Result, bufferSize),
	}
}

// Start 启动工作池
// 启动指定数量的 worker goroutine
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// worker 是工作协程
// id: worker 编号（用于调试）
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	// 持续从 tasks channel 获取任务，直到 channel 关闭
	for task := range wp.tasks {
		// 使用 recover 捕获 panic，防止单个任务影响整个工作池
		func() {
			defer func() {
				if r := recover(); r != nil {
					// 发送 panic 错误到 results channel
					wp.results <- Result{
						TaskID: task.ID,
						Output: nil,
						Err:    fmt.Errorf("task %d panicked: %v", task.ID, r),
					}
				}
			}()

			// 执行任务
			output, err := task.Handler(task.Data)

			// 发送结果到 results channel
			wp.results <- Result{
				TaskID: task.ID,
				Output: output,
				Err:    err,
			}
		}()
	}
}

// Submit 提交任务到工作池
// 如果 tasks channel 已满，会阻塞直到有空间
func (wp *WorkerPool) Submit(task Task) {
	wp.tasks <- task
}

// SubmitAsync 异步提交任务
// 返回 false 表示 channel 已满，任务未被提交
func (wp *WorkerPool) SubmitAsync(task Task) bool {
	select {
	case wp.tasks <- task:
		return true
	default:
		return false
	}
}

// Close 关闭任务提交
// 关闭 tasks channel，通知所有 worker 没有更多任务
// 必须在所有任务提交完成后调用
func (wp *WorkerPool) Close() {
	close(wp.tasks)
}

// Wait 等待所有任务完成
// 阻塞直到所有 worker 处理完所有任务
func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
	close(wp.results)
}

// Results 返回结果 channel
// 调用者可以 range 这个 channel 获取所有结果
func (wp *WorkerPool) Results() <-chan Result {
	return wp.results
}

// NumWorkers 返回 worker 数量
func (wp *WorkerPool) NumWorkers() int {
	return wp.numWorkers
}

// ProcessBatch 批量处理任务的便捷方法
// 自动创建、启动、提交、等待和收集结果
//
// 示例：
//
//	tasks := []Task{...}
//	results := ProcessBatch(4, tasks)
func ProcessBatch(numWorkers int, tasks []Task) []Result {
	wp := New(numWorkers, len(tasks))
	wp.Start()

	// 提交所有任务
	for _, task := range tasks {
		wp.Submit(task)
	}
	wp.Close()

	// 启动收集结果的 goroutine
	go wp.Wait()

	// 收集所有结果
	var results []Result
	for result := range wp.Results() {
		results = append(results, result)
	}

	return results
}

// ============ 示例函数 ============

// ExampleProcessor 示例：处理任务并返回结果
func ExampleProcessor(data interface{}) (interface{}, error) {
	num, ok := data.(int)
	if !ok {
		return nil, fmt.Errorf("expected int, got %T", data)
	}

	// 模拟处理
	result := num * num
	return result, nil
}
