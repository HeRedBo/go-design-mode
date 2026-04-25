package decorator

import "testing"

func TestLoggingDecorator(t *testing.T) {
	handler := NewBasicHandler("API")
	decorator := &LoggingDecorator{}
	wrapped := decorator.Wrap(handler)

	result := wrapped.Handle("GET /users")
	if result == "" {
		t.Error("Expected result")
	}
}

func TestAuthDecorator(t *testing.T) {
	handler := NewBasicHandler("SecureAPI")
	decorator := &AuthDecorator{}
	wrapped := decorator.Wrap(handler)

	// 无 token，应该失败
	result := wrapped.Handle("GET /admin")
	if result != "[ERROR] Unauthorized" {
		t.Errorf("Expected unauthorized error")
	}

	// 有 token，应该成功
	result = wrapped.Handle("GET /admin token=abc123")
	if result == "[ERROR] Unauthorized" {
		t.Errorf("Expected success with token")
	}
}

func TestDecoratorChain(t *testing.T) {
	handler := NewBasicHandler("API")
	
	chain := NewHandlerChain()
	chain.AddDecorator(&LoggingDecorator{})
	chain.AddDecorator(&TimingDecorator{})
	
	wrapped := chain.Wrap(handler)
	result := wrapped.Handle("GET /data")
	
	if result == "" {
		t.Error("Expected result")
	}
}

func TestStringDecorators(t *testing.T) {
	processor := &BaseProcessor{}
	
	// 大写装饰
	upper := &UpperCaseDecorator{}
	upperProc := upper.Decorate(processor)
	
	result := upperProc.Process("hello")
	if result != "HELLO" {
		t.Errorf("Expected 'HELLO', got '%s'", result)
	}
	
	// 反转装饰
	reverse := &ReverseDecorator{}
	reverseProc := reverse.Decorate(processor)
	
	result = reverseProc.Process("hello")
	if result != "olleh" {
		t.Errorf("Expected 'olleh', got '%s'", result)
	}
	
	// 组合装饰
	combined := reverse.Decorate(upperProc)
	result = combined.Process("hello")
	if result != "OLLEH" {
		t.Errorf("Expected 'OLLEH', got '%s'", result)
	}
}

func TestMemoryStream(t *testing.T) {
	stream := NewMemoryStream()
	
	err := stream.Write("Hello")
	if err != nil {
		t.Errorf("Write failed: %v", err)
	}
	
	err = stream.Write(" World")
	if err != nil {
		t.Errorf("Write failed: %v", err)
	}
	
	data := stream.Read()
	if data != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", data)
	}
	
	err = stream.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}
	
	// 关闭后写入应该失败
	err = stream.Write("test")
	if err == nil {
		t.Error("Expected error after close")
	}
}

func TestEncryptedStream(t *testing.T) {
	original := NewMemoryStream()
	encryptor := &EncryptedStream{}
	encrypted := encryptor.DecorateStream(original)
	
	err := encrypted.Write("secret")
	if err != nil {
		t.Errorf("Write failed: %v", err)
	}
	
	data := encrypted.Read()
	if data != "secret" {
		t.Errorf("Expected 'secret', got '%s'", data)
	}
}

func BenchmarkLoggingDecorator(b *testing.B) {
	handler := NewBasicHandler("API")
	decorator := &LoggingDecorator{}
	wrapped := decorator.Wrap(handler)
	
	for i := 0; i < b.N; i++ {
		_ = wrapped.Handle("GET /test")
	}
}

func BenchmarkStringDecorators(b *testing.B) {
	processor := &BaseProcessor{}
	upper := &UpperCaseDecorator{}
	reverse := &ReverseDecorator{}
	combined := reverse.Decorate(upper.Decorate(processor))
	
	for i := 0; i < b.N; i++ {
		_ = combined.Process("hello world")
	}
}
