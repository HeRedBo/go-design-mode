package pipeline

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// TestPipeline_Basic 测试基本流水线
func TestPipeline_Basic(t *testing.T) {
	// 创建流水线：平方 -> 过滤偶数
	p := New(
		Square(),
		IsEven(),
	)

	// 输入数据
	in := make(chan interface{})
	go func() {
		defer close(in)
		for i := 1; i <= 5; i++ {
			in <- i
		}
	}()

	// 执行流水线
	ctx := context.Background()
	out := p.Execute(ctx, in)

	// 收集结果
	var results []int
	for val := range out {
		results = append(results, val.(int))
	}

	// 1-5 的平方：1, 4, 9, 16, 25
	// 偶数：4, 16
	expected := []int{4, 16}

	if len(results) != len(expected) {
		t.Errorf("expected %d results, got %d", len(expected), len(results))
	}

	for i, v := range expected {
		if results[i] != v {
			t.Errorf("expected %d at index %d, got %d", v, i, results[i])
		}
	}
}

// TestPipeline_Map 测试 Map 阶段
func TestPipeline_Map(t *testing.T) {
	// 将所有数字乘以 2
	stage := Map(func(val interface{}) interface{} {
		n := val.(int)
		return n * 2
	})

	ctx := context.Background()
	in := make(chan interface{}, 5)
	for i := 1; i <= 5; i++ {
		in <- i
	}
	close(in)

	out := stage(ctx, in)

	var results []int
	for val := range out {
		results = append(results, val.(int))
	}

	expected := []int{2, 4, 6, 8, 10}
	if len(results) != len(expected) {
		t.Errorf("expected %d results, got %d", len(expected), len(results))
	}
}

// TestPipeline_Filter 测试 Filter 阶段
func TestPipeline_Filter(t *testing.T) {
	// 只保留大于 3 的数字
	stage := Filter(func(val interface{}) bool {
		n := val.(int)
		return n > 3
	})

	ctx := context.Background()
	in := make(chan interface{}, 10)
	for i := 1; i <= 10; i++ {
		in <- i
	}
	close(in)

	out := stage(ctx, in)

	var results []int
	for val := range out {
		results = append(results, val.(int))
	}

	expected := []int{4, 5, 6, 7, 8, 9, 10}
	if len(results) != len(expected) {
		t.Errorf("expected %d results, got %d", len(expected), len(results))
	}
}

// TestPipeline_Generator 测试 Generator 阶段
func TestPipeline_Generator(t *testing.T) {
	ctx := context.Background()

	gen := Generator(1, 2, 3, 4, 5)
	out := gen(ctx, nil)

	var results []int
	for val := range out {
		results = append(results, val.(int))
	}

	expected := []int{1, 2, 3, 4, 5}
	if len(results) != len(expected) {
		t.Errorf("expected %d results, got %d", len(expected), len(results))
	}
}

// TestPipeline_Cancellation 测试取消操作
func TestPipeline_Cancellation(t *testing.T) {
	// 创建慢速处理阶段
	slowStage := func(ctx context.Context, in <-chan interface{}) <-chan interface{} {
		out := make(chan interface{})
		go func() {
			defer close(out)
			for val := range in {
				time.Sleep(100 * time.Millisecond) // 模拟慢处理
				select {
				case out <- val:
				case <-ctx.Done():
					return
				}
			}
		}()
		return out
	}

	p := New(slowStage)

	// 输入数据
	in := make(chan interface{})
	go func() {
		defer close(in)
		for i := 0; i < 10; i++ {
			in <- i
		}
	}()

	// 50ms 后取消（由于有缓冲和并发，可能会处理少量数据）
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	out := p.Execute(ctx, in)

	// 收集结果
	count := 0
	for range out {
		count++
	}

	// 由于 50ms 超时，应该只处理了少量数据（0-3 个都是合理的）
	// 这个测试主要验证取消机制生效，不验证精确数量
	if count > 5 {
		t.Errorf("expected at most 5 results due to cancellation, got %d", count)
	}
}

// TestPipeline_ProcessSlice 测试 ProcessSlice 便捷函数
func TestPipeline_ProcessSlice(t *testing.T) {
	data := []interface{}{1, 2, 3, 4, 5}

	results := ProcessSlice(data,
		Square(),
		IsEven(),
	)

	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}

	expected := []interface{}{4, 16}
	for i, v := range expected {
		if results[i] != v {
			t.Errorf("expected %v at index %d, got %v", v, i, results[i])
		}
	}
}

// TestPipeline_ProcessInts 测试 ProcessInts 便捷函数
func TestPipeline_ProcessInts(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	results := ProcessInts(data,
		Map(func(val interface{}) interface{} {
			n := val.(int)
			return n * 10
		}),
		Filter(func(val interface{}) bool {
			n := val.(int)
			return n <= 50
		}),
	)

	expected := []int{10, 20, 30, 40, 50}
	if len(results) != len(expected) {
		t.Errorf("expected %d results, got %d", len(expected), len(results))
	}

	for i, v := range expected {
		if results[i] != v {
			t.Errorf("expected %d at index %d, got %d", v, i, results[i])
		}
	}
}

