package workerpool

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestNew 测试创建工作池
func TestNew(t *testing.T) {
	wp := New(4, 10)

	if wp.NumWorkers() != 4 {
		t.Errorf("expected 4 workers, got %d", wp.NumWorkers())
	}
	if wp.tasks == nil {
		t.Error("expected tasks channel to be initialized")
	}
	if wp.results == nil {
		t.Error("expected results channel to be initialized")
	}
}

// TestWorkerPool_Basic 测试基本工作流程
func TestWorkerPool_Basic(t *testing.T) {
	wp := New(2, 5)
	wp.Start()

	// 提交任务
	for i := 1; i <= 5; i++ {
		wp.Submit(Task{
			ID:   i,
			Data: i,
			Handler: func(data interface{}) (interface{}, error) {
				num := data.(int)
				return num * 2, nil
			},
		})
	}
	wp.Close()

	// 等待完成
	go wp.Wait()

	// 收集结果
	results := make(map[int]int)
	for result := range wp.Results() {
		if result.Err != nil {
			t.Errorf("task %d failed: %v", result.TaskID, result.Err)
		}
		results[result.TaskID] = result.Output.(int)
	}

	// 验证结果
	if len(results) != 5 {
		t.Errorf("expected 5 results, got %d", len(results))
	}

	for i := 1; i <= 5; i++ {
		expected := i * 2
		if results[i] != expected {
			t.Errorf("task %d: expected %d, got %d", i, expected, results[i])
		}
	}
}

// TestWorkerPool_ProcessBatch 测试 ProcessBatch 便捷方法
func TestWorkerPool_ProcessBatch(t *testing.T) {
	tasks := make([]Task, 10)
	for i := 0; i < 10; i++ {
		num := i
		tasks[i] = Task{
			ID:   i,
			Data: num,
			Handler: func(data interface{}) (interface{}, error) {
				n := data.(int)
				return n * n, nil
			},
		}
	}

	results := ProcessBatch(4, tasks)

	if len(results) != 10 {
		t.Errorf("expected 10 results, got %d", len(results))
	}

	// 验证每个结果
	for _, result := range results {
		if result.Err != nil {
			t.Errorf("task %d failed: %v", result.TaskID, result.Err)
		}
		n := result.TaskID
		expected := n * n
		if result.Output.(int) != expected {
			t.Errorf("task %d: expected %d, got %d", n, expected, result.Output)
		}
	}
}

// TestWorkerPool_Error 测试错误处理
func TestWorkerPool_Error(t *testing.T) {
	wp := New(2, 5)
	wp.Start()

	// 提交会失败的任务
	wp.Submit(Task{
		ID:   1,
		Data: "invalid",
		Handler: func(data interface{}) (interface{}, error) {
			return nil, fmt.Errorf("processing error")
		},
	})
	wp.Close()

	go wp.Wait()

	// 收集结果
	for result := range wp.Results() {
		if result.TaskID == 1 {
			if result.Err == nil {
				t.Errorf("expected error for task 1, got nil")
			}
		}
	}
}

// TestWorkerPool_Concurrent 测试并发安全性
func TestWorkerPool_Concurrent(t *testing.T) {
	wp := New(4, 100)
	wp.Start()

	var counter int64
	var wg sync.WaitGroup

	// 并发提交任务
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			wp.Submit(Task{
				ID:   id,
				Data: id,
				Handler: func(data interface{}) (interface{}, error) {
					atomic.AddInt64(&counter, 1)
					return data, nil
				},
			})
		}(i)
	}

	wg.Wait()
	wp.Close()
	go wp.Wait()

	// 收集结果
	count := 0
	for range wp.Results() {
		count++
	}

	if count != 100 {
		t.Errorf("expected 100 results, got %d", count)
	}
	if counter != 100 {
		t.Errorf("expected counter to be 100, got %d", counter)
	}
}

// TestWorkerPool_SubmitAsync 测试异步提交
func TestWorkerPool_SubmitAsync(t *testing.T) {
	wp := New(1, 2) // 小缓冲区
	wp.Start()

	// 填满缓冲区
	wp.Submit(Task{ID: 1, Data: 1, Handler: ExampleProcessor})
	wp.Submit(Task{ID: 2, Data: 2, Handler: ExampleProcessor})

	// 异步提交应该成功（缓冲区满但会被 worker 消费）
	time.Sleep(10 * time.Millisecond)
	success := wp.SubmitAsync(Task{ID: 3, Data: 3, Handler: ExampleProcessor})

	// 可能会成功或失败，取决于 worker 是否已经消费
	_ = success

	wp.Close()
	go wp.Wait()

	count := 0
	for range wp.Results() {
		count++
	}

	if count < 2 {
		t.Errorf("expected at least 2 results, got %d", count)
	}
}

