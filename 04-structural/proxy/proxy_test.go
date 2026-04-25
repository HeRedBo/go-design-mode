package proxy

import "testing"

func TestImageProxy(t *testing.T) {
	proxy := NewImageProxy("photo.jpg")

	// 初始未加载
	if proxy.GetSize() != 0 {
		t.Error("Expected size 0 before loading")
	}

	// 首次访问，触发加载
	proxy.Display()

	// 验证已加载
	if proxy.GetSize() != 1024*1024 {
		t.Errorf("Expected size %d, got %d", 1024*1024, proxy.GetSize())
	}
}

func TestUserServiceProxy(t *testing.T) {
	service := NewRealUserService()

	// 未登录用户
	proxy := NewUserServiceProxy(service, nil)
	_, err := proxy.GetUser(1)
	if err == nil {
		t.Error("Expected authentication error")
	}

	// 普通用户
	bob := &User{ID: 2, Name: "Bob", Role: "user"}
	proxyBob := NewUserServiceProxy(service, bob)

	// 可以查看用户（密码被隐藏）
	user, err := proxyBob.GetUser(1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if user.Password != "***" {
		t.Error("Expected password to be hidden")
	}

	// 不能删除用户（权限不足）
	err = proxyBob.DeleteUser(1)
	if err == nil {
		t.Error("Expected permission error")
	}

	// 管理员
	admin := &User{ID: 1, Name: "Alice", Role: "admin"}
	proxyAdmin := NewUserServiceProxy(service, admin)

	// 可以删除用户
	err = proxyAdmin.DeleteUser(2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestCacheProxy(t *testing.T) {
	service := &RealDataService{}
	proxy := NewCacheProxy(service)

	// 首次访问（缓存未命中）
	data1, err := proxy.FetchData("key1")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// 第二次访问（缓存命中）
	data2, err := proxy.FetchData("key1")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if data1 != data2 {
		t.Error("Expected same data from cache")
	}

	// 验证缓存统计
	hits, misses := proxy.GetCacheStats()
	if hits != 1 || misses != 1 {
		t.Errorf("Expected 1 hit and 1 miss, got %d hits and %d misses", hits, misses)
	}
}

func TestLoggingProxy(t *testing.T) {
	calc := &RealCalculator{}
	proxy := NewLoggingProxy(calc)

	// 测试加法
	result := proxy.Add(10, 5)
	if result != 15 {
		t.Errorf("Expected 15, got %d", result)
	}

	// 测试除法
	result, err := proxy.Divide(10, 2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != 5 {
		t.Errorf("Expected 5, got %d", result)
	}

	// 测试除零错误
	_, err = proxy.Divide(10, 0)
	if err == nil {
		t.Error("Expected division by zero error")
	}
}

func BenchmarkImageProxy(b *testing.B) {
	proxy := NewImageProxy("test.jpg")

	for i := 0; i < b.N; i++ {
		proxy.Display()
	}
}

func BenchmarkCacheProxy(b *testing.B) {
	service := &RealDataService{}
	proxy := NewCacheProxy(service)

	// 预热缓存
	_, _ = proxy.FetchData("key1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = proxy.FetchData("key1")
	}
}

func BenchmarkLoggingProxy(b *testing.B) {
	calc := &RealCalculator{}
	proxy := NewLoggingProxy(calc)

	for i := 0; i < b.N; i++ {
		_ = proxy.Add(i, i+1)
	}
}
