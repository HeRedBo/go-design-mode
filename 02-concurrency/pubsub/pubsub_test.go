package pubsub

import (
	"testing"
	"time"
)

func TestPubSub(t *testing.T) {
	ps := New()
	defer ps.Close()

	// 订阅
	ch1 := ps.Subscribe("news")
	ch2 := ps.Subscribe("news")

	// 发布
	ps.Publish("news", "Hello World")

	// 接收
	select {
	case msg := <-ch1:
		if msg != "Hello World" {
			t.Errorf("expected 'Hello World', got %v", msg)
		}
	case <-time.After(1 * time.Second):
		t.Error("timeout waiting for message")
	}

	select {
	case msg := <-ch2:
		if msg != "Hello World" {
			t.Errorf("expected 'Hello World', got %v", msg)
		}
	case <-time.After(1 * time.Second):
		t.Error("timeout waiting for message")
	}
}

func TestUnsubscribe(t *testing.T) {
	ps := New()
	defer ps.Close()

	ch := make(chan interface{}, 10)
	ps.subscribers["test"] = append(ps.subscribers["test"], ch)

	ps.Unsubscribe("test", ch)

	if len(ps.subscribers["test"]) != 0 {
		t.Errorf("expected 0 subscribers, got %d", len(ps.subscribers["test"]))
	}
}
