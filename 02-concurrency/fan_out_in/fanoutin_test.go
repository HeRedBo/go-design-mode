package fanoutin

import (
	"context"
	"testing"
)

func TestFanOutIn_Basic(t *testing.T) {
	ctx := context.Background()

	// 创建输入
	input := make(chan interface{}, 10)
	for i := 1; i <= 10; i++ {
		input <- i
	}
	close(input)

	// 执行 Fan-Out/Fan-In
	foi := New(3)
	output := foi.Process(ctx, input, func(val interface{}) interface{} {
		n := val.(int)
		return n * 2
	})

	// 收集结果
	var results []int
	for val := range output {
		results = append(results, val.(int))
	}

	if len(results) != 10 {
		t.Errorf("expected 10 results, got %d", len(results))
	}
}

func TestProcessSlice(t *testing.T) {
	ctx := context.Background()
	data := []interface{}{1, 2, 3, 4, 5}

	results := ProcessSlice(ctx, data, 3, func(val interface{}) interface{} {
		return val.(int) * 10
	})

	if len(results) != 5 {
		t.Errorf("expected 5 results, got %d", len(results))
	}
}
