// Package backoffretry 演示退避重试模式 (Backoff Retry)
package backoffretry

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"
)

// BackOff 退避策略
type BackOff struct {
	InitialInterval     time.Duration // 初始间隔
	MaxInterval         time.Duration // 最大间隔
	MaxRetries          int           // 最大重试次数
	Multiplier          float64       // 倍增系数
	RandomizationFactor float64       // 随机化因子
}

// DefaultBackOff 默认退避策略
func DefaultBackOff() *BackOff {
	return &BackOff{
		InitialInterval:     100 * time.Millisecond,
		MaxInterval:         10 * time.Second,
		MaxRetries:          5,
		Multiplier:          2.0,
		RandomizationFactor: 0.5,
	}
}

// CalculateBackoff 计算下次退避时间
func (b *BackOff) CalculateBackoff(attempt int) time.Duration {
	// 指数增长
	interval := b.InitialInterval * time.Duration(math.Pow(b.Multiplier, float64(attempt)))

	// 限制最大值
	if interval > b.MaxInterval {
		interval = b.MaxInterval
	}

	// 添加随机化（避免惊群效应）
	if b.RandomizationFactor > 0 {
		delta := b.RandomizationFactor * float64(interval)
		minInterval := float64(interval) - delta
		maxInterval := float64(interval) + delta
		randomInterval := minInterval + rand.Float64()*(maxInterval-minInterval)
		interval = time.Duration(randomInterval)
	}

	return interval
}

// Retry 执行重试逻辑
func Retry(ctx context.Context, backoff *BackOff, fn func() error) error {
	var lastErr error

	for attempt := 0; attempt <= backoff.MaxRetries; attempt++ {
		// 检查是否取消
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// 执行函数
		lastErr = fn()
		if lastErr == nil {
			return nil
		}

		// 最后一次失败，不再等待
		if attempt == backoff.MaxRetries {
			break
		}

		// 计算退避时间
		wait := backoff.CalculateBackoff(attempt)

		// 等待重试（可取消）
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(wait):
			// 继续重试
		}
	}

	return fmt.Errorf("after %d retries: %w", backoff.MaxRetries, lastErr)
}

// RetryWithConstant 固定间隔重试
func RetryWithConstant(ctx context.Context, interval time.Duration, maxRetries int, fn func() error) error {
	return Retry(ctx, &BackOff{
		InitialInterval: interval,
		MaxInterval:     interval,
		Multiplier:      1.0,
		MaxRetries:      maxRetries,
	}, fn)
}
