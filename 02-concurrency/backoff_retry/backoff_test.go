package backoffretry

import (
	"context"
	"testing"
	"time"
)

func TestBackoffRetry_Success(t *testing.T) {
	ctx := context.Background()
	attempt := 0

	err := Retry(ctx, DefaultBackOff(), func() error {
		attempt++
		if attempt < 3 {
			return nil
		}
		return nil
	})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestBackoffRetry_Failure(t *testing.T) {
	ctx := context.Background()

	err := Retry(ctx, &BackOff{
		InitialInterval: 10 * time.Millisecond,
		MaxInterval:     50 * time.Millisecond,
		MaxRetries:      2,
		Multiplier:      2.0,
	}, func() error {
		return context.DeadlineExceeded
	})

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestBackoffRetry_Cancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	err := Retry(ctx, DefaultBackOff(), func() error {
		return context.DeadlineExceeded
	})

	if err != context.Canceled {
		t.Errorf("expected Canceled error, got %v", err)
	}
}

func TestCalculateBackoff(t *testing.T) {
	bo := DefaultBackOff()

	// 测试增长
	d1 := bo.CalculateBackoff(0)
	d2 := bo.CalculateBackoff(1)

	if d2 <= d1 {
		t.Errorf("expected d2 > d1, got %v <= %v", d2, d1)
	}
}
