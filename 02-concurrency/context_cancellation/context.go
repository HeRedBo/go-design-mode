// Package contextcancel 演示 Go 语言的 Context 取消模式
//
// Context 模式解决的问题：
// 1. goroutine 生命周期管理：如何优雅地停止 goroutine
// 2. 超时控制：避免操作无限期等待
// 3. 级联取消：一个操作取消时，所有相关操作都应取消
// 4. 请求范围的值传递：在请求链路中传递元数据
//
// 使用场景：
// - HTTP 请求处理（请求取消时自动停止处理）
// - 数据库查询超时控制
// - RPC 调用取消
// - 后台任务管理
// - 优雅关闭服务
//
// Context 的四种创建方式：
// 1. context.Background(): 根 context，永不取消
// 2. context.TODO(): 临时占位，后续需要替换
// 3. context.WithCancel(): 手动取消
// 4. context.WithTimeout(): 超时自动取消
// 5. context.WithDeadline(): 在指定时间取消
// 6. context.WithValue(): 携带键值对
//
// 最佳实践：
// 1. Context 应该是第一个参数，命名为 ctx
// 2. 不要将 Context 放在结构体中，应该显式传递
// 3. 不要传递 nil Context，使用 context.TODO()
// 4. Context 的 Value 只用于请求范围的元数据，不要用于传递可选参数
// 5. 相同的 Context 可以传递给多个 goroutine
// 6. 总是调用 cancel 函数，使用 defer 确保执行
// 7. 检查 ctx.Done() 而不是 ctx.Err()
//
// Context 取消传播：
// - 父 context 取消时，所有子 context 都会收到取消信号
// - 形成 context tree，任何节点取消都会影响其子树
package contextcancel

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ============ 基本使用示例 ============

