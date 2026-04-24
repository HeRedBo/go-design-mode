package functional_options

import (
	"testing"
	"time"
)

// TestNewServer_Default 测试默认配置
func TestNewServer_Default(t *testing.T) {
	server := NewServer()

	// 验证默认值
	if server.Host != "127.0.0.1" {
		t.Errorf("expected Host '127.0.0.1', got '%s'", server.Host)
	}
	if server.Port != 8080 {
		t.Errorf("expected Port 8080, got %d", server.Port)
	}
	if server.Timeout != 30*time.Second {
		t.Errorf("expected Timeout 30s, got %v", server.Timeout)
	}
	if server.MaxConn != 100 {
		t.Errorf("expected MaxConn 100, got %d", server.MaxConn)
	}
	if server.Debug != false {
		t.Errorf("expected Debug false, got %v", server.Debug)
	}
	if server.TLS != false {
		t.Errorf("expected TLS false, got %v", server.TLS)
	}
}

// TestNewServer_WithHost 测试自定义 Host
func TestNewServer_WithHost(t *testing.T) {
	server := NewServer(WithHost("0.0.0.0"))

	if server.Host != "0.0.0.0" {
		t.Errorf("expected Host '0.0.0.0', got '%s'", server.Host)
	}
	// 其他字段应保持默认值
	if server.Port != 8080 {
		t.Errorf("expected Port 8080, got %d", server.Port)
	}
}

// TestNewServer_WithPort 测试自定义 Port
func TestNewServer_WithPort(t *testing.T) {
	server := NewServer(WithPort(3000))

	if server.Port != 3000 {
		t.Errorf("expected Port 3000, got %d", server.Port)
	}
}

// TestNewServer_MultipleOptions 测试多个选项组合
func TestNewServer_MultipleOptions(t *testing.T) {
	server := NewServer(
		WithHost("localhost"),
		WithPort(9090),
		WithTimeout(time.Minute),
		WithMaxConn(500),
	)

	if server.Host != "localhost" {
		t.Errorf("expected Host 'localhost', got '%s'", server.Host)
	}
	if server.Port != 9090 {
		t.Errorf("expected Port 9090, got %d", server.Port)
	}
	if server.Timeout != time.Minute {
		t.Errorf("expected Timeout 1m, got %v", server.Timeout)
	}
	if server.MaxConn != 500 {
		t.Errorf("expected MaxConn 500, got %d", server.MaxConn)
	}
}

// TestNewServer_WithOptions 测试布尔选项
func TestNewServer_WithOptions(t *testing.T) {
	server := NewServer(
		WithDebug(),
		WithTLS(),
	)

	if server.Debug != true {
		t.Errorf("expected Debug true, got %v", server.Debug)
	}
	if server.TLS != true {
		t.Errorf("expected TLS true, got %v", server.TLS)
	}
}

// TestNewServer_AllOptions 测试所有选项
func TestNewServer_AllOptions(t *testing.T) {
	server := NewServer(
		WithHost("example.com"),
		WithPort(443),
		WithTimeout(60*time.Second),
		WithMaxConn(1000),
		WithDebug(),
		WithTLS(),
	)

	if server.Host != "example.com" {
		t.Errorf("expected Host 'example.com', got '%s'", server.Host)
	}
	if server.Port != 443 {
		t.Errorf("expected Port 443, got %d", server.Port)
	}
	if server.Timeout != 60*time.Second {
		t.Errorf("expected Timeout 60s, got %v", server.Timeout)
	}
	if server.MaxConn != 1000 {
		t.Errorf("expected MaxConn 1000, got %d", server.MaxConn)
	}
	if server.Debug != true {
		t.Errorf("expected Debug true, got %v", server.Debug)
	}
	if server.TLS != true {
		t.Errorf("expected TLS true, got %v", server.TLS)
	}
}

