// Package functional_options 演示 Go 语言特有的函数选项模式 (Functional Options Pattern)
//
// 该模式解决的问题：
// 1. Go 不支持函数重载，无法提供多个不同参数的构造函数
// 2. Go 不支持命名参数，参数多时容易混淆
// 3. 可选参数的优雅处理方式
//
// 使用场景：
// - 配置项较多且大部分有合理默认值的对象创建
// - 需要灵活配置但又不想暴露内部字段
// - 标准库中的广泛应用：http.Server, grpc.ServerOption 等
//
// 核心思想：
// - 定义一个函数类型 Option，接收一个指向目标类型的指针
// - 提供多个返回 Option 的工厂函数
// - 在构造函数中接收可变参数 ...Option，依次应用
package functional_options

import (
	"fmt"
	"time"
)

// Server 代表一个 HTTP 服务器配置
type Server struct {
	Host    string        // 主机地址
	Port    int           // 端口号
	Timeout time.Duration // 超时时间
	MaxConn int           // 最大连接数
	Debug   bool          // 调试模式
	TLS     bool          // 是否启用 TLS
}

// Option 是一个函数类型，用于配置 Server
// 接收 *Server 指针，可以修改其字段
type Option func(*Server)

// WithHost 设置主机地址
func WithHost(host string) Option {
	return func(s *Server) {
		s.Host = host
	}
}

// WithPort 设置端口号
func WithPort(port int) Option {
	return func(s *Server) {
		s.Port = port
	}
}

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.Timeout = timeout
	}
}

// WithMaxConn 设置最大连接数
func WithMaxConn(maxConn int) Option {
	return func(s *Server) {
		s.MaxConn = maxConn
	}
}

// WithDebug 开启调试模式
func WithDebug() Option {
	return func(s *Server) {
		s.Debug = true
	}
}

// WithTLS 启用 TLS
func WithTLS() Option {
	return func(s *Server) {
		s.TLS = true
	}
}

// NewServer 创建新的 Server 实例
// 使用函数选项模式，支持灵活的配置
//
// 示例：
//
//	// 使用默认配置
//	s1 := NewServer()
//
//	// 自定义部分配置
//	s2 := NewServer(
//	    WithHost("localhost"),
//	    WithPort(8080),
//	    WithTimeout(30*time.Second),
//	)
//
//	// 完整配置
//	s3 := NewServer(
//	    WithHost("0.0.0.0"),
//	    WithPort(443),
//	    WithMaxConn(1000),
//	    WithTimeout(time.Minute),
//	    WithTLS(),
//	    WithDebug(),
//	)
func NewServer(opts ...Option) *Server {
	// 1. 设置合理的默认值
	s := &Server{
		Host:    "127.0.0.1",
		Port:    8080,
		Timeout: 30 * time.Second,
		MaxConn: 100,
		Debug:   false,
		TLS:     false,
	}

	// 2. 依次应用所有选项
	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Start 模拟启动服务器
func (s *Server) Start() error {
	protocol := "http"
	if s.TLS {
		protocol = "https"
	}

	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	fmt.Printf("[Server] Starting %s server on %s\n", protocol, addr)
	fmt.Printf("[Server] Config: timeout=%v, maxConn=%d, debug=%v\n",
		s.Timeout, s.MaxConn, s.Debug)

	// 实际业务逻辑...
	return nil
}

// String 返回服务器配置的字符串表示
func (s *Server) String() string {
	return fmt.Sprintf("Server{Host:%s, Port:%d, Timeout:%v, MaxConn:%d, Debug:%v, TLS:%v}",
		s.Host, s.Port, s.Timeout, s.MaxConn, s.Debug, s.TLS)
}