// ManualCancel 演示手动取消
func ManualCancel() error {
	// 创建可取消的 context
	ctx, cancel := context.WithCancel(context.Background())
	// 确保 cancel 被调用
	defer cancel()

	// 启动 goroutine
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			select {
			case <-ctx.Done():
				// 收到取消信号
				fmt.Println("Goroutine cancelled")
				return
			default:
				// 模拟工作
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()

	// 模拟一些工作后取消
	time.Sleep(50 * time.Millisecond)
	cancel() // 触发取消

	// 等待 goroutine 退出
	<-done
	return nil
}

// TimeoutControl 演示超时控制
func TimeoutControl() error {
	// 创建带超时的 context（2秒超时）
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // 确保资源被释放

	// 模拟长时间操作
	select {
	case <-time.After(3 * time.Second):
		// 超时前未完成
		return fmt.Errorf("operation took too long")
	case <-ctx.Done():
		// 超时或被取消
		return fmt.Errorf("context done: %v", ctx.Err())
	}
}

// Deadline 演示截止时间控制
func Deadline() error {
	// 创建带截止时间的 context
	deadline := time.Now().Add(1 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	// 检查是否已经过期
	if time.Now().After(deadline) {
		return fmt.Errorf("deadline already passed")
	}

	// 模拟工作
	select {
	case <-time.After(2 * time.Second):
		return fmt.Errorf("operation exceeded deadline")
	case <-ctx.Done():
		return fmt.Errorf("context done: %v", ctx.Err())
	}
}

// ============ 级联取消 ============

// CascadingCancel 演示级联取消
// 父 context 取消时，所有子 context 都会收到信号
func CascadingCancel() {
	// 创建父 context
	parentCtx, parentCancel := context.WithCancel(context.Background())
	defer parentCancel()

	// 创建子 context
	childCtx1, childCancel1 := context.WithCancel(parentCtx)
	defer childCancel1()

	childCtx2, childCancel2 := context.WithCancel(parentCtx)
	defer childCancel2()

	// 启动多个 goroutine
	var wg sync.WaitGroup
	wg.Add(3)

	// 父 goroutine
	go func() {
		defer wg.Done()
		<-parentCtx.Done()
		fmt.Println("Parent cancelled")
	}()

	// 子 goroutine 1
	go func() {
		defer wg.Done()
		<-childCtx1.Done()
		fmt.Println("Child 1 cancelled")
	}()

	// 子 goroutine 2
	go func() {
		defer wg.Done()
		<-childCtx2.Done()
		fmt.Println("Child 2 cancelled")
	}()

	// 取消父 context，所有子 context 都会收到信号
	parentCancel()

	// 等待所有 goroutine 退出
	wg.Wait()
}

// ============ Context 传值 ============

// ContextKey 定义 context key 类型（避免冲突）
type ContextKey string

const (
	// RequestIDKey 请求 ID 的 key
	RequestIDKey ContextKey = "request_id"
	// UserIDKey 用户 ID 的 key
	UserIDKey ContextKey = "user_id"
	// TraceIDKey 追踪 ID 的 key
	TraceIDKey ContextKey = "trace_id"
)

// WithRequestID 添加请求 ID 到 context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// GetRequestID 从 context 获取请求 ID
func GetRequestID(ctx context.Context) string {
	if val := ctx.Value(RequestIDKey); val != nil {
		return val.(string)
	}
	return ""
}

// WithUserID 添加用户 ID 到 context
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// GetUserID 从 context 获取用户 ID
func GetUserID(ctx context.Context) string {
	if val := ctx.Value(UserIDKey); val != nil {
		return val.(string)
	}
	return ""
}

// ============ 实际应用场景 ============

// LongRunningTask 长时间运行的任务，支持取消
func LongRunningTask(ctx context.Context) error {
	for i := 0; i < 10; i++ {
		// 每次迭代都检查 context
		select {
		case <-ctx.Done():
			return fmt.Errorf("task cancelled at step %d: %v", i, ctx.Err())
		default:
			// 模拟工作
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("Step %d completed\n", i+1)
		}
	}
	return nil
}

// HTTPHandler 模拟 HTTP 处理器（支持请求取消）
func HTTPHandler(ctx context.Context) error {
	// 模拟数据库查询
	result, err := DatabaseQuery(ctx, "SELECT * FROM users")
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	// 处理结果
	fmt.Printf("Query result: %v\n", result)
	return nil
}

// DatabaseQuery 模拟数据库查询（支持超时和取消）
func DatabaseQuery(ctx context.Context, query string) (string, error) {
	// 创建带超时的 context（1秒）
	queryCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	// 模拟查询
	select {
	case <-time.After(2 * time.Second):
		// 查询超时
		return "", fmt.Errorf("query timeout")
	case <-queryCtx.Done():
		// 查询被取消或超时
		return "", fmt.Errorf("query cancelled: %v", queryCtx.Err())
	case <-time.After(500 * time.Millisecond):
		// 查询成功
		return "result", nil
	}
}

// GracefulShutdown 演示优雅关闭
func GracefulShutdown() {
	// 创建可取消的 context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动多个服务
	var wg sync.WaitGroup
	wg.Add(3)

	// 服务 1
	go func() {
		defer wg.Done()
		RunService(ctx, "Service-1")
	}()

	// 服务 2
	go func() {
		defer wg.Done()
		RunService(ctx, "Service-2")
	}()

	// 服务 3
	go func() {
		defer wg.Done()
		RunService(ctx, "Service-3")
	}()

	// 模拟运行一段时间后关闭
	time.Sleep(200 * time.Millisecond)
	fmt.Println("Initiating graceful shutdown...")
	cancel() // 触发关闭

	// 等待所有服务退出
	wg.Wait()
	fmt.Println("All services stopped")
}

// RunService 运行服务（支持取消）
func RunService(ctx context.Context, name string) {
	fmt.Printf("%s starting...\n", name)

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%s shutting down gracefully\n", name)
			return
		case <-ticker.C:
			// 模拟处理请求
			fmt.Printf("%s processing...\n", name)
		}
	}
}

// ============ 辅助函数 ============

// DoWithTimeout 带超时的操作
func DoWithTimeout(ctx context.Context, timeout time.Duration, fn func() error) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- fn()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

// Race 多个操作中第一个完成的获胜，其他被取消
func Race(ctx context.Context, operations ...func(context.Context) error) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error, len(operations))

	for _, op := range operations {
		go func(fn func(context.Context) error) {
			errCh <- fn(ctx)
		}(op)
	}

	// 等待第一个结果
	return <-errCh
}

// RetryWithBackoff 带退避的重试，支持取消
func RetryWithBackoff(ctx context.Context, maxRetries int, fn func() error) error {
	var lastErr error

	for i := 0; i <= maxRetries; i++ {
		// 检查是否取消
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// 尝试执行
		lastErr = fn()
		if lastErr == nil {
			return nil
		}

		// 计算退避时间
		backoff := time.Duration(i*i) * 100 * time.Millisecond
		if backoff > 5*time.Second {
			backoff = 5 * time.Second
		}

		fmt.Printf("Attempt %d failed: %v, retrying in %v...\n", i+1, lastErr, backoff)

		// 等待重试（可取消）
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
			// 继续重试
		}
	}

	return fmt.Errorf("after %d retries: %w", maxRetries, lastErr)
}
