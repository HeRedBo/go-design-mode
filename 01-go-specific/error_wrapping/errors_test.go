package errorwrapping

import (
	"errors"
	"fmt"
	"testing"
)

// TestSentinelErrors 测试哨兵错误
func TestSentinelErrors(t *testing.T) {
	// 验证哨兵错误是否定义正确
	tests := []struct {
		name string
		err  error
	}{
		{"ErrNotFound", ErrNotFound},
		{"ErrUnauthorized", ErrUnauthorized},
		{"ErrInvalidInput", ErrInvalidInput},
		{"ErrDatabase", ErrDatabase},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Errorf("expected error to be defined, got nil")
			}
			if tt.err.Error() == "" {
				t.Errorf("expected error message to be non-empty")
			}
		})
	}
}

// TestFindUser 测试用户查找和错误包装
func TestFindUser(t *testing.T) {
	tests := []struct {
		name        string
		id          int
		expectError bool
		expectUser  string
	}{
		{"valid ID", 100, false, "User_100"},
		{"invalid ID (negative)", -1, true, ""},
		{"invalid ID (zero)", 0, true, ""},
		{"not found (too large)", 2000, true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := FindUser(tt.id)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				// 验证错误包装是否正确
				if !errors.Is(err, ErrInvalidInput) && !errors.Is(err, ErrNotFound) {
					t.Errorf("expected error to wrap ErrInvalidInput or ErrNotFound, got: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if user != tt.expectUser {
					t.Errorf("expected user '%s', got '%s'", tt.expectUser, user)
				}
			}
		})
	}
}