// TestNewServer_OverrideOptions 测试选项覆盖（后面的覆盖前面的）
func TestNewServer_OverrideOptions(t *testing.T) {
	server := NewServer(
		WithPort(8080),
		WithPort(9090), // 覆盖前面的设置
	)

	if server.Port != 9090 {
		t.Errorf("expected Port 9090 (overridden), got %d", server.Port)
	}
}

// TestServer_Start 测试启动服务器
func TestServer_Start(t *testing.T) {
	server := NewServer(
		WithHost("localhost"),
		WithPort(8080),
		WithTimeout(30*time.Second),
	)

	// Start 应该不返回错误
	err := server.Start()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// TestServer_Start_WithTLS 测试 TLS 服务器启动
func TestServer_Start_WithTLS(t *testing.T) {
	server := NewServer(
		WithHost("localhost"),
		WithPort(443),
		WithTLS(),
	)

	// Start 应该不返回错误
	err := server.Start()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// TestServer_String 测试 String 方法
func TestServer_String(t *testing.T) {
	server := NewServer(
		WithHost("localhost"),
		WithPort(8080),
		WithTimeout(30*time.Second),
		WithMaxConn(100),
		WithDebug(),
		WithTLS(),
	)

	expected := "Server{Host:localhost, Port:8080, Timeout:30s, MaxConn:100, Debug:true, TLS:true}"
	if server.String() != expected {
		t.Errorf("expected String '%s', got '%s'", expected, server.String())
	}
}

// TestNewServer_NoOptions 测试空选项（不传任何选项）
func TestNewServer_NoOptions(t *testing.T) {
	server := NewServer()

	// 应该使用所有默认值
	expected := &Server{
		Host:    "127.0.0.1",
		Port:    8080,
		Timeout: 30 * time.Second,
		MaxConn: 100,
		Debug:   false,
		TLS:     false,
	}

	if *server != *expected {
		t.Errorf("expected %+v, got %+v", expected, server)
	}
}

// TestOption_Functionality 测试 Option 函数的正确性
func TestOption_Functionality(t *testing.T) {
	tests := []struct {
		name     string
		opts     []Option
		expected *Server
	}{
		{
			name: "只设置 Host",
			opts: []Option{WithHost("test.com")},
			expected: &Server{
				Host:    "test.com",
				Port:    8080,
				Timeout: 30 * time.Second,
				MaxConn: 100,
				Debug:   false,
				TLS:     false,
			},
		},
		{
			name: "只设置 Timeout",
			opts: []Option{WithTimeout(time.Hour)},
			expected: &Server{
				Host:    "127.0.0.1",
				Port:    8080,
				Timeout: time.Hour,
				MaxConn: 100,
				Debug:   false,
				TLS:     false,
			},
		},
		{
			name: "开启 Debug 和 TLS",
			opts: []Option{WithDebug(), WithTLS()},
			expected: &Server{
				Host:    "127.0.0.1",
				Port:    8080,
				Timeout: 30 * time.Second,
				MaxConn: 100,
				Debug:   true,
				TLS:     true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewServer(tt.opts...)

			if server.Host != tt.expected.Host {
				t.Errorf("Host: expected '%s', got '%s'", tt.expected.Host, server.Host)
			}
			if server.Port != tt.expected.Port {
				t.Errorf("Port: expected %d, got %d", tt.expected.Port, server.Port)
			}
			if server.Timeout != tt.expected.Timeout {
				t.Errorf("Timeout: expected %v, got %v", tt.expected.Timeout, server.Timeout)
			}
			if server.MaxConn != tt.expected.MaxConn {
				t.Errorf("MaxConn: expected %d, got %d", tt.expected.MaxConn, server.MaxConn)
			}
			if server.Debug != tt.expected.Debug {
				t.Errorf("Debug: expected %v, got %v", tt.expected.Debug, server.Debug)
			}
			if server.TLS != tt.expected.TLS {
				t.Errorf("TLS: expected %v, got %v", tt.expected.TLS, server.TLS)
			}
		})
	}
}
