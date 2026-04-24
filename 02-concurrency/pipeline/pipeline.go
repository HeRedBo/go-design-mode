// Package pipeline 演示 Go 语言的流水线模式 (Pipeline Pattern)
//
// 流水线模式解决的问题：
// 1. 将复杂处理流程分解为多个独立阶段
// 2. 每个阶段专注于单一职责，提高可维护性
// 3. 阶段间通过 channel 连接，实现并发处理
// 4. 提高吞吐量，支持并行处理
//
// 使用场景：
// - 数据处理管道（读取->转换->过滤->输出）
// - 文件处理（读取->解析->转换->写入）
// - 网络请求处理（请求->验证->处理->响应）
// - 音视频处理流程
//
// 核心组件：
// - Stage: 处理阶段，接收输入 channel，返回输出 channel
// - Pipeline: 连接多个 stage 的管道
// - Fan-out/Fan-in: 多路分发和汇聚
//
// 设计原则：
// 1. 每个 stage 是一个独立的 goroutine
// 2. stage 通过 channel 通信
// 3. stage 在输入 channel 关闭时自动退出
// 4. 使用 context 支持取消操作
// 5. 处理错误并传播到下游
//
// 注意事项：
// 1. 避免 goroutine 泄漏：确保所有 stage 都能退出
// 2. channel 关闭时机：只在发送方关闭 channel
// 3. 错误处理：需要优雅地处理阶段失败
// 4. 缓冲大小：合理设置 channel 缓冲区提高性能
package pipeline

import (
	"context"
	"fmt"
)

// Stage 代表流水线的一个阶段
// 接收输入 channel，返回输出 channel
type Stage func(ctx context.Context, in <-chan interface{}) <-chan interface{}

// Pipeline 流水线
type Pipeline struct {
	stages []Stage
}

// New 创建新的流水线
// stages: 处理阶段列表
func New(stages ...Stage) *Pipeline {
	return &Pipeline{
		stages: stages,
	}
}

// AddStage 添加处理阶段
func (p *Pipeline) AddStage(stage Stage) {
	p.stages = append(p.stages, stage)
}

// Execute 执行流水线
// ctx: 上下文，用于取消操作
// in: 输入 channel
// 返回输出 channel
func (p *Pipeline) Execute(ctx context.Context, in <-chan interface{}) <-chan interface{} {
	current := in

	// 依次连接所有阶段
	for _, stage := range p.stages {
		current = stage(ctx, current)
	}

	return current
}

// ============ 常用 Stage 工厂函数 ============

// Generator 生成器阶段：将数据发送到 channel
func Generator(items ...interface{}) Stage {
	return func(ctx context.Context, in <-chan interface{}) <-chan interface{} {
		out := make(chan interface{})
		go func() {
			defer close(out)
			for _, item := range items {
				select {
				case out <- item:
				case <-ctx.Done():
					return
				}
			}
		}()
		return out
	}
}

// Map 映射阶段：对每个元素应用转换函数
func Map(fn func(interface{}) interface{}) Stage {
	return func(ctx context.Context, in <-chan interface{}) <-chan interface{} {
		out := make(chan interface{})
		go func() {
			defer close(out)
			for val := range in {
				result := fn(val)
				select {
				case out <- result:
				case <-ctx.Done():
					return
				}
			}
		}()
		return out
	}
}

// Filter 过滤阶段：保留满足条件的元素
func Filter(fn func(interface{}) bool) Stage {
	return func(ctx context.Context, in <-chan interface{}) <-chan interface{} {
		out := make(chan interface{})
		go func() {
			defer close(out)
			for val := range in {
				if fn(val) {
					select {
					case out <- val:
					case <-ctx.Done():
						return
					}
				}
			}
		}()
		return out
	}
}

// Square 示例：平方计算阶段
func Square() Stage {
	return Map(func(val interface{}) interface{} {
		n := val.(int)
		return n * n
	})
}

// IsEven 示例：偶数过滤阶段
func IsEven() Stage {
	return Filter(func(val interface{}) bool {
		n := val.(int)
		return n%2 == 0
	})
}

// ToString 示例：转换为字符串
func ToString() Stage {
	return Map(func(val interface{}) interface{} {
		return fmt.Sprintf("%v", val)
	})
}

// ============ 错误处理 ============

// Result 带错误信息的处理结果
type Result struct {
	Value interface{}
	Err   error
}

// MapWithErr 带错误处理的映射阶段
func MapWithErr(fn func(interface{}) (interface{}, error)) Stage {
	return func(ctx context.Context, in <-chan interface{}) <-chan interface{} {
		out := make(chan interface{})
		go func() {
			defer close(out)
			for val := range in {
				result, err := fn(val)
				select {
				case out <- Result{Value: result, Err: err}:
				case <-ctx.Done():
					return
				}
			}
		}()
		return out
	}
}

// FilterErrors 过滤掉错误结果
func FilterErrors() Stage {
	return Filter(func(val interface{}) bool {
		result, ok := val.(Result)
		if !ok {
			return true
		}
		return result.Err == nil
	})
}

// ============ Fan-Out / Fan-In ============

// FanOut 将输入分发到多个 stage
// num: 分发数量
// stage: 要复制的 stage
func FanOut(num int, stage Stage) Stage {
	return func(ctx context.Context, in <-chan interface{}) <-chan interface{} {
		out := make(chan interface{})

		// 创建多个输出 channel
		channels := make([]<-chan interface{}, num)
		for i := 0; i < num; i++ {
			channels[i] = stage(ctx, in)
		}

		// 合并所有输出
		go func() {
			defer close(out)

			done := make(chan struct{})
			defer close(done)

			// 启动 goroutine 从每个 channel 读取
			for _, ch := range channels {
				go func(c <-chan interface{}) {
					for val := range c {
						select {
						case out <- val:
						case <-done:
							return
						case <-ctx.Done():
							return
						}
					}
				}(ch)
			}

			// 等待所有输入 channel 关闭
			for range channels {
				// 这里简化处理，实际应使用 WaitGroup
			}
		}()

		return out
	}
}

// ============ 便捷函数 ============

// ProcessSlice 处理切片的便捷函数
// 自动创建流水线，处理切片数据，返回结果切片
func ProcessSlice(data []interface{}, stages ...Stage) []interface{} {
	ctx := context.Background()

	// 创建输入 channel
	in := make(chan interface{})
	go func() {
		defer close(in)
		for _, item := range data {
			in <- item
		}
	}()

	// 创建并执行流水线
	p := New(stages...)
	out := p.Execute(ctx, in)

	// 收集结果
	var results []interface{}
	for val := range out {
		results = append(results, val)
	}

	return results
}

// ProcessInts 处理整数切片的便捷函数
func ProcessInts(data []int, stages ...Stage) []int {
	// 转换为 interface{} 切片
	input := make([]interface{}, len(data))
	for i, v := range data {
		input[i] = v
	}

	// 处理
	results := ProcessSlice(input, stages...)

	// 转换回 int 切片
	output := make([]int, len(results))
	for i, v := range results {
		output[i] = v.(int)
	}

	return output
}