// TestPipeline_MultipleStages 测试多阶段流水线
func TestPipeline_MultipleStages(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}

	// 复杂流水线：*2 -> +1 -> 过滤奇数
	results := ProcessInts(data,
		Map(func(val interface{}) interface{} {
			return val.(int) * 2
		}),
		Map(func(val interface{}) interface{} {
			return val.(int) + 1
		}),
		Filter(func(val interface{}) bool {
			return val.(int)%2 == 0
		}),
	)

	// 1->2->3(奇), 2->4->5(奇), 3->6->7(奇), 4->8->9(奇), 5->10->11(奇)
	// 没有偶数
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

// TestPipeline_ToString 测试 ToString 阶段
func TestPipeline_ToString(t *testing.T) {
	data := []interface{}{1, 2, 3}

	results := ProcessSlice(data, ToString())

	if len(results) != 3 {
		t.Errorf("expected 3 results, got %d", len(results))
	}

	expected := []interface{}{"1", "2", "3"}
	for i, v := range expected {
		if results[i] != v {
			t.Errorf("expected %v at index %d, got %v", v, i, results[i])
		}
	}
}

// TestPipeline_WithError 测试带错误处理的流水线
func TestPipeline_WithError(t *testing.T) {
	// 会失败的转换函数
	riskyMap := MapWithErr(func(val interface{}) (interface{}, error) {
		n := val.(int)
		if n == 3 {
			return nil, fmt.Errorf("cannot process 3")
		}
		return n * 2, nil
	})

	data := []interface{}{1, 2, 3, 4, 5}
	results := ProcessSlice(data, riskyMap)

	// 应该有 5 个结果（包括错误）
	if len(results) != 5 {
		t.Errorf("expected 5 results, got %d", len(results))
	}

	// 验证错误处理
	for _, val := range results {
		result := val.(Result)
		if result.Value != nil && result.Value.(int) == 6 {
			t.Errorf("value 6 should have error, but got success")
		}
	}
}

// TestPipeline_FilterErrors 测试过滤错误
func TestPipeline_FilterErrors(t *testing.T) {
	riskyMap := MapWithErr(func(val interface{}) (interface{}, error) {
		n := val.(int)
		if n == 3 {
			return nil, fmt.Errorf("error on 3")
		}
		return n * 2, nil
	})

	data := []interface{}{1, 2, 3, 4, 5}
	results := ProcessSlice(data, riskyMap, FilterErrors())

	// 过滤后应该只有 4 个成功结果
	if len(results) != 4 {
		t.Errorf("expected 4 results after filtering errors, got %d", len(results))
	}
}

// TestPipeline_Concurrent 测试并发处理
func TestPipeline_Concurrent(t *testing.T) {
	data := make([]int, 100)
	for i := 0; i < 100; i++ {
		data[i] = i + 1
	}

	start := time.Now()

	// 流水线：平方 -> 过滤大于 1000
	results := ProcessInts(data,
		Square(),
		Filter(func(val interface{}) bool {
			return val.(int) > 1000
		}),
	)

	elapsed := time.Since(start)

	// 验证结果数量
	// 32^2 = 1024，所以从 32 开始都大于 1000
	// 32-100 共 69 个数字
	if len(results) != 69 {
		t.Errorf("expected 69 results, got %d", len(results))
	}

	t.Logf("processed 100 items in %v", elapsed)
}

// TestPipeline_AddStage 测试动态添加阶段
func TestPipeline_AddStage(t *testing.T) {
	p := New(Square())

	// 动态添加阶段
	p.AddStage(IsEven())

	data := []int{1, 2, 3, 4, 5}
	results := ProcessInts(data, p.stages...)

	expected := []int{4, 16}
	if len(results) != len(expected) {
		t.Errorf("expected %d results, got %d", len(expected), len(results))
	}
}

// TestPipeline_Empty 测试空流水线
func TestPipeline_Empty(t *testing.T) {
	p := New() // 没有阶段

	data := []int{1, 2, 3}
	results := ProcessInts(data, p.stages...)

	// 应该原样返回
	if len(results) != 3 {
		t.Errorf("expected 3 results, got %d", len(results))
	}

	expected := []int{1, 2, 3}
	for i, v := range expected {
		if results[i] != v {
			t.Errorf("expected %d at index %d, got %d", v, i, results[i])
		}
	}
}

// TestPipeline_ContextCancel 测试 context 取消传播
func TestPipeline_ContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	// 输入数据
	in := make(chan interface{}, 100)
	go func() {
		for i := 0; i < 100; i++ {
			in <- i
		}
	}()

	// 慢速阶段
	slowStage := func(ctx context.Context, in <-chan interface{}) <-chan interface{} {
		out := make(chan interface{})
		go func() {
			defer close(out)
			for val := range in {
				time.Sleep(10 * time.Millisecond)
				select {
				case out <- val:
				case <-ctx.Done():
					return
				}
			}
		}()
		return out
	}

	p := New(slowStage, slowStage)
	out := p.Execute(ctx, in)

	// 处理一些数据后取消
	count := 0
	for val := range out {
		count++
		if count == 3 {
			cancel() // 取消
			break
		}
		_ = val
	}

	// 验证取消生效
	if count > 5 {
		t.Errorf("expected cancellation to stop processing, got %d results", count)
	}
}

// ExamplePipeline 演示流水线的基本使用
func ExamplePipeline() {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// 创建流水线：平方 -> 过滤偶数
	results := ProcessInts(data,
		Square(),
		IsEven(),
	)

	fmt.Println("Results:", results)
	// Output:
	// Results: [4 16 36 64 100]
}
