// Package gracefulshutdown 演示 Go 语言的优雅退出模式
//
// 优雅退出解决的问题：
// 1. 等待正在处理的请求完成
// 2. 清理资源（关闭连接、保存状态等）
// 3. 通知相关组件停止工作
// 4. 避免数据丢失或服务中断
//
// 使用场景：
// - HTTP 服务器关闭
// - 后台任务停止
// - 微服务优雅下线
// - 消息队列消费者退出
//
// 核心组件：
// - os/signal: 监听系统信号（SIGINT, SIGTERM）
// - context.Context: 控制生命周期
// - sync.WaitGroup: 等待 goroutine 完成
// - chan: 协调关闭流程
//
// 关闭流程：
// 1. 接收关闭信号
// 2. 停止接收新请求
// 3. 等待正在处理的请求完成
// 4. 清理资源
// 5. 退出程序
package gracefulshutdown

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// ShutdownManager 优雅退出管理器
type ShutdownManager struct {
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	services   []Service
	quitSignal chan os.Signal
	timeout    time.Duration
}

// Service 服务接口
type Service interface {
	Name() string
	Start(ctx context.Context) error
	Stop() error
}

// NewShutdownManager 创建优雅退出管理器
func NewShutdownManager(timeout time.Duration) *ShutdownManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &ShutdownManager{
		ctx:        ctx,
		cancel:     cancel,
		services:   make([]Service, 0),
		quitSignal: make(chan os.Signal, 1),
		timeout:    timeout,
	}
}

// Register 注册服务
func (sm *ShutdownManager) Register(service Service) {
	sm.services = append(sm.services, service)
}

// Start 启动所有服务并监听信号
func (sm *ShutdownManager) Start() error {
	// 监听系统信号
	signal.Notify(sm.quitSignal, syscall.SIGINT, syscall.SIGTERM)

	// 启动所有服务
	for _, svc := range sm.services {
		fmt.Printf("[ShutdownManager] Starting service: %s\n", svc.Name())
		if err := svc.Start(sm.ctx); err != nil {
			return fmt.Errorf("failed to start service %s: %w", svc.Name(), err)
		}
	}

	fmt.Println("[ShutdownManager] All services started, waiting for shutdown signal...")

	// 等待关闭信号
	<-sm.quitSignal
	fmt.Println("\n[ShutdownManager] Received shutdown signal, initiating graceful shutdown...")

	// 执行关闭
	return sm.Shutdown()
}

// Shutdown 执行优雅关闭
func (sm *ShutdownManager) Shutdown() error {
	// 取消 context，通知所有 goroutine 停止工作
	sm.cancel()

	// 创建带超时的 context
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), sm.timeout)
	defer shutdownCancel()

	// 停止所有服务
	done := make(chan error, len(sm.services))
	for _, svc := range sm.services {
		go func(service Service) {
			fmt.Printf("[ShutdownManager] Stopping service: %s\n", service.Name())
			if err := service.Stop(); err != nil {
				done <- fmt.Errorf("failed to stop service %s: %w", service.Name(), err)
			} else {
				fmt.Printf("[ShutdownManager] Service %s stopped gracefully\n", service.Name())
				done <- nil
			}
		}(svc)
	}

	// 等待所有服务停止或超时
	var firstErr error
	for i := 0; i < len(sm.services); i++ {
		select {
		case err := <-done:
			if err != nil && firstErr == nil {
				firstErr = err
			}
		case <-shutdownCtx.Done():
			return fmt.Errorf("shutdown timed out after %v", sm.timeout)
		}
	}

	// 等待所有 goroutine 完成
	sm.wg.Wait()

	fmt.Println("[ShutdownManager] Graceful shutdown completed")
	return firstErr
}

// Context 返回管理器上下文
func (sm *ShutdownManager) Context() context.Context {
	return sm.ctx
}

// AddWorker 添加工作协程
func (sm *ShutdownManager) AddWorker(fn func(ctx context.Context)) {
	sm.wg.Add(1)
	go func() {
		defer sm.wg.Done()
		fn(sm.ctx)
	}()
}

// ============ 示例服务 ============

// HTTPService 模拟 HTTP 服务
type HTTPService struct {
	name string
}

// NewHTTPService 创建 HTTP 服务
func NewHTTPService(name string) *HTTPService {
	return &HTTPService{name: name}
}

// Name 返回服务名称
func (h *HTTPService) Name() string {
	return h.name
}

// Start 启动服务
func (h *HTTPService) Start(ctx context.Context) error {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				fmt.Printf("[%s] Shutting down...\n", h.name)
				return
			case <-ticker.C:
				fmt.Printf("[%s] Processing requests...\n", h.name)
			}
		}
	}()
	return nil
}

// Stop 停止服务
func (h *HTTPService) Stop() error {
	fmt.Printf("[%s] Cleaning up resources...\n", h.name)
	time.Sleep(500 * time.Millisecond) // 模拟清理
	return nil
}

// WorkerService 模拟后台工作服务
type WorkerService struct {
	name string
}

// NewWorkerService 创建工作服务
func NewWorkerService(name string) *WorkerService {
	return &WorkerService{name: name}
}

// Name 返回服务名称
func (w *WorkerService) Name() string {
	return w.name
}

// Start 启动服务
func (w *WorkerService) Start(ctx context.Context) error {
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				fmt.Printf("[%s] Finishing current task...\n", w.name)
				return
			case <-ticker.C:
				fmt.Printf("[%s] Processing tasks...\n", w.name)
			}
		}
	}()
	return nil
}

// Stop 停止服务
func (w *WorkerService) Stop() error {
	fmt.Printf("[%s] Saving state...\n", w.name)
	time.Sleep(300 * time.Millisecond)
	return nil
}

// ============ 便捷函数 ============

// QuickShutdown 快速优雅退出示例
func QuickShutdown() {
	fmt.Println("=== Quick Shutdown Demo ===")

	// 创建管理器（超时 3 秒）
	manager := NewShutdownManager(3 * time.Second)

	// 注册服务
	manager.Register(NewHTTPService("HTTP-Server"))
	manager.Register(NewWorkerService("Background-Worker"))

	// 模拟运行 2 秒后发送关闭信号
	go func() {
		time.Sleep(2 * time.Second)
		// 发送 SIGTERM 信号
		manager.quitSignal <- syscall.SIGTERM
	}()

	// 启动（会阻塞直到关闭完成）
	if err := manager.Start(); err != nil {
		fmt.Printf("Shutdown error: %v\n", err)
	}
}
