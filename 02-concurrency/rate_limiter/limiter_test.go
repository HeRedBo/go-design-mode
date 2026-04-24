package ratelimiter

import (
	"context"
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	rl := New(100 * time.Millisecond)
	defer rl.Close()

	ctx := context.Background()

	// 第一次应该立即成功
	err := rl.Wait(ctx)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// 第二次应该等待
	err = rl.Wait(ctx)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
