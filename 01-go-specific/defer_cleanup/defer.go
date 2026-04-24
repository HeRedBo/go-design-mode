// Package deferpattern 演示 Go 语言的 Defer 资源清理模式
//
// Defer 是 Go 独有的特性，用于延迟执行语句
// 特点：
// 1. 延迟到函数返回前执行
// 2. 多个 defer 按 LIFO（后进先出）顺序执行
// 3. 即使函数 panic，defer 仍会执行
// 4. 常用于资源清理：关闭文件、释放锁、关闭连接等
//
// 最佳实践：
// 1. defer 紧跟资源获取之后
// 2. 使用匿名函数处理复杂清理逻辑
// 3. defer 可以用于错误恢复（配合 recover）
// 4. 注意 defer 的参数在注册时就已求值
// 5. 避免在循环中使用 defer（会导致资源延迟释放）
package deferpattern

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// ============ 基础示例 ============

// FileHandling 演示文件操作中的 defer 清理
func FileHandling(filename string) error {
	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	// 立即注册清理，确保文件被关闭
	defer file.Close()

	// 读取文件内容
	data := make([]byte, 1024)
	n, err := file.Read(data)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read file: %w", err)
	}

	fmt.Printf("Read %d bytes\n", n)
	return nil
}

// ============ Mutex 锁清理 ============

// SafeCounter 线程安全的计数器
type SafeCounter struct {
	mu    sync.Mutex
	count int
}

// Increment 增加计数（使用 defer 解锁）
func (c *SafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock() // 确保锁被释放，即使发生 panic

	c.count++
}

// Decrement 减少计数
func (c *SafeCounter) Decrement() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.count--
}

// GetCount 获取计数
func (c *SafeCounter) GetCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.count
}

// ============ 多资源清理 ============

// MultiResourceCleanup 演示多个资源的清理
// defer 按 LIFO 顺序执行（后进先出）
func MultiResourceCleanup() error {
	// 模拟多个资源
	fmt.Println("Acquiring resource 1")
	// resource1 := acquireResource1()
	defer func() {
		fmt.Println("Releasing resource 1")
		// resource1.Release()
	}()

	fmt.Println("Acquiring resource 2")
	// resource2 := acquireResource2()
	defer func() {
		fmt.Println("Releasing resource 2")
		// resource2.Release()
	}()

	fmt.Println("Acquiring resource 3")
	// resource3 := acquireResource3()
	defer func() {
		fmt.Println("Releasing resource 3")
		// resource3.Release()
	}()

	// 执行实际操作
	fmt.Println("Doing work with all resources")

	// 返回时，defer 按 3 -> 2 -> 1 的顺序执行
	return nil
}

// ============ Panic 恢复 ============

// RecoverFromPanic 演示使用 defer 从 panic 恢复
func RecoverFromPanic() (err error) {
	// 命名返回值 + defer 捕获 panic
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
		}
	}()

	// 模拟可能 panic 的代码
	panic("something went wrong")

	// 这行不会执行
	fmt.Println("This will not be printed")
	return nil
}

// SafeFunction 安全执行函数（包装可能 panic 的函数）
func SafeFunction(fn func()) (panicked bool, panicValue interface{}) {
	defer func() {
		if r := recover(); r != nil {
			panicked = true
			panicValue = r
		}
	}()

	fn()
	return false, nil
}

// ============ 注意事项 ============

// DeferParameterEvaluation 演示 defer 参数的求值时机
// 重要：defer 的参数在注册时就已求值，不是在执行时！
func DeferParameterEvaluation() {
	a := 1
	defer fmt.Println("Deferred value:", a) // 输出 1，不是 2

	a = 2 // 这个修改不会影响 defer 的输出
	fmt.Println("Current value:", a)
}

// DeferWithPointer 使用指针或闭包获取最新值
func DeferWithPointer() {
	a := 1
	defer func() {
		fmt.Println("Deferred value with closure:", a) // 输出 2
	}()

	a = 2
}

// DeferInLoop 循环中的 defer 陷阱
// 警告：不要在循环中使用 defer，会导致资源延迟释放
func DeferInLoop(filenames []string) error {
	// 错误示例（不推荐）：
	// for _, filename := range filenames {
	//     file, err := os.Open(filename)
	//     if err != nil {
	//         return err
	//     }
	//     defer file.Close() // 文件要等到函数返回才关闭！
	//     // 处理文件...
	// }

	// 正确示例：
	for _, filename := range filenames {
		if err := processFile(filename); err != nil {
			return err
		}
	}

	return nil
}

// processFile 处理单个文件（推荐方式）
func processFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close() // 在这个小函数中立即关闭

	// 处理文件...
	return nil
}

// ============ HTTP 响应清理 ============

// HTTPResponseCleanup 演示 HTTP 响应的清理
func HTTPResponseCleanup() error {
	// 模拟 HTTP 响应
	// resp, err := http.Get("https://example.com")
	// if err != nil {
	//     return err
	// }
	// defer resp.Body.Close() // 确保响应体被关闭

	// 实际代码：
	fmt.Println("Making HTTP request...")
	fmt.Println("Response body will be closed automatically")

	return nil
}

// ============ 数据库连接清理 ============

// DatabaseQuery 演示数据库查询的连接清理
func DatabaseQuery() error {
	// 模拟数据库连接
	// db, err := sql.Open("postgres", "connection_string")
	// if err != nil {
	//     return err
	// }
	// defer db.Close()

	// 模拟查询
	// rows, err := db.Query("SELECT * FROM users")
	// if err != nil {
	//     return err
	// }
	// defer rows.Close() // 确保 rows 被关闭

	fmt.Println("Database query with proper cleanup")
	return nil
}
