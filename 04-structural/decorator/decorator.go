// Package decorator 演示装饰器模式 (Decorator Pattern)
//
// 装饰器模式解决的问题：
// 1. 动态地给对象添加额外职责
// 2. 比继承更灵活的扩展方式
// 3. 可以组合多个装饰器
// 4. 遵循开闭原则（对扩展开放，对修改关闭）
//
// 与继承的区别：
// - 继承：静态扩展，编译时确定
// - 装饰器：动态扩展，运行时组合
//
// 使用场景：
// - HTTP 中间件（日志、认证、限流）
// - I/O 流（缓冲、压缩、加密）
// - GUI 组件（边框、滚动条、阴影）
// - 数据处理（过滤、转换、验证）
//
// 核心组件：
// 1. Component 组件接口：定义接口
// 2. ConcreteComponent 具体组件：被装饰的对象
// 3. Decorator 装饰器：实现组件接口，持有组件引用
// 4. ConcreteDecorator 具体装饰器：添加额外行为
package decorator

import (
	"fmt"
	"strings"
	"time"
)

// ============ 示例 1: HTTP 处理器装饰器 ============

// HTTPHandler HTTP 处理器接口
type HTTPHandler interface {
	Handle(request string) string
}

// BasicHandler 基础处理器
type BasicHandler struct {
	name string
}

// NewBasicHandler 创建基础处理器
func NewBasicHandler(name string) *BasicHandler {
	return &BasicHandler{name: name}
}

// Handle 处理请求
func (h *BasicHandler) Handle(request string) string {
	return fmt.Sprintf("[%s] Processing: %s", h.name, request)
}

// HandlerDecorator 处理器装饰器接口
type HandlerDecorator interface {
	Wrap(HTTPHandler) HTTPHandler
}

// LoggingDecorator 日志装饰器
type LoggingDecorator struct{}

// Wrap 包装处理器
func (d *LoggingDecorator) Wrap(handler HTTPHandler) HTTPHandler {
	return &loggingHandler{handler: handler}
}

type loggingHandler struct {
	handler HTTPHandler
}

// Handle 处理请求并记录日志
func (h *loggingHandler) Handle(request string) string {
	fmt.Printf("[LOG] Request: %s\n", request)
	result := h.handler.Handle(request)
	fmt.Printf("[LOG] Response: %s\n", result)
	return result
}

// AuthDecorator 认证装饰器
type AuthDecorator struct{}

// Wrap 包装处理器
func (d *AuthDecorator) Wrap(handler HTTPHandler) HTTPHandler {
	return &authHandler{handler: handler}
}

type authHandler struct {
	handler HTTPHandler
}

// Handle 处理请求前进行认证
func (h *authHandler) Handle(request string) string {
	if !authenticate(request) {
		return "[ERROR] Unauthorized"
	}
	return h.handler.Handle(request)
}

func authenticate(request string) bool {
	return strings.Contains(request, "token=")
}

// TimingDecorator 计时装饰器
type TimingDecorator struct{}

// Wrap 包装处理器
func (d *TimingDecorator) Wrap(handler HTTPHandler) HTTPHandler {
	return &timingHandler{handler: handler}
}

type timingHandler struct {
	handler HTTPHandler
}

// Handle 处理请求并计时
func (h *timingHandler) Handle(request string) string {
	start := time.Now()
	result := h.handler.Handle(request)
	elapsed := time.Since(start)
	fmt.Printf("[TIMING] Took %v\n", elapsed)
	return result
}

// RateLimitDecorator 限流装饰器
type RateLimitDecorator struct {
	maxRequests int
	current     int
}

// NewRateLimitDecorator 创建限流装饰器
func NewRateLimitDecorator(maxRequests int) *RateLimitDecorator {
	return &RateLimitDecorator{maxRequests: maxRequests}
}

// Wrap 包装处理器
func (d *RateLimitDecorator) Wrap(handler HTTPHandler) HTTPHandler {
	return &rateLimitHandler{
		handler:     handler,
		maxRequests: d.maxRequests,
	}
}

type rateLimitHandler struct {
	handler     HTTPHandler
	maxRequests int
	current     int
}

// Handle 处理请求并进行限流
func (h *rateLimitHandler) Handle(request string) string {
	h.current++
	if h.current > h.maxRequests {
		return "[ERROR] Rate limit exceeded"
	}
	return h.handler.Handle(request)
}

// ============ 示例 2: 字符串处理装饰器 ============

// StringProcessor 字符串处理器接口
type StringProcessor interface {
	Process(input string) string
}

// BaseProcessor 基础处理器
type BaseProcessor struct{}

// Process 处理字符串
func (p *BaseProcessor) Process(input string) string {
	return input
}

// ProcessorDecorator 处理器装饰器
type ProcessorDecorator interface {
	StringProcessor
	Decorate(StringProcessor) StringProcessor
}

// UpperCaseDecorator 大写装饰器
type UpperCaseDecorator struct{}

// Decorate 装饰处理器
func (d *UpperCaseDecorator) Decorate(processor StringProcessor) StringProcessor {
	return &upperCaseProcessor{processor: processor}
}

type upperCaseProcessor struct {
	processor StringProcessor
}

