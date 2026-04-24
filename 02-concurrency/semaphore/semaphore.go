// Package semaphore 演示信号量模式 (Semaphore)
// 使用带缓冲的 channel 实现信号量
package semaphore

import "context"

// Semaphore 信号量
type Semaphore struct {
	sem chan struct{}
}

// New 创建信号量
// size: 并发数量限制
func New(size int) *Semaphore {
	return &Semaphore{
		sem: make(chan struct{}, size),
	}
}

// Acquire 获取信号量
func (s *Semaphore) Acquire(ctx context.Context) error {
	select {
	case s.sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Release 释放信号量
func (s *Semaphore) Release() {
	<-s.sem
}

// TryAcquire 尝试获取信号量（非阻塞）
func (s *Semaphore) TryAcquire() bool {
	select {
	case s.sem <- struct{}{}:
		return true
	default:
		return false
	}
}

// Available 返回可用信号量数量
func (s *Semaphore) Available() int {
	return cap(s.sem) - len(s.sem)
}
