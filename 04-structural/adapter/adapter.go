// Package adapter 演示适配器模式 (Adapter Pattern)
//
// 适配器模式解决的问题：
// 1. 兼容不兼容的接口，使原本不能一起工作的类可以一起工作
// 2. 复用现有代码，无需修改原有类
// 3. 在客户端和目标类之间提供转换层
//
// 使用场景：
// - 集成第三方库，接口不匹配
// - 新旧系统兼容
// - 统一不同支付/消息/存储服务的接口
//
// Go 中的实现特点：
// - 使用接口定义目标接口
// - 适配器实现目标接口，内部调用被适配者
// - 对象适配器（组合）比类适配器（继承）更常用
//
// 适配器类型：
// 1. 对象适配器：通过组合持有被适配者对象
// 2. 接口适配器：实现接口的大部分方法，只提供默认实现
//
// 核心组件：
// 1. Target 接口：客户端期望的接口
// 2. Adaptee 被适配者：已存在的接口/类
// 3. Adapter 适配器：实现 Target 接口，包装 Adaptee
package adapter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ============ 示例 1: 支付系统适配器 ============

// LegacyPayment 旧版支付系统（不兼容的接口）
type LegacyPayment struct{}

// MakeOldPayment 旧版支付方法
func (p *LegacyPayment) MakeOldPayment(amount float64, currency string, description string) bool {
	fmt.Printf("[LegacyPayment] Processing payment: %.2f %s - %s\n", amount, currency, description)
	return amount > 0
}

// PaymentProcessor 目标支付接口（新版标准）
type PaymentProcessor interface {
	Pay(amount float64, currency string) error
	GetName() string
}

// LegacyPaymentAdapter 旧版支付适配器
type LegacyPaymentAdapter struct {
	legacy *LegacyPayment
}

// NewLegacyPaymentAdapter 创建旧版支付适配器
func NewLegacyPaymentAdapter() *LegacyPaymentAdapter {
	return &LegacyPaymentAdapter{
		legacy: &LegacyPayment{},
	}
}

// Pay 实现 PaymentProcessor 接口
func (a *LegacyPaymentAdapter) Pay(amount float64, currency string) error {
	success := a.legacy.MakeOldPayment(amount, currency, "Adapted payment")
	if !success {
		return fmt.Errorf("payment failed")
	}
	return nil
}

// GetName 返回支付系统名称
func (a *LegacyPaymentAdapter) GetName() string {
	return "Legacy Payment System (Adapted)"
}

// ModernPayment 现代支付系统（已经兼容）
type ModernPayment struct {
	Name string
}

// Pay 实现 PaymentProcessor 接口
func (p *ModernPayment) Pay(amount float64, currency string) error {
	fmt.Printf("[ModernPayment] Processing payment: %.2f %s\n", amount, currency)
	return nil
}

// GetName 返回支付系统名称
func (p *ModernPayment) GetName() string {
	return p.Name
}

// ProcessPayment 使用支付处理器（客户端代码）
func ProcessPayment(processor PaymentProcessor, amount float64, currency string) error {
	fmt.Printf("Using payment processor: %s\n", processor.GetName())
	return processor.Pay(amount, currency)
}

// ============ 示例 2: 消息队列适配器 ============

// RabbitMQMessage 旧版 RabbitMQ 消息格式
type RabbitMQMessage struct {
	Body    string
	Headers map[string]string
}

// RabbitMQProducer 旧版 RabbitMQ 生产者
type RabbitMQProducer struct{}

// Publish 发布消息（旧版接口）
func (p *RabbitMQProducer) Publish(exchange, routingKey string, message RabbitMQMessage) error {
	fmt.Printf("[RabbitMQ] Publishing to %s/%s: %s\n", exchange, routingKey, message.Body)
	return nil
}

// MessageProducer 目标消息生产者接口
type MessageProducer interface {
	Send(queue string, message map[string]interface{}) error
	GetType() string
}

// RabbitMQAdapter RabbitMQ 适配器
type RabbitMQAdapter struct {
	producer *RabbitMQProducer
}

// NewRabbitMQAdapter 创建 RabbitMQ 适配器
func NewRabbitMQAdapter() *RabbitMQAdapter {
	return &RabbitMQAdapter{
		producer: &RabbitMQProducer{},
	}
}

// Send 实现 MessageProducer 接口
func (a *RabbitMQAdapter) Send(queue string, message map[string]interface{}) error {
	// 转换消息格式
	body, _ := json.Marshal(message)
	rmqMessage := RabbitMQMessage{
		Body:    string(body),
		Headers: make(map[string]string),
	}

	return a.producer.Publish("", queue, rmqMessage)
}

// GetType 返回类型
func (a *RabbitMQAdapter) GetType() string {
	return "RabbitMQ (Adapted)"
}

// KafkaProducer 现代 Kafka 生产者（已兼容）
type KafkaProducer struct{}

// Send 实现 MessageProducer 接口
func (p *KafkaProducer) Send(queue string, message map[string]interface{}) error {
	fmt.Printf("[Kafka] Sending to %s: %v\n", queue, message)
	return nil
}

// GetType 返回类型
func (p *KafkaProducer) GetType() string {
	return "Kafka"
}

