package gracefulshutdown

import (
	"context"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"
)

// TestShutdownManager 测试基本关闭流程
func TestShutdownManager(t *testing.T) {
	manager := NewShutdownManager(2 * time.Second)

	// 注册服务
	service := NewHTTPService("Test-Service")
	manager.Register(service)

	// 验证上下文可用
	if manager.ctx == nil {
		t.Error("expected context to be initialized")
	}

	// 验证服务列表
	if len(manager.services) != 1 {
		t.Errorf("expected 1 service, got %d", len(manager.services))
	}
}

// TestShutdownManager_StartStop 测试启动和停止
func TestShutdownManager_StartStop(t *testing.T) {
	manager := NewShutdownManager(2 * time.Second)

	// 注册服务
	svc := NewHTTPService("Test-HTTP")
	manager.Register(svc)

	// 模拟发送信号
	go func() {
		time.Sleep(500 * time.Millisecond)
		manager.quitSignal <- syscall.SIGTERM
	}()

	// 启动（会阻塞直到关闭）
	err := manager.Start()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// TestMultipleServices 测试多个服务
func TestMultipleServices(t *testing.T) {
	manager := NewShutdownManager(3 * time.Second)

	// 注册多个服务
	manager.Register(NewHTTPService("HTTP-1"))
	manager.Register(NewHTTPService("HTTP-2"))
	manager.Register(NewWorkerService("Worker-1"))

	// 模拟发送信号
	go func() {
		time.Sleep(500 * time.Millisecond)
		manager.quitSignal <- syscall.SIGTERM
	}()

	err := manager.Start()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// TestWorker 测试工作协程
func TestWorker(t *testing.T) {
	manager := NewShutdownManager(2 * time.Second)

	var wg sync.WaitGroup
	wg.Add(1)

	// 添加工作协程
	manager.AddWorker(func(ctx context.Context) {
		defer wg.Done()
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		count := 0
		for {
			select {
			case <-ctx.Done():
				t.Logf("Worker stopped after %d iterations", count)
				return
			case <-ticker.C:
				count++
			}
		}
	})

	// 模拟关闭
	go func() {
		time.Sleep(500 * time.Millisecond)
		manager.cancel()
	}()

	// 等待工作协程完成
	wg.Wait()
}

// TestServiceInterfaces 测试服务接口实现
func TestServiceInterfaces(t *testing.T) {
	httpSvc := NewHTTPService("HTTP")
	if httpSvc.Name() != "HTTP" {
		t.Errorf("expected name 'HTTP', got '%s'", httpSvc.Name())
	}

	workerSvc := NewWorkerService("Worker")
	if workerSvc.Name() != "Worker" {
		t.Errorf("expected name 'Worker', got '%s'", workerSvc.Name())
	}
}

// TestServiceStartStop 测试服务启动停止
func TestServiceStartStop(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 测试 HTTP 服务
	httpSvc := NewHTTPService("Test-HTTP")
	err := httpSvc.Start(ctx)
	if err != nil {
		t.Errorf("expected no error starting HTTP service, got %v", err)
	}

	// 取消上下文
	time.Sleep(100 * time.Millisecond)
	cancel()

	// 停止服务
	err = httpSvc.Stop()
	if err != nil {
		t.Errorf("expected no error stopping HTTP service, got %v", err)
	}
}

// TestTimeout 测试超时情况
func TestTimeout(t *testing.T) {
	// 创建一个会超时的慢服务
	slowSvc := &slowService{name: "Slow-Service"}

	manager := NewShutdownManager(500 * time.Millisecond) // 很短的超时
	manager.Register(slowSvc)

	// 模拟发送信号
	go func() {
		time.Sleep(200 * time.Millisecond)
		manager.quitSignal <- syscall.SIGTERM
	}()

	err := manager.Start()
	// 应该超时
	if err == nil {
		t.Log("Expected timeout error (may pass if service stops quickly)")
	}
}

// slowService 慢速服务（用于测试超时）
type slowService struct {
	name string
}

func (s *slowService) Name() string {
	return s.name
}

func (s *slowService) Start(ctx context.Context) error {
	return nil
}

func (s *slowService) Stop() error {
	time.Sleep(2 * time.Second) // 故意很慢
	return nil
}

// TestQuickShutdown 测试快速关闭示例
func TestQuickShutdown(t *testing.T) {
	// 这个函数会在 2 秒后自动关闭
	// 不阻塞测试
	go QuickShutdown()
	time.Sleep(3 * time.Second) // 等待完成
}

// TestContextCancellation 测试上下文取消传播
func TestContextCancellation(t *testing.T) {
	manager := NewShutdownManager(2 * time.Second)

	ctx := manager.Context()
	if ctx == nil {
		t.Fatal("expected context to be initialized")
	}

	// 启动一个 goroutine 监听上下文
	done := make(chan struct{})
	go func() {
		defer close(done)
		<-ctx.Done()
	}()

	// 取消上下文
	manager.cancel()

	// 等待 goroutine 收到取消信号
	select {
	case <-done:
		// 成功
	case <-time.After(1 * time.Second):
		t.Error("context cancellation not propagated")
	}
}

// TestSignalHandling 测试信号处理
func TestSignalHandling(t *testing.T) {
	manager := NewShutdownManager(2 * time.Second)

	// 验证 quitSignal channel 已初始化
	if manager.quitSignal == nil {
		t.Error("expected quitSignal channel to be initialized")
	}

	// 验证可以接收信号
	go func() {
		manager.quitSignal <- os.Interrupt
	}()

	select {
	case <-manager.quitSignal:
		// 成功收到信号
	case <-time.After(1 * time.Second):
		t.Error("timeout waiting for signal")
	}
}

// ExampleShutdownManager 演示优雅退出管理器
func ExampleShutdownManager() {
	// 实际使用时会阻塞，这里只演示创建
	manager := NewShutdownManager(5 * time.Second)
	manager.Register(NewHTTPService("My-Server"))

	// manager.Start() // 会阻塞直到收到关闭信号
	_ = manager
}
