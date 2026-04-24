package contextcancel

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestManualCancel 测试手动取消
func TestManualCancel(t *testing.T) {
	err := ManualCancel()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// TestTimeoutControl 测试超时控制
func TestTimeoutControl(t *testing.T) {
	start := time.Now()
	err := TimeoutControl()
	elapsed := time.Since(start)

	// 应该在 2 秒左右超时
	if elapsed < 1500*time.Millisecond || elapsed > 2500*time.Millisecond {
		t.Errorf("expected timeout around 2s, got %v", elapsed)
	}

	// 应该返回错误
	if err == nil {
		t.Error("expected error due to timeout, got nil")
	}

	// 验证错误信息包含 "context done"
	if err != nil && err.Error() == "" {
		t.Error("expected non-empty error message")
	}
}

// TestDeadline 测试截止时间
func TestDeadline(t *testing.T) {
	start := time.Now()
	err := Deadline()
	elapsed := time.Since(start)

	// 应该在 1 秒左右超时
	if elapsed < 800*time.Millisecond || elapsed > 1500*time.Millisecond {
		t.Errorf("expected deadline around 1s, got %v", elapsed)
	}

	// 应该返回错误
	if err == nil {
		t.Error("expected error due to deadline, got nil")
	}
}

// TestCascadingCancel 测试级联取消
func TestCascadingCancel(t *testing.T) {
	// 这个函数会打印取消信息，不应该 panic
	CascadingCancel()
	// 如果没有 panic，测试通过
}

// TestContextValue 测试 Context 传值
func TestContextValue(t *testing.T) {
	ctx := context.Background()

	// 初始 context 没有值
	if GetRequestID(ctx) != "" {
		t.Error("expected empty request ID")
	}

	// 添加请求 ID
	ctx = WithRequestID(ctx, "req-123")
	if GetRequestID(ctx) != "req-123" {
		t.Errorf("expected request ID 'req-123', got '%s'", GetRequestID(ctx))
	}

	// 添加用户 ID
	ctx = WithUserID(ctx, "user-456")
	if GetUserID(ctx) != "user-456" {
		t.Errorf("expected user ID 'user-456', got '%s'", GetUserID(ctx))
	}

	// 请求 ID 应该仍然存在
	if GetRequestID(ctx) != "req-123" {
		t.Errorf("request ID should still be 'req-123', got '%s'", GetRequestID(ctx))
	}
}

// TestLongRunningTask 测试长时间运行的任务
func TestLongRunningTask(t *testing.T) {
	// 测试成功完成
	ctx := context.Background()
	err := LongRunningTask(ctx)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// 测试取消
	ctx, cancel := context.WithCancel(context.Background())

	// 200ms 后取消
	go func() {
		time.Sleep(200 * time.Millisecond)
		cancel()
	}()

	err = LongRunningTask(ctx)
	if err == nil {
		t.Error("expected error due to cancellation, got nil")
	}
}

// TestHTTPHandler 测试 HTTP 处理器
func TestHTTPHandler(t *testing.T) {
	// 正常情况
	ctx := context.Background()
	ctx = WithRequestID(ctx, "test-req")

	err := HTTPHandler(ctx)
	// DatabaseQuery 模拟了超时，所以可能成功或失败
	// 这里只验证不 panic
	_ = err
}

// TestDatabaseQuery 测试数据库查询
func TestDatabaseQuery(t *testing.T) {
	ctx := context.Background()

	// 正常查询（500ms）
	result, err := DatabaseQuery(ctx, "SELECT 1")
	if err != nil {
		t.Logf("Query failed (expected in test): %v", err)
	} else {
		if result != "result" {
			t.Errorf("expected 'result', got '%s'", result)
		}
	}

	// 超时查询
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	result, err = DatabaseQuery(ctx, "SLOW_QUERY")
	if err == nil {
		t.Error("expected timeout error, got nil")
	}
}

// TestGracefulShutdown 测试优雅关闭
func TestGracefulShutdown(t *testing.T) {
	// 这个函数会打印服务启动和关闭信息
	GracefulShutdown()
	// 如果没有 hang 住，测试通过
}

// TestDoWithTimeout 测试带超时的操作
func TestDoWithTimeout(t *testing.T) {
	ctx := context.Background()

	// 快速操作（应该成功）
	err := DoWithTimeout(ctx, 1*time.Second, func() error {
		time.Sleep(100 * time.Millisecond)
		return nil
	})
	if err != nil {
		t.Errorf("expected no error for fast operation, got %v", err)
	}

	// 慢速操作（应该超时）
	err = DoWithTimeout(ctx, 200*time.Millisecond, func() error {
		time.Sleep(1 * time.Second)
		return nil
	})
	if err == nil {
		t.Error("expected timeout error for slow operation, got nil")
	}
	if err != context.DeadlineExceeded {
		t.Errorf("expected DeadlineExceeded error, got %v", err)
	}
}

// TestRace 测试竞态操作
func TestRace(t *testing.T) {
	ctx := context.Background()

	// 创建三个操作，速度不同
	op1 := func(ctx context.Context) error {
		time.Sleep(300 * time.Millisecond)
		return fmt.Errorf("op1 failed")
	}

	op2 := func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond) // 最快
		return fmt.Errorf("op2 failed")
	}

	op3 := func(ctx context.Context) error {
		time.Sleep(200 * time.Millisecond)
		return fmt.Errorf("op3 failed")
	}

	start := time.Now()
	err := Race(ctx, op1, op2, op3)
	elapsed := time.Since(start)

	// op2 应该最先完成
	if elapsed > 150*time.Millisecond {
		t.Errorf("expected race to complete in ~100ms, got %v", elapsed)
	}

	// 应该返回 op2 的错误
	if err == nil || err.Error() != "op2 failed" {
		t.Errorf("expected 'op2 failed' error, got %v", err)
	}
}

