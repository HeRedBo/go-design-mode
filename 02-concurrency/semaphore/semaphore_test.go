package semaphore

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestSemaphore(t *testing.T) {
	sem := New(3)

	// 获取 3 个信号量
	ctx := context.Background()
	for i := 0; i < 3; i++ {
		if err := sem.Acquire(ctx); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	}

	// 应该没有可用信号量
	if sem.Available() != 0 {
		t.Errorf("expected 0 available, got %d", sem.Available())
	}

	// 释放一个
	sem.Release()
	if sem.Available() != 1 {
		t.Errorf("expected 1 available, got %d", sem.Available())
	}
}

func TestSemaphore_Concurrent(t *testing.T) {
	sem := New(5)
	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			if err := sem.Acquire(ctx); err != nil {
				return
			}
			defer sem.Release()

			time.Sleep(10 * time.Millisecond)
		}()
	}

	wg.Wait()
}

func TestTryAcquire(t *testing.T) {
	sem := New(2)

	// 前两个应该成功
	if !sem.TryAcquire() {
		t.Error("expected TryAcquire to succeed")
	}
	if !sem.TryAcquire() {
		t.Error("expected TryAcquire to succeed")
	}

	// 第三个应该失败
	if sem.TryAcquire() {
		t.Error("expected TryAcquire to fail")
	}
}