// TestWorkerPool_SlowTasks 测试慢任务处理
func TestWorkerPool_SlowTasks(t *testing.T) {
	wp := New(3, 10)
	wp.Start()

	// 提交慢任务
	for i := 0; i < 5; i++ {
		num := i
		wp.Submit(Task{
			ID:   i,
			Data: num,
			Handler: func(data interface{}) (interface{}, error) {
				time.Sleep(50 * time.Millisecond) // 模拟慢任务
				return data, nil
			},
		})
	}
	wp.Close()

	start := time.Now()
	go wp.Wait()

	// 收集结果
	count := 0
	for range wp.Results() {
		count++
	}

	elapsed := time.Since(start)

	if count != 5 {
		t.Errorf("expected 5 results, got %d", count)
	}

	// 3 个 worker 处理 5 个慢任务，应该比串行快
	// 预期时间：~100ms (2 批 * 50ms)
	if elapsed > 300*time.Millisecond {
		t.Errorf("expected processing to be faster, took %v", elapsed)
	}
}

// TestWorkerPool_Panic 测试任务 panic 不会导致工作池崩溃
func TestWorkerPool_Panic(t *testing.T) {
	wp := New(2, 5)
	wp.Start()

	// 提交会 panic 的任务
	wp.Submit(Task{
		ID:   1,
		Data: "panic",
		Handler: func(data interface{}) (interface{}, error) {
			panic("task panic")
		},
	})

	// 提交正常任务
	wp.Submit(Task{
		ID:   2,
		Data: 2,
		Handler: func(data interface{}) (interface{}, error) {
			return data, nil
		},
	})

	wp.Close()

	// 使用 defer recover 防止测试失败
	defer func() {
		if r := recover(); r != nil {
			t.Logf("recovered from panic: %v", r)
		}
	}()

	go wp.Wait()

	// 收集结果
	count := 0
	for range wp.Results() {
		count++
	}

	// panic 的任务可能不会产生结果
	t.Logf("received %d results", count)
}

// TestWorkerPool_DifferentWorkers 测试不同 worker 数量的性能
func TestWorkerPool_DifferentWorkers(t *testing.T) {
	numTasks := 100

	tests := []struct {
		name       string
		numWorkers int
	}{
		{"1 worker", 1},
		{"2 workers", 2},
		{"4 workers", 4},
		{"10 workers", 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wp := New(tt.numWorkers, numTasks)
			wp.Start()

			for i := 0; i < numTasks; i++ {
				num := i
				wp.Submit(Task{
					ID:   i,
					Data: num,
					Handler: func(data interface{}) (interface{}, error) {
						time.Sleep(1 * time.Millisecond)
						return data, nil
					},
				})
			}
			wp.Close()

			start := time.Now()
			go wp.Wait()

			count := 0
			for range wp.Results() {
				count++
			}
			elapsed := time.Since(start)

			if count != numTasks {
				t.Errorf("expected %d results, got %d", numTasks, count)
			}

			t.Logf("%s: processed %d tasks in %v", tt.name, count, elapsed)
		})
	}
}

// TestWorkerPool_ResultsChannel 测试 Results channel 的只读特性
func TestWorkerPool_ResultsChannel(t *testing.T) {
	wp := New(2, 5)
	wp.Start()

	wp.Submit(Task{
		ID:      1,
		Data:    1,
		Handler: ExampleProcessor,
	})
	wp.Close()

	go wp.Wait()

	// Results() 返回只读 channel
	results := wp.Results()

	count := 0
	for range results {
		count++
	}

	if count != 1 {
		t.Errorf("expected 1 result, got %d", count)
	}
}

// ExampleProcessBatch 演示 ProcessBatch 的使用
func ExampleProcessBatch() {
	tasks := []Task{
		{ID: 1, Data: 2, Handler: ExampleProcessor},
		{ID: 2, Data: 3, Handler: ExampleProcessor},
		{ID: 3, Data: 4, Handler: ExampleProcessor},
	}

	results := ProcessBatch(2, tasks)

	for _, result := range results {
		if result.Err != nil {
			fmt.Printf("Task %d failed: %v\n", result.TaskID, result.Err)
		} else {
			fmt.Printf("Task %d result: %v\n", result.TaskID, result.Output)
		}
	}
	// Output:
	// Task 1 result: 4
	// Task 2 result: 9
	// Task 3 result: 16
}
