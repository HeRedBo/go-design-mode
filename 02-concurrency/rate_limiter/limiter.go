// Package ratelimiter 演示限流器模式 (Rate Limiter)
package ratelimiter

import (
	"context"
	"time"
)

// RateLimiter 限流器
type RateLimiter struct {
	rate       time.Duration // 请求间隔
	tokens     chan struct{} // 令牌桶
	cancel     context.CancelFunc
}

// New 创建限流器
// rate: 允许的请求间隔
func New(rate time.Duration) *RateLimiter {
	ctx, cancel := context.WithCancel(context.Background())
	rl := &RateLimiter{
		rate:   rate,
		tokens: make(chan struct{}, 1),
		cancel: cancel,
	}

	// 启动令牌生成器
	go rl.generateTokens(ctx)

	// 初始令牌
	rl.tokens <- struct{}{}

	return rl
}

// generateTokens 生成令牌
func (rl *RateLimiter) generateTokens(ctx context.Context) {
	ticker := time.NewTicker(rl.rate)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			select {
			case rl.tokens <- struct{}{}:
			default:
				// 令牌桶已满
			}
		}
	}
}

// Wait 等待令牌
func (rl *RateLimiter) Wait(ctx context.Context) error {
	select {
	case <-rl.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Close 关闭限流器
func (rl *RateLimiter) Close() {
	rl.cancel()
}
