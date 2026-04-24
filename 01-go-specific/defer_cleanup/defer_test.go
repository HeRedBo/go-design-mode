package deferpattern

import (
	"fmt"
	"sync"
	"testing"
)

// TestSafeCounter 测试线程安全计数器
func TestSafeCounter(t *testing.T) {
	counter := &SafeCounter{}
	var wg sync.WaitGroup

	// 并发增加计数
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}

	// 并发减少计数
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Decrement()
		}()
	}

	wg.Wait()

	// 最终应该是 100 - 50 = 50
	expected := 50
	if counter.GetCount() != expected {
		t.Errorf("expected count %d, got %d", expected, counter.GetCount())
	}
}

// TestMultiResourceCleanup 测试多资源清理顺序
func TestMultiResourceCleanup(t *testing.T) {
	err := MultiResourceCleanup()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	// 输出应该显示资源按 3->2->1 顺序释放
}

// TestRecoverFromPanic 测试 panic 恢复
func TestRecoverFromPanic(t *testing.T) {
	err := RecoverFromPanic()
	if err == nil {
		t.Error("expected error from recovered panic, got nil")
	}

	expected := "recovered from panic: something went wrong"
	if err.Error() != expected {
		t.Errorf("expected '%s', got '%s'", expected, err.Error())
	}
}

// TestSafeFunction 测试安全执行函数
func TestSafeFunction(t *testing.T) {
	// 测试会 panic 的函数
	panicked, value := SafeFunction(func() {
		panic("test panic")
	})

	if !panicked {
		t.Error("expected panicked to be true")
	}
	if value != "test panic" {
		t.Errorf("expected panic value 'test panic', got %v", value)
	}

	// 测试正常的函数
	panicked, value = SafeFunction(func() {
		fmt.Println("normal execution")
	})

	if panicked {
		t.Error("expected panicked to be false")
	}
	if value != nil {
		t.Errorf("expected nil panic value, got %v", value)
	}
}

// TestDeferParameterEvaluation 测试 defer 参数求值时机
func TestDeferParameterEvaluation(t *testing.T) {
	// 这个测试主要验证不 panic
	DeferParameterEvaluation()
}

// TestDeferWithPointer 测试闭包方式获取最新值
func TestDeferWithPointer(t *testing.T) {
	// 这个测试主要验证不 panic
	DeferWithPointer()
}

// TestDeferInLoop 测试循环中的 defer
func TestDeferInLoop(t *testing.T) {
	// 使用不存在的文件，测试错误处理
	files := []string{"/tmp/nonexistent1.txt", "/tmp/nonexistent2.txt"}

	// 应该返回错误（文件不存在）
	err := DeferInLoop(files)
	if err == nil {
		t.Error("expected error for non-existent files, got nil")
	}
}

// TestHTTPResponseCleanup 测试 HTTP 响应清理
func TestHTTPResponseCleanup(t *testing.T) {
	err := HTTPResponseCleanup()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// TestDatabaseQuery 测试数据库查询清理
func TestDatabaseQuery(t *testing.T) {
	err := DatabaseQuery()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// TestDeferLIFO 测试 defer 的 LIFO 顺序
func TestDeferLIFO(t *testing.T) {
	var order []int

	func() {
		defer func() { order = append(order, 1) }()
		defer func() { order = append(order, 2) }()
		defer func() { order = append(order, 3) }()
	}()

	expected := []int{3, 2, 1}
	if len(order) != len(expected) {
		t.Errorf("expected order length %d, got %d", len(expected), len(order))
	}

	for i, v := range expected {
		if order[i] != v {
			t.Errorf("expected order[%d] = %d, got %d", i, v, order[i])
		}
	}
}

// TestDeferPanicRecovery 测试 defer 在 panic 时仍会执行
func TestDeferPanicRecovery(t *testing.T) {
	executed := false

	func() {
		defer func() {
			executed = true
			recover() // 恢复 panic
		}()

		panic("test")
	}()

	if !executed {
		t.Error("expected defer to execute even during panic")
	}
}

// TestConcurrentSafeCounter 测试计数器的并发安全性
func TestConcurrentSafeCounter(t *testing.T) {
	counter := &SafeCounter{}
	var wg sync.WaitGroup

	// 大量并发操作
	for i := 0; i < 1000; i++ {
		wg.Add(3)

		go func() {
			defer wg.Done()
			counter.Increment()
		}()

		go func() {
			defer wg.Done()
			counter.Decrement()
		}()

		go func() {
			defer wg.Done()
			_ = counter.GetCount()
		}()
	}

	wg.Wait()

	// 应该最终是 0（1000 次增加 - 1000 次减少）
	if counter.GetCount() != 0 {
		t.Errorf("expected count 0, got %d", counter.GetCount())
	}
}

// ExampleSafeCounter 演示 SafeCounter 的使用
func ExampleSafeCounter() {
	counter := &SafeCounter{}
	counter.Increment()
	counter.Increment()
	counter.Decrement()

	fmt.Println("Count:", counter.GetCount())
	// Output:
	// Count: 1
}

// ExampleRecoverFromPanic 演示 panic 恢复
func ExampleRecoverFromPanic() {
	err := RecoverFromPanic()
	if err != nil {
		fmt.Println("Recovered:", err)
	}
	// Output:
	// Recovered: recovered from panic: something went wrong
}
