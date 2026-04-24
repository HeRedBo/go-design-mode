// Package pubsub 演示发布订阅模式 (Pub/Sub)
package pubsub

import (
	"sync"
)

// PubSub 发布订阅系统
type PubSub struct {
	mu         sync.RWMutex
	subscribers map[string][]chan interface{}
}

// New 创建 PubSub
func New() *PubSub {
	return &PubSub{
		subscribers: make(map[string][]chan interface{}),
	}
}

// Subscribe 订阅主题
func (ps *PubSub) Subscribe(topic string) <-chan interface{} {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan interface{}, 10)
	ps.subscribers[topic] = append(ps.subscribers[topic], ch)
	return ch
}

// Publish 发布消息到主题
func (ps *PubSub) Publish(topic string, msg interface{}) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, ch := range ps.subscribers[topic] {
		select {
		case ch <- msg:
		default:
			// 订阅者处理慢，跳过
		}
	}
}

// Unsubscribe 取消订阅
func (ps *PubSub) Unsubscribe(topic string, ch chan interface{}) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	subs := ps.subscribers[topic]
	for i, sub := range subs {
		if sub == ch {
			ps.subscribers[topic] = append(subs[:i], subs[i+1:]...)
			close(ch)
			break
		}
	}
}

// Close 关闭所有订阅
func (ps *PubSub) Close() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for topic, subs := range ps.subscribers {
		for _, ch := range subs {
			close(ch)
		}
		delete(ps.subscribers, topic)
	}
}
