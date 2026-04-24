// Package fanoutin 演示 Fan-Out/Fan-In 并发模式
//
// Fan-Out/Fan-In 模式解决的问题：
// 1. Fan-Out: 将任务分发给多个 worker 并行处理
// 2. Fan-In: 将多个 worker 的结果汇聚到一个 channel
// 3. 提高吞吐量，充分利用多核
//
// 使用场景：
// - 批量数据处理
// - 并发网络请求
// - 并行计算
// - 数据聚合
//
// 工作流程：
// 1. 创建输入 channel
// 2. Fan-Out: 启动多个 goroutine 从输入 channel 读取
// 3. 每个 goroutine 处理数据并写入输出 channel
// 4. Fan-In: 将所有输出汇聚到一个 channel
// 5. 收集最终结果
package fanoutin

import (
	"context"
	"sync"
)

// FanOutIn Fan-Out/Fan-In 处理器
type FanOutIn struct {
	numWorkers int
}

// New 创建 Fan-Out/Fan-In 处理器
func New(numWorkers int) *FanOutIn {
	return &FanOutIn{
		numWorkers: numWorkers,
	}
}

// Process 执行 Fan-Out/Fan-In 处理
// input: 输入 channel
// processor: 处理函数
// 返回输出 channel
func (foi *FanOutIn) Process(ctx context.Context, input <-chan interface{}, processor func(interface{}) interface{}) <-chan interface{} {
	// Fan-Out: 启动多个 worker
	channels := make([]<-chan interface{}, foi.numWorkers)
	for i := 0; i < foi.numWorkers; i++ {
		channels[i] = foi.worker(ctx, input, processor)
	}

	// Fan-In: 汇聚所有输出
	return foi.merge(ctx, channels)
}

// worker 单个工作协程
func (foi *FanOutIn) worker(ctx context.Context, input <-chan interface{}, processor func(interface{}) interface{}) <-chan interface{} {
	output := make(chan interface{})

	go func() {
		defer close(output)
		for {
			select {
			case <-ctx.Done():
				return
			case val, ok := <-input:
				if !ok {
					return
				}
				// 处理数据
				result := processor(val)
				select {
				case output <- result:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return output
}

// merge 将多个 channel 汇聚到一个
func (foi *FanOutIn) merge(ctx context.Context, channels []<-chan interface{}) <-chan interface{} {
	output := make(chan interface{})
	var wg sync.WaitGroup

	// 为每个输入 channel 启动一个 goroutine
	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan interface{}) {
			defer wg.Done()
			for val := range c {
				select {
				case output <- val:
				case <-ctx.Done():
					return
				}
			}
		}(ch)
	}

	// 等待所有 goroutine 完成后关闭输出
	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}

// ProcessSlice 便捷函数：处理切片
func ProcessSlice(ctx context.Context, data []interface{}, numWorkers int, processor func(interface{}) interface{}) []interface{} {
	// 创建输入 channel
	input := make(chan interface{}, len(data))
	for _, item := range data {
		input <- item
	}
	close(input)

	// 执行 Fan-Out/Fan-In
	foi := New(numWorkers)
	output := foi.Process(ctx, input, processor)

	// 收集结果
	var results []interface{}
	for val := range output {
		results = append(results, val)
	}

	return results
}