// Process 处理并转为大写
func (p *upperCaseProcessor) Process(input string) string {
	result := p.processor.Process(input)
	return strings.ToUpper(result)
}

// TrimDecorator 去除空格装饰器
type TrimDecorator struct{}

// Decorate 装饰处理器
func (d *TrimDecorator) Decorate(processor StringProcessor) StringProcessor {
	return &trimProcessor{processor: processor}
}

type trimProcessor struct {
	processor StringProcessor
}

// Process 处理并去除空格
func (p *trimProcessor) Process(input string) string {
	result := p.processor.Process(input)
	return strings.TrimSpace(result)
}

// ReverseDecorator 反转装饰器
type ReverseDecorator struct{}

// Decorate 装饰处理器
func (d *ReverseDecorator) Decorate(processor StringProcessor) StringProcessor {
	return &reverseProcessor{processor: processor}
}

type reverseProcessor struct {
	processor StringProcessor
}

// Process 处理并反转
func (p *reverseProcessor) Process(input string) string {
	result := p.processor.Process(input)
	runes := []rune(result)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// ============ 示例 3: 数据流装饰器 ============

// DataStream 数据流接口
type DataStream interface {
	Read() string
	Write(data string) error
	Close() error
}

// MemoryStream 内存流
type MemoryStream struct {
	data   string
	closed bool
}

// NewMemoryStream 创建内存流
func NewMemoryStream() *MemoryStream {
	return &MemoryStream{}
}

// Read 读取数据
func (s *MemoryStream) Read() string {
	return s.data
}

// Write 写入数据
func (s *MemoryStream) Write(data string) error {
	if s.closed {
		return fmt.Errorf("stream is closed")
	}
	s.data += data
	return nil
}

// Close 关闭流
func (s *MemoryStream) Close() error {
	s.closed = true
	return nil
}

// StreamDecorator 流装饰器接口
type StreamDecorator interface {
	DataStream
	DecorateStream(DataStream) DataStream
}

// BufferedStream 缓冲流装饰器
type BufferedStream struct {
	bufferSize int
}

// NewBufferedStream 创建缓冲流装饰器
func NewBufferedStream(bufferSize int) *BufferedStream {
	return &BufferedStream{bufferSize: bufferSize}
}

// DecorateStream 装饰流
func (d *BufferedStream) DecorateStream(stream DataStream) DataStream {
	return &bufferedStreamImpl{
		stream:     stream,
		buffer:     make([]string, 0),
		bufferSize: d.bufferSize,
	}
}

type bufferedStreamImpl struct {
	stream     DataStream
	buffer     []string
	bufferSize int
}

// Read 读取数据
func (s *bufferedStreamImpl) Read() string {
	return s.stream.Read()
}

// Write 写入数据（带缓冲）
func (s *bufferedStreamImpl) Write(data string) error {
	s.buffer = append(s.buffer, data)
	if len(s.buffer) >= s.bufferSize {
		return s.flush()
	}
	return nil
}

// Close 关闭并刷新缓冲
func (s *bufferedStreamImpl) Close() error {
	if err := s.flush(); err != nil {
		return err
	}
	return s.stream.Close()
}

func (s *bufferedStreamImpl) flush() error {
	if len(s.buffer) > 0 {
		combined := strings.Join(s.buffer, "")
		if err := s.stream.Write(combined); err != nil {
			return err
		}
		s.buffer = make([]string, 0)
	}
	return nil
}

// EncryptedStream 加密流装饰器
type EncryptedStream struct{}

// DecorateStream 装饰流
func (d *EncryptedStream) DecorateStream(stream DataStream) DataStream {
	return &encryptedStreamImpl{stream: stream}
}

type encryptedStreamImpl struct {
	stream DataStream
}

// Read 读取并解密
func (s *encryptedStreamImpl) Read() string {
	data := s.stream.Read()
	return decrypt(data)
}

// Write 加密后写入
func (s *encryptedStreamImpl) Write(data string) error {
	encrypted := encrypt(data)
	return s.stream.Write(encrypted)
}

// Close 关闭
func (s *encryptedStreamImpl) Close() error {
	return s.stream.Close()
}

func encrypt(data string) string {
	// 简单加密示例
	encrypted := ""
	for _, c := range data {
		encrypted += string(c + 1)
	}
	return encrypted
}

func decrypt(data string) string {
	// 简单解密示例
	decrypted := ""
	for _, c := range data {
		decrypted += string(c - 1)
	}
	return decrypted
}

// ============ 装饰器链 ============

// HandlerChain 处理器装饰器链
type HandlerChain struct {
	decorators []HandlerDecorator
}

// NewHandlerChain 创建装饰器链
func NewHandlerChain() *HandlerChain {
	return &HandlerChain{
		decorators: make([]HandlerDecorator, 0),
	}
}

// AddDecorator 添加装饰器
func (c *HandlerChain) AddDecorator(decorator HandlerDecorator) {
	c.decorators = append(c.decorators, decorator)
}

// Wrap 应用所有装饰器
func (c *HandlerChain) Wrap(handler HTTPHandler) HTTPHandler {
	result := handler
	// 从后往前应用装饰器
	for i := len(c.decorators) - 1; i >= 0; i-- {
		result = c.decorators[i].Wrap(result)
	}
	return result
}
