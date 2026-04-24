// Package singleton 演示单例模式 (Singleton Pattern)
//
// 单例模式确保一个类只有一个实例，并提供全局访问点
// Go 中的实现方式：
// 1. 使用 sync.Once 保证线程安全（推荐）
// 2. 使用 init 函数（不推荐，无法延迟初始化）
// 3. 使用包级变量（简单但不安全）
//
// 使用场景：
// - 配置管理器
// - 数据库连接池
// - 日志记录器
// - 缓存管理器
// - 设备驱动
//
// sync.Once 的优势：
// 1. 线程安全
// 2. 延迟初始化（首次使用时创建）
// 3. 高性能（ Once 内部使用原子操作）
// 4. 简洁易用
package singleton

import (
	"fmt"
	"sync"
)

// ============ 方式 1: sync.Once 实现（推荐） ============

// ConfigManager 配置管理器单例
type ConfigManager struct {
	configs map[string]string
	mu      sync.RWMutex
}

var (
	configInstance *ConfigManager
	configOnce     sync.Once
)

// GetConfigManager 获取配置管理器单例（线程安全）
func GetConfigManager() *ConfigManager {
	configOnce.Do(func() {
		configInstance = &ConfigManager{
			configs: make(map[string]string),
		}
		// 初始化配置
		configInstance.configs["app_name"] = "MyApp"
		configInstance.configs["version"] = "1.0.0"
	})
	return configInstance
}

// Set 设置配置
func (cm *ConfigManager) Set(key, value string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.configs[key] = value
}

// Get 获取配置
func (cm *ConfigManager) Get(key string) (string, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	val, ok := cm.configs[key]
	return val, ok
}

// ============ 方式 2: 带初始化的单例 ============

// DatabasePool 数据库连接池单例
type DatabasePool struct {
	maxConnections int
	currentConnections int
	mu sync.Mutex
}

var (
	dbPoolInstance *DatabasePool
	dbPoolOnce     sync.Once
)

// GetDatabasePool 获取数据库连接池（带初始化参数）
func GetDatabasePool(maxConn int) *DatabasePool {
	dbPoolOnce.Do(func() {
		dbPoolInstance = &DatabasePool{
			maxConnections: maxConn,
		}
		fmt.Printf("DatabasePool initialized with max connections: %d\n", maxConn)
	})
	return dbPoolInstance
}

// GetConnection 获取连接
func (dp *DatabasePool) GetConnection() (int, error) {
	dp.mu.Lock()
	defer dp.mu.Unlock()

	if dp.currentConnections >= dp.maxConnections {
		return 0, fmt.Errorf("connection pool exhausted")
	}

	dp.currentConnections++
	return dp.currentConnections, nil
}

// ReleaseConnection 释放连接
func (dp *DatabasePool) ReleaseConnection() {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	if dp.currentConnections > 0 {
		dp.currentConnections--
	}
}

// ============ 方式 3: 通用单例工厂 ============

// Singleton 通用单例接口
type Singleton interface {
	Name() string
}

// SingletonFactory 单例工厂
type SingletonFactory struct {
	instances map[string]Singleton
	mu        sync.RWMutex
	once      sync.Once
}

var (
	factoryInstance *SingletonFactory
	factoryOnce     sync.Once
)

// GetSingletonFactory 获取单例工厂
func GetSingletonFactory() *SingletonFactory {
	factoryOnce.Do(func() {
		factoryInstance = &SingletonFactory{
			instances: make(map[string]Singleton),
		}
	})
	return factoryInstance
}

// Register 注册单例
func (sf *SingletonFactory) Register(name string, s Singleton) {
	sf.mu.Lock()
	defer sf.mu.Unlock()
	sf.instances[name] = s
}

// Get 获取单例
func (sf *SingletonFactory) Get(name string) (Singleton, bool) {
	sf.mu.RLock()
	defer sf.mu.RUnlock()
	s, ok := sf.instances[name]
	return s, ok
}

// ============ 示例：日志记录器单例 ============

// Logger 日志记录器单例
type Logger struct {
	prefix string
	mu     sync.Mutex
}

var (
	loggerInstance *Logger
	loggerOnce     sync.Once
)

// GetLogger 获取日志记录器单例
func GetLogger() *Logger {
	loggerOnce.Do(func() {
		loggerInstance = &Logger{
			prefix: "[LOG]",
		}
	})
	return loggerInstance
}

// Log 记录日志
func (l *Logger) Log(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Printf("%s %s\n", l.prefix, msg)
}

// SetPrefix 设置前缀
func (l *Logger) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prefix = prefix
}