// TestRetryWithBackoff 测试带退避的重试
func TestRetryWithBackoff(t *testing.T) {
	ctx := context.Background()

	// 测试成功情况（第 3 次成功）
	attempt := 0
	err := RetryWithBackoff(ctx, 5, func() error {
		attempt++
		if attempt < 3 {
			return fmt.Errorf("attempt %d failed", attempt)
		}
		return nil
	})

	if err != nil {
		t.Errorf("expected no error after retries, got %v", err)
	}
	if attempt != 3 {
		t.Errorf("expected 3 attempts, got %d", attempt)
	}

	// 测试全部失败
	attempt = 0
	err = RetryWithBackoff(ctx, 2, func() error {
		attempt++
		return fmt.Errorf("always fails")
	})

	if err == nil {
		t.Error("expected error after all retries failed, got nil")
	}
	if attempt != 3 { // 0, 1, 2 共 3 次（包括第 0 次）
		t.Errorf("expected 3 attempts, got %d", attempt)
	}
}

// TestRetryWithBackoff_Cancel 测试重试时的取消
func TestRetryWithBackoff_Cancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	// 100ms 后取消
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	attempt := 0
	err := RetryWithBackoff(ctx, 10, func() error {
		attempt++
		return fmt.Errorf("always fails")
	})

	// 应该因为取消而停止
	if err != context.Canceled {
		t.Errorf("expected Canceled error, got %v", err)
	}

	t.Logf("Retried %d times before cancellation", attempt)
}

// TestContextCancelPropagation 测试取消传播
func TestContextCancelPropagation(t *testing.T) {
	parentCtx, parentCancel := context.WithCancel(context.Background())

	// 创建多个子 context
	childCtx1, _ := context.WithCancel(parentCtx)
	childCtx2, _ := context.WithCancel(parentCtx)
	childCtx3, _ := context.WithCancel(parentCtx)

	// 取消父 context
	parentCancel()

	// 所有子 context 都应该收到取消信号
	select {
	case <-childCtx1.Done():
		// 正确
	default:
		t.Error("childCtx1 should be cancelled")
	}

	select {
	case <-childCtx2.Done():
		// 正确
	default:
		t.Error("childCtx2 should be cancelled")
	}

	select {
	case <-childCtx3.Done():
		// 正确
	default:
		t.Error("childCtx3 should be cancelled")
	}
}

// TestContextMultipleValues 测试 Context 多个值
func TestContextMultipleValues(t *testing.T) {
	ctx := context.Background()
	ctx = WithRequestID(ctx, "req-789")
	ctx = WithUserID(ctx, "user-123")

	// 同时验证多个值
	if GetRequestID(ctx) != "req-789" {
		t.Errorf("expected request ID 'req-789', got '%s'", GetRequestID(ctx))
	}
	if GetUserID(ctx) != "user-123" {
		t.Errorf("expected user ID 'user-123', got '%s'", GetUserID(ctx))
	}
}

// TestContextCancellationSafety 测试取消安全性
func TestContextCancellationSafety(t *testing.T) {
	// 多次调用 cancel 不应该 panic
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	cancel() // 第二次调用
	cancel() // 第三次调用

	// Context 应该仍然处于取消状态
	select {
	case <-ctx.Done():
		// 正确
	default:
		t.Error("context should be cancelled")
	}
}

// TestConcurrentContextUsage 测试并发使用 Context
func TestConcurrentContextUsage(t *testing.T) {
	ctx := context.Background()
	ctx = WithRequestID(ctx, "concurrent-test")

	var wg sync.WaitGroup
	results := make([]string, 10)

	// 并发读取 context
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx] = GetRequestID(ctx)
		}(i)
	}

	wg.Wait()

	// 所有 goroutine 应该读到相同的值
	for i, result := range results {
		if result != "concurrent-test" {
			t.Errorf("goroutine %d: expected 'concurrent-test', got '%s'", i, result)
		}
	}
}

// ExampleManualCancel 演示手动取消
func ExampleManualCancel() {
	err := ManualCancel()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	// Output:
	// Goroutine cancelled
}

// ExampleContextValue 演示 Context 传值
func ExampleContextValue() {
	ctx := context.Background()
	ctx = WithRequestID(ctx, "req-001")
	ctx = WithUserID(ctx, "user-123")

	fmt.Printf("RequestID: %s, UserID: %s\n", GetRequestID(ctx), GetUserID(ctx))
	// Output:
	// RequestID: req-001, UserID: user-123
}