// ============ 示例 3: 存储系统适配器 ============

// FTPStorage 旧版 FTP 存储
type FTPStorage struct{}

// UploadFile FTP 上传文件（旧版接口）
func (f *FTPStorage) UploadFile(path string, data []byte, mode string) error {
	fmt.Printf("[FTP] Uploading %d bytes to %s (mode: %s)\n", len(data), path, mode)
	return nil
}

// DownloadFile FTP 下载文件
func (f *FTPStorage) DownloadFile(path string) ([]byte, error) {
	fmt.Printf("[FTP] Downloading from %s\n", path)
	return []byte("file content"), nil
}

// Storage 目标存储接口
type Storage interface {
	Save(key string, data []byte) error
	Load(key string) ([]byte, error)
	Delete(key string) error
	GetType() string
}

// FTPStorageAdapter FTP 存储适配器
type FTPStorageAdapter struct {
	ftp *FTPStorage
}

// NewFTPStorageAdapter 创建 FTP 存储适配器
func NewFTPStorageAdapter() *FTPStorageAdapter {
	return &FTPStorageAdapter{
		ftp: &FTPStorage{},
	}
}

// Save 实现 Storage 接口
func (a *FTPStorageAdapter) Save(key string, data []byte) error {
	return a.ftp.UploadFile(key, data, "binary")
}

// Load 实现 Storage 接口
func (a *FTPStorageAdapter) Load(key string) ([]byte, error) {
	return a.ftp.DownloadFile(key)
}

// Delete 实现 Storage 接口（FTP 不支持，返回错误）
func (a *FTPStorageAdapter) Delete(key string) error {
	return fmt.Errorf("FTP storage does not support delete operation")
}

// GetType 返回类型
func (a *FTPStorageAdapter) GetType() string {
	return "FTP Storage (Adapted)"
}

// S3Storage 现代 S3 存储（已兼容）
type S3Storage struct{}

// Save 实现 Storage 接口
func (s *S3Storage) Save(key string, data []byte) error {
	fmt.Printf("[S3] Saving %d bytes to %s\n", len(data), key)
	return nil
}

// Load 实现 Storage 接口
func (s *S3Storage) Load(key string) ([]byte, error) {
	fmt.Printf("[S3] Loading from %s\n", key)
	return []byte("s3 content"), nil
}

// Delete 实现 Storage 接口
func (s *S3Storage) Delete(key string) error {
	fmt.Printf("[S3] Deleting %s\n", key)
	return nil
}

// GetType 返回类型
func (s *S3Storage) GetType() string {
	return "S3 Storage"
}

// ============ 示例 4: 字符串处理适配器 ============

// StringProcessor 字符串处理器接口
type StringProcessor interface {
	Process(input string) string
	GetName() string
}

// ToUpperProcessor 转大写处理器
type ToUpperProcessor struct{}

// Process 处理字符串
func (p *ToUpperProcessor) Process(input string) string {
	return strings.ToUpper(input)
}

// GetName 返回名称
func (p *ToUpperProcessor) GetName() string {
	return "ToUpper"
}

// TrimProcessor 去除空格处理器
type TrimProcessor struct{}

// Process 处理字符串
func (p *TrimProcessor) Process(input string) string {
	return strings.TrimSpace(input)
}

// GetName 返回名称
func (p *TrimProcessor) GetName() string {
	return "Trim"
}

// ReverseProcessor 反转字符串处理器
type ReverseProcessor struct{}

// Process 处理字符串
func (p *ReverseProcessor) Process(input string) string {
	runes := []rune(input)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// GetName 返回名称
func (p *ReverseProcessor) GetName() string {
	return "Reverse"
}

// ProcessorChain 处理器链（组合多个处理器）
type ProcessorChain struct {
	processors []StringProcessor
}

// NewProcessorChain 创建处理器链
func NewProcessorChain() *ProcessorChain {
	return &ProcessorChain{
		processors: make([]StringProcessor, 0),
	}
}

// AddProcessor 添加处理器
func (pc *ProcessorChain) AddProcessor(processor StringProcessor) {
	pc.processors = append(pc.processors, processor)
}

// Process 依次执行所有处理器
func (pc *ProcessorChain) Process(input string) string {
	result := input
	for _, processor := range pc.processors {
		result = processor.Process(result)
	}
	return result
}

// GetName 返回名称
func (pc *ProcessorChain) GetName() string {
	names := make([]string, len(pc.processors))
	for i, p := range pc.processors {
		names[i] = p.GetName()
	}
	return fmt.Sprintf("Chain[%s]", strings.Join(names, " -> "))
}

// StringAdapter 字符串处理适配器（将普通函数转换为处理器）
type StringAdapter struct {
	name    string
	handler func(string) string
}

// NewStringAdapter 创建字符串适配器
func NewStringAdapter(name string, handler func(string) string) StringProcessor {
	return &StringAdapter{
		name:    name,
		handler: handler,
	}
}

// Process 实现 StringProcessor 接口
func (a *StringAdapter) Process(input string) string {
	return a.handler(input)
}

// GetName 返回名称
func (a *StringAdapter) GetName() string {
	return a.name
}
