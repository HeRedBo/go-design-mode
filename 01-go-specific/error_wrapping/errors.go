// Package errorwrapping 演示 Go 1.13+ 的错误包装特性
//
// Go 错误处理的独特之处：
// 1. 错误是一等公民，作为返回值传递
// 2. fmt.Errorf 支持 %w 动词包装错误
// 3. errors.Is() 判断错误链中是否包含特定错误
// 4. errors.As() 从错误链中提取特定类型的错误
// 5. 支持自定义错误类型，提供上下文信息
//
// 与传统方式的区别：
// - 旧方式：使用 errors.New() 或 fmt.Errorf() 创建错误，无法追踪错误链
// - 新方式：使用 %w 包装错误，保持错误链，可追溯根本原因
//
// 最佳实践：
// 1. 只在需要添加上下文时才包装错误
// 2. 使用哨兵错误 (sentinel errors) 作为可检查的错误值
// 3. 定义自定义错误类型以携带额外信息
// 4. 使用 errors.Is() 和 errors.As() 进行错误判断
// 5. 不要包装已经包含足够上下文的错误
package errorwrapping

import (
	"errors"
	"fmt"
)

// ============ 哨兵错误 (Sentinel Errors) ============

// 预定义的错误值，用于错误比较
var (
	ErrNotFound     = errors.New("resource not found")
	ErrUnauthorized = errors.New("unauthorized access")
	ErrInvalidInput = errors.New("invalid input")
	ErrDatabase     = errors.New("database error")
)

// ============ 自定义错误类型 ============

// ValidationError 验证错误类型，携带字段信息
type ValidationError struct {
	Field   string
	Message string
	Err     error // 原始错误
}

// Error 实现 error 接口
func (e *ValidationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("validation error on field '%s': %s (cause: %v)",
			e.Field, e.Message, e.Err)
	}
	return fmt.Sprintf("validation error on field '%s': %s",
		e.Field, e.Message)
}

// Unwrap 实现错误解包，支持 errors.Is 和 errors.As
func (e *ValidationError) Unwrap() error {
	return e.Err
}

// DatabaseError 数据库错误类型，携带更多信息
type DatabaseError struct {
	Operation string // 操作类型：query, insert, update, delete
	Query     string // SQL 查询
	Err       error  // 原始错误
}

// Error 实现 error 接口
func (e *DatabaseError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("database error during %s: %v", e.Operation, e.Err)
	}
	return fmt.Sprintf("database error during %s", e.Operation)
}

// Unwrap 实现错误解包
func (e *DatabaseError) Unwrap() error {
	return e.Err
}

// ============ 错误包装示例函数 ============

// FindUser 查找用户 - 演示基本错误包装
func FindUser(id int) (string, error) {
	if id <= 0 {
		// 包装错误，添加上下文
		return "", fmt.Errorf("FindUser(%d): %w", id, ErrInvalidInput)
	}

	if id > 1000 {
		// 包装错误，模拟未找到
		return "", fmt.Errorf("FindUser: %w (id=%d)", ErrNotFound, id)
	}

	return "User_" + fmt.Sprint(id), nil
}

// AuthenticateUser 认证用户 - 演示多层错误包装
func AuthenticateUser(username, password string) error {
	if username == "" {
		return fmt.Errorf("authenticate: username is empty: %w", ErrInvalidInput)
	}

	if password == "" {
		return fmt.Errorf("authenticate: password is empty: %w", ErrInvalidInput)
	}

	// 模拟认证失败
	if username != "admin" || password != "secret" {
		return fmt.Errorf("authenticate(%s): %w", username, ErrUnauthorized)
	}

	return nil
}

// QueryDatabase 查询数据库 - 演示自定义错误类型
func QueryDatabase(query string) ([]string, error) {
	if query == "" {
		return nil, &DatabaseError{
			Operation: "query",
			Query:     query,
			Err:       ErrInvalidInput,
		}
	}

	// 模拟查询失败
	if query == "FAIL" {
		return nil, &DatabaseError{
			Operation: "query",
			Query:     query,
			Err:       errors.New("connection timeout"),
		}
	}

	return []string{"result1", "result2"}, nil
}

// ValidateUser 验证用户 - 演示 ValidationError
func ValidateUser(name string, age int) error {
	if name == "" {
		return &ValidationError{
			Field:   "name",
			Message: "name cannot be empty",
		}
	}

	if age < 0 || age > 150 {
		return &ValidationError{
			Field:   "age",
			Message: fmt.Sprintf("age %d is out of range [0, 150]", age),
		}
	}

	return nil
}

// ProcessOrder 处理订单 - 演示复杂错误处理流程
func ProcessOrder(orderID int) error {
	// 步骤 1: 验证订单
	if orderID <= 0 {
		return fmt.Errorf("process order: invalid order ID: %w", ErrInvalidInput)
	}

	// 步骤 2: 查找订单
	_, err := FindUser(orderID)
	if err != nil {
		return fmt.Errorf("process order: find user failed: %w", err)
	}

	// 步骤 3: 查询数据库
	_, err = QueryDatabase("FAIL")
	if err != nil {
		return fmt.Errorf("process order: database query failed: %w", err)
	}

	return nil
}

// ============ 错误处理辅助函数 ============

// HandleError 演示如何使用 errors.Is() 和 errors.As() 处理错误
func HandleError(err error) string {
	if err == nil {
		return "no error"
	}

	// 使用 errors.Is() 检查错误链中是否包含特定错误
	if errors.Is(err, ErrNotFound) {
		return "Resource not found"
	}

	if errors.Is(err, ErrUnauthorized) {
		return "Unauthorized access denied"
	}

	if errors.Is(err, ErrInvalidInput) {
		return "Invalid input provided"
	}

	// 使用 errors.As() 提取特定类型的错误
	var dbErr *DatabaseError
	if errors.As(err, &dbErr) {
		return fmt.Sprintf("Database error: operation=%s, query=%s",
			dbErr.Operation, dbErr.Query)
	}

	var valErr *ValidationError
	if errors.As(err, &valErr) {
		return fmt.Sprintf("Validation error: field=%s, message=%s",
			valErr.Field, valErr.Message)
	}

	// 默认错误信息
	return fmt.Sprintf("Unknown error: %v", err)
}

// GetRootCause 获取错误的根本原因
func GetRootCause(err error) error {
	for err != nil {
		unwrapped := errors.Unwrap(err)
		if unwrapped == nil {
			return err
		}
		err = unwrapped
	}
	return err
}
