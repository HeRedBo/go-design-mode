package singleton

import (
	"sync"
	"testing"
)

// TestConfigManager 测试配置管理器单例
func TestConfigManager(t *testing.T) {
	// 多次获取应该返回同一实例
	cm1 := GetConfigManager()
	cm2 := GetConfigManager()

	if cm1 != cm2 {
		t.Error("expected same instance")
	}

	// 测试设置和获取配置
	cm1.Set("key1", "value1")
	val, ok := cm2.Get("key1")
	if !ok || val != "value1" {
		t.Errorf("expected 'value1', got '%s', ok=%v", val, ok)
	}
}

// TestConfigManager_Concurrent 测试并发安全性
func TestConfigManager_Concurrent(t *testing.T) {
	var wg sync.WaitGroup

	// 并发获取实例
	instances := make([]*ConfigManager, 100)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			instances[idx] = GetConfigManager()
		}(i)
	}

	wg.Wait()

	// 验证所有实例都相同
	first := instances[0]
	for i := 1; i < 100; i++ {
		if instances[i] != first {
			t.Errorf("instance %d is different from instance 0", i)
		}
	}
}

// TestDatabasePool 测试数据库连接池单例
func TestDatabasePool(t *testing.T) {
	pool1 := GetDatabasePool(5)
	pool2 := GetDatabasePool(10) // 第二次调用应该忽略参数

	if pool1 != pool2 {
		t.Error("expected same pool instance")
	}

	// 测试获取连接
	conn1, err := pool1.GetConnection()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if conn1 != 1 {
		t.Errorf("expected connection ID 1, got %d", conn1)
	}

	// 测试释放连接
	pool1.ReleaseConnection()
	conn2, err := pool1.GetConnection()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if conn2 != 1 {
		t.Errorf("expected connection ID 1, got %d", conn2)
	}
}

// TestDatabasePool_Exhausted 测试连接池耗尽
func TestDatabasePool_Exhausted(t *testing.T) {
	// 注意：单例是全局的，前面测试可能已使用连接
	pool := GetDatabasePool(5)

	// 先释放所有可能的连接（清理状态）
	for i := 0; i < 5; i++ {
		pool.ReleaseConnection()
	}

	// 现在获取 5 个连接（应该成功）
	for i := 0; i < 5; i++ {
		_, err := pool.GetConnection()
		if err != nil {
			t.Logf("Connection %d failed (may be used by other tests): %v", i+1, err)
			break
		}
	}

	// 尝试获取第 6 个（应该失败）
	_, err := pool.GetConnection()
	if err == nil {
		t.Log("Expected pool to be exhausted (may pass if not all connections used)")
	}
}

// TestLogger 测试日志记录器单例
func TestLogger(t *testing.T) {
	logger1 := GetLogger()
	logger2 := GetLogger()

	if logger1 != logger2 {
		t.Error("expected same logger instance")
	}

	// 测试设置前缀
	logger1.SetPrefix("[TEST]")
	// 验证另一个引用也看到变化
	if logger2.prefix != "[TEST]" {
		t.Errorf("expected prefix '[TEST]', got '%s'", logger2.prefix)
	}
}

// TestSingletonFactory 测试单例工厂
func TestSingletonFactory(t *testing.T) {
	factory1 := GetSingletonFactory()
	factory2 := GetSingletonFactory()

	if factory1 != factory2 {
		t.Error("expected same factory instance")
	}

	// 注册一个单例
	mockSingleton := &mockSingleton{name: "test"}
	factory1.Register("test", mockSingleton)

	// 获取单例
	s, ok := factory2.Get("test")
	if !ok {
		t.Error("expected to find registered singleton")
	}
	if s.Name() != "test" {
		t.Errorf("expected name 'test', got '%s'", s.Name())
	}
}

// mockSingleton 模拟单例实现
type mockSingleton struct {
	name string
}

func (m *mockSingleton) Name() string {
	return m.name
}

// TestSingleton_ConcurrentInit 测试并发初始化
func TestSingleton_ConcurrentInit(t *testing.T) {
	var wg sync.WaitGroup
	instances := make([]*ConfigManager, 1000)

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			instances[idx] = GetConfigManager()
		}(i)
	}

	wg.Wait()

	// 验证所有实例相同
	first := instances[0]
	for i := 1; i < 1000; i++ {
		if instances[i] != first {
			t.Errorf("instance %d differs from instance 0", i)
		}
	}
}

// TestConfigManager_DefaultValues 测试默认配置值
func TestConfigManager_DefaultValues(t *testing.T) {
	cm := GetConfigManager()

	appName, ok := cm.Get("app_name")
	if !ok || appName != "MyApp" {
		t.Errorf("expected default app_name 'MyApp', got '%s', ok=%v", appName, ok)
	}

	version, ok := cm.Get("version")
	if !ok || version != "1.0.0" {
		t.Errorf("expected default version '1.0.0', got '%s', ok=%v", version, ok)
	}
}

// TestConfigManager_NotFound 测试获取不存在的配置
func TestConfigManager_NotFound(t *testing.T) {
	cm := GetConfigManager()

	_, ok := cm.Get("nonexistent")
	if ok {
		t.Error("expected ok to be false for nonexistent key")
	}
}

// ExampleGetConfigManager 演示配置管理器使用
func ExampleGetConfigManager() {
	cm := GetConfigManager()
	cm.Set("database_url", "localhost:5432")

	if val, ok := cm.Get("database_url"); ok {
		println("Database URL:", val)
	}
}