// TestAuthenticateUser 测试用户认证
func TestAuthenticateUser(t *testing.T) {
	tests := []struct {
		name        string
		username    string
		password    string
		expectError bool
	}{
		{"valid credentials", "admin", "secret", false},
		{"empty username", "", "secret", true},
		{"empty password", "admin", "", true},
		{"invalid credentials", "user", "wrong", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := AuthenticateUser(tt.username, tt.password)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				// 验证错误类型
				if !errors.Is(err, ErrInvalidInput) && !errors.Is(err, ErrUnauthorized) {
					t.Errorf("expected error to wrap ErrInvalidInput or ErrUnauthorized, got: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}

// TestQueryDatabase 测试数据库查询
func TestQueryDatabase(t *testing.T) {
	tests := []struct {
		name          string
		query         string
		expectError   bool
		expectResults []string
	}{
		{"valid query", "SELECT * FROM users", false, []string{"result1", "result2"}},
		{"empty query", "", true, nil},
		{"failing query", "FAIL", true, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := QueryDatabase(tt.query)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				// 验证错误类型是 DatabaseError
				var dbErr *DatabaseError
				if !errors.As(err, &dbErr) {
					t.Errorf("expected error to be *DatabaseError, got: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if len(results) != len(tt.expectResults) {
					t.Errorf("expected %d results, got %d", len(tt.expectResults), len(results))
				}
			}
		})
	}
}

// TestValidateUser 测试用户验证
func TestValidateUser(t *testing.T) {
	tests := []struct {
		name        string
		userName    string
		age         int
		expectError bool
	}{
		{"valid user", "Alice", 25, false},
		{"empty name", "", 25, true},
		{"negative age", "Bob", -1, true},
		{"too old", "Charlie", 200, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUser(tt.userName, tt.age)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				// 验证错误类型是 ValidationError
				var valErr *ValidationError
				if !errors.As(err, &valErr) {
					t.Errorf("expected error to be *ValidationError, got: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}

// TestProcessOrder 测试订单处理（复杂错误链）
func TestProcessOrder(t *testing.T) {
	tests := []struct {
		name        string
		orderID     int
		expectError bool
	}{
		{"invalid order ID", -1, true},
		{"valid order but db fails", 100, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ProcessOrder(tt.orderID)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}

// TestHandleError 测试错误处理辅助函数
func TestHandleError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: "no error",
		},
		{
			name:     "ErrNotFound",
			err:      ErrNotFound,
			expected: "Resource not found",
		},
		{
			name:     "wrapped ErrNotFound",
			err:      errors.New("context: " + ErrNotFound.Error()),
			expected: "Unknown error: context: resource not found",
		},
		{
			name:     "wrapped with %w",
			err:      errors.New("find user: " + ErrNotFound.Error()),
			expected: "Unknown error: find user: resource not found",
		},
		{
			name:     "ErrUnauthorized",
			err:      ErrUnauthorized,
			expected: "Unauthorized access denied",
		},
		{
			name:     "ErrInvalidInput",
			err:      ErrInvalidInput,
			expected: "Invalid input provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HandleError(tt.err)
			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// TestHandleError_WithDatabaseError 测试处理数据库错误
func TestHandleError_WithDatabaseError(t *testing.T) {
	dbErr := &DatabaseError{
		Operation: "query",
		Query:     "SELECT * FROM users",
		Err:       errors.New("timeout"),
	}

	result := HandleError(dbErr)
	expected := "Database error: operation=query, query=SELECT * FROM users"

	if result != expected {
		t.Errorf("expected '%s', got '%s'", expected, result)
	}
}

// TestHandleError_WithValidationError 测试处理验证错误
func TestHandleError_WithValidationError(t *testing.T) {
	valErr := &ValidationError{
		Field:   "age",
		Message: "age must be positive",
	}

	result := HandleError(valErr)
	expected := "Validation error: field=age, message=age must be positive"

	if result != expected {
		t.Errorf("expected '%s', got '%s'", expected, result)
	}
}

// TestGetRootCause 测试获取错误根本原因
func TestGetRootCause(t *testing.T) {
	// 创建多层包装的错误
	rootErr := errors.New("root cause")
	wrapped1 := fmt.Errorf("layer 1: %w", rootErr)
	wrapped2 := fmt.Errorf("layer 2: %w", wrapped1)
	wrapped3 := fmt.Errorf("layer 3: %w", wrapped2)

	// 测试获取根本原因
	cause := GetRootCause(wrapped3)
	if cause != rootErr {
		t.Errorf("expected root cause to be 'root cause', got: %v", cause)
	}

	// 测试无包装的错误
	simpleErr := errors.New("simple error")
	cause = GetRootCause(simpleErr)
	if cause != simpleErr {
		t.Errorf("expected root cause to be 'simple error', got: %v", cause)
	}

	// 测试 nil
	cause = GetRootCause(nil)
	if cause != nil {
		t.Errorf("expected root cause to be nil, got: %v", cause)
	}
}

// TestErrorsIs 测试 errors.Is 的功能
func TestErrorsIs(t *testing.T) {
	// 创建包装错误
	err := fmt.Errorf("FindUser(100): %w", ErrNotFound)

	// 验证 errors.Is 能正确识别包装链中的错误
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected errors.Is to find ErrNotFound in error chain")
	}

	// 验证不匹配的错误
	if errors.Is(err, ErrUnauthorized) {
		t.Errorf("expected errors.Is to NOT find ErrUnauthorized in error chain")
	}
}

// TestErrorsAs 测试 errors.As 的功能
func TestErrorsAs(t *testing.T) {
	// 创建自定义错误类型
	dbErr := &DatabaseError{
		Operation: "insert",
		Query:     "INSERT INTO users",
		Err:       errors.New("duplicate key"),
	}

	// 包装错误
	wrapped := fmt.Errorf("process order: %w", dbErr)

	// 使用 errors.As 提取错误类型
	var extracted *DatabaseError
	if !errors.As(wrapped, &extracted) {
		t.Fatalf("expected errors.As to extract *DatabaseError")
	}

	// 验证提取的错误
	if extracted.Operation != "insert" {
		t.Errorf("expected Operation 'insert', got '%s'", extracted.Operation)
	}
	if extracted.Query != "INSERT INTO users" {
		t.Errorf("expected Query 'INSERT INTO users', got '%s'", extracted.Query)
	}
}

// TestValidationError_Error 测试 ValidationError 的 Error 方法
func TestValidationError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *ValidationError
		expected string
	}{
		{
			name: "with cause",
			err: &ValidationError{
				Field:   "email",
				Message: "invalid format",
				Err:     errors.New("regex mismatch"),
			},
			expected: "validation error on field 'email': invalid format (cause: regex mismatch)",
		},
		{
			name: "without cause",
			err: &ValidationError{
				Field:   "name",
				Message: "required field",
			},
			expected: "validation error on field 'name': required field",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()
			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// TestDatabaseError_Error 测试 DatabaseError 的 Error 方法
func TestDatabaseError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *DatabaseError
		expected string
	}{
		{
			name: "with cause",
			err: &DatabaseError{
				Operation: "update",
				Query:     "UPDATE users SET name='test'",
				Err:       errors.New("lock timeout"),
			},
			expected: "database error during update: lock timeout",
		},
		{
			name: "without cause",
			err: &DatabaseError{
				Operation: "delete",
				Query:     "DELETE FROM users WHERE id=1",
			},
			expected: "database error during delete",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()
			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
