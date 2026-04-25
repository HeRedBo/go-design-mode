package adapter

import (
	"testing"
)

// ============ 支付系统适配器测试 ============

func TestLegacyPaymentAdapter(t *testing.T) {
	adapter := NewLegacyPaymentAdapter()

	if adapter.GetName() != "Legacy Payment System (Adapted)" {
		t.Errorf("Unexpected name: %s", adapter.GetName())
	}

	err := adapter.Pay(100.50, "USD")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// 测试负数支付
	err = adapter.Pay(-10.0, "USD")
	if err == nil {
		t.Error("Expected error for negative payment")
	}
}

func TestModernPayment(t *testing.T) {
	payment := &ModernPayment{Name: "Stripe"}

	if payment.GetName() != "Stripe" {
		t.Errorf("Expected name 'Stripe', got '%s'", payment.GetName())
	}

	err := payment.Pay(50.0, "EUR")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestProcessPayment(t *testing.T) {
	// 测试旧版支付适配器
	legacyAdapter := NewLegacyPaymentAdapter()
	err := ProcessPayment(legacyAdapter, 100.0, "USD")
	if err != nil {
		t.Errorf("Legacy payment failed: %v", err)
	}

	// 测试现代支付
	modernPayment := &ModernPayment{Name: "PayPal"}
	err = ProcessPayment(modernPayment, 200.0, "EUR")
	if err != nil {
		t.Errorf("Modern payment failed: %v", err)
	}
}

// ============ 消息队列适配器测试 ============

func TestRabbitMQAdapter(t *testing.T) {
	adapter := NewRabbitMQAdapter()

	if adapter.GetType() != "RabbitMQ (Adapted)" {
		t.Errorf("Unexpected type: %s", adapter.GetType())
	}

	message := map[string]interface{}{
		"user_id":  123,
		"action":   "create",
		"priority": "high",
	}

	err := adapter.Send("user_events", message)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestKafkaProducer(t *testing.T) {
	producer := &KafkaProducer{}

	if producer.GetType() != "Kafka" {
		t.Errorf("Expected type 'Kafka', got '%s'", producer.GetType())
	}

	message := map[string]interface{}{
		"order_id": 456,
		"status":   "shipped",
	}

	err := producer.Send("order_updates", message)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

// ============ 存储系统适配器测试 ============

func TestFTPStorageAdapter(t *testing.T) {
	adapter := NewFTPStorageAdapter()

	if adapter.GetType() != "FTP Storage (Adapted)" {
		t.Errorf("Unexpected type: %s", adapter.GetType())
	}

	// 测试保存
	err := adapter.Save("test.txt", []byte("hello world"))
	if err != nil {
		t.Errorf("Save failed: %v", err)
	}

	// 测试加载
	data, err := adapter.Load("test.txt")
	if err != nil {
		t.Errorf("Load failed: %v", err)
	}
	if string(data) != "file content" {
		t.Errorf("Unexpected data: %s", string(data))
	}

	// 测试删除（应该失败）
	err = adapter.Delete("test.txt")
	if err == nil {
		t.Error("Expected error for delete operation")
	}
}

func TestS3Storage(t *testing.T) {
	storage := &S3Storage{}

	if storage.GetType() != "S3 Storage" {
		t.Errorf("Expected type 'S3 Storage', got '%s'", storage.GetType())
	}

	// 测试保存
	err := storage.Save("file.pdf", []byte("pdf content"))
	if err != nil {
		t.Errorf("Save failed: %v", err)
	}

	// 测试加载
	_, err = storage.Load("file.pdf")
	if err != nil {
		t.Errorf("Load failed: %v", err)
	}

	// 测试删除
	err = storage.Delete("file.pdf")
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}
}

// ============ 字符串处理适配器测试 ============

func TestToUpperProcessor(t *testing.T) {
	processor := &ToUpperProcessor{}

	if processor.GetName() != "ToUpper" {
		t.Errorf("Expected name 'ToUpper', got '%s'", processor.GetName())
	}

	result := processor.Process("hello world")
	if result != "HELLO WORLD" {
		t.Errorf("Expected 'HELLO WORLD', got '%s'", result)
	}
}

func TestTrimProcessor(t *testing.T) {
	processor := &TrimProcessor{}

	if processor.GetName() != "Trim" {
		t.Errorf("Expected name 'Trim', got '%s'", processor.GetName())
	}

	result := processor.Process("  hello  ")
	if result != "hello" {
		t.Errorf("Expected 'hello', got '%s'", result)
	}
}

func TestReverseProcessor(t *testing.T) {
	processor := &ReverseProcessor{}

	if processor.GetName() != "Reverse" {
		t.Errorf("Expected name 'Reverse', got '%s'", processor.GetName())
	}

	result := processor.Process("hello")
	if result != "olleh" {
		t.Errorf("Expected 'olleh', got '%s'", result)
	}

	// 测试中文
	result = processor.Process("你好")
	if result != "好你" {
		t.Errorf("Expected '好你', got '%s'", result)
	}
}

func TestProcessorChain(t *testing.T) {
	chain := NewProcessorChain()
	chain.AddProcessor(&ToUpperProcessor{})
	chain.AddProcessor(&ReverseProcessor{})

	if chain.GetName() != "Chain[ToUpper -> Reverse]" {
		t.Errorf("Unexpected chain name: %s", chain.GetName())
	}

	result := chain.Process("hello")
	// "hello" -> "HELLO" -> "OLLEH"
	if result != "OLLEH" {
		t.Errorf("Expected 'OLLEH', got '%s'", result)
	}
}

func TestStringAdapter(t *testing.T) {
	// 将普通函数适配为 StringProcessor
	replaceSpaces := func(s string) string {
		result := ""
		for _, c := range s {
			if c == ' ' {
				result += "_"
			} else {
				result += string(c)
			}
		}
		return result
	}

	adapter := NewStringAdapter("ReplaceSpaces", replaceSpaces)

	if adapter.GetName() != "ReplaceSpaces" {
		t.Errorf("Expected name 'ReplaceSpaces', got '%s'", adapter.GetName())
	}

	result := adapter.Process("hello world foo")
	if result != "hello_world_foo" {
		t.Errorf("Expected 'hello_world_foo', got '%s'", result)
	}
}

func TestComplexProcessorChain(t *testing.T) {
	chain := NewProcessorChain()
	chain.AddProcessor(&TrimProcessor{})
	chain.AddProcessor(&ToUpperProcessor{})

	// 使用适配器添加自定义处理器
	replaceSpaces := NewStringAdapter("ReplaceSpaces", func(s string) string {
		result := ""
		for _, c := range s {
			if c == ' ' {
				result += "-"
			} else {
				result += string(c)
			}
		}
		return result
	})
	chain.AddProcessor(replaceSpaces)

	if chain.GetName() != "Chain[Trim -> ToUpper -> ReplaceSpaces]" {
		t.Errorf("Unexpected chain name: %s", chain.GetName())
	}

	result := chain.Process("  hello world  ")
	// "  hello world  " -> "hello world" -> "HELLO WORLD" -> "HELLO-WORLD"
	if result != "HELLO-WORLD" {
		t.Errorf("Expected 'HELLO-WORLD', got '%s'", result)
	}
}

// ============ 基准测试 ============

func BenchmarkLegacyPaymentAdapter(b *testing.B) {
	adapter := NewLegacyPaymentAdapter()
	for i := 0; i < b.N; i++ {
		_ = adapter.Pay(100.0, "USD")
	}
}

func BenchmarkRabbitMQAdapter(b *testing.B) {
	adapter := NewRabbitMQAdapter()
	message := map[string]interface{}{"data": "test"}
	for i := 0; i < b.N; i++ {
		_ = adapter.Send("queue", message)
	}
}

func BenchmarkProcessorChain(b *testing.B) {
	chain := NewProcessorChain()
	chain.AddProcessor(&TrimProcessor{})
	chain.AddProcessor(&ToUpperProcessor{})
	chain.AddProcessor(&ReverseProcessor{})

	for i := 0; i < b.N; i++ {
		_ = chain.Process("  hello world  ")
	}
}
