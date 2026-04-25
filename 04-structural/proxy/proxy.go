// Package proxy 演示代理模式 (Proxy Pattern)
//
// 代理模式解决的问题：
// 1. 控制对对象的访问
// 2. 在访问前后添加额外逻辑
// 3. 延迟初始化（虚拟代理）
// 4. 远程访问（远程代理）
// 5. 权限控制（保护代理）
//
// 使用场景：
// - 延迟加载大对象
// - 访问权限控制
// - 缓存代理
// - 日志记录
// - 远程调用
//
// 代理类型：
// 1. 虚拟代理：延迟初始化
// 2. 保护代理：权限控制
// 3. 缓存代理：缓存结果
// 4. 远程代理：远程对象本地代表
//
// 与装饰器的区别：
// - 装饰器：为对象添加行为
// - 代理：控制对对象的访问
package proxy

import (
	"fmt"
	"sync"
	"time"
)

// ============ 示例 1: 图像虚拟代理 ============

// Image 图像接口
type Image interface {
	Display()
	GetSize() int
}

// RealImage 真实图像
type RealImage struct {
	filename string
	size     int
	loaded   bool
}

// NewRealImage 创建真实图像
func NewRealImage(filename string) *RealImage {
	fmt.Printf("[RealImage] Loading %s from disk...\n", filename)
	// 模拟加载大文件
	time.Sleep(100 * time.Millisecond)
	return &RealImage{
		filename: filename,
		size:     1024 * 1024, // 1MB
		loaded:   true,
	}
}

// Display 显示图像
func (r *RealImage) Display() {
	fmt.Printf("[RealImage] Displaying %s (%d bytes)\n", r.filename, r.size)
}

// GetSize 获取大小
func (r *RealImage) GetSize() int {
	return r.size
}

// ImageProxy 图像代理（虚拟代理）
type ImageProxy struct {
	filename  string
	realImage *RealImage
	size      int
}

// NewImageProxy 创建图像代理
func NewImageProxy(filename string) *ImageProxy {
	return &ImageProxy{filename: filename}
}

// Display 显示图像（延迟加载）
func (p *ImageProxy) Display() {
	if p.realImage == nil {
		p.realImage = NewRealImage(p.filename)
		p.size = p.realImage.GetSize()
	}
	p.realImage.Display()
}

// GetSize 获取大小
func (p *ImageProxy) GetSize() int {
	if p.realImage == nil {
		return 0 // 未加载
	}
	return p.size
}

// ============ 示例 2: 用户服务保护代理 ============

// User 用户
type User struct {
	ID       int
	Name     string
	Role     string
	Password string
}

// UserService 用户服务接口
type UserService interface {
	GetUser(id int) (*User, error)
	DeleteUser(id int) error
	UpdateUser(user *User) error
}

// RealUserService 真实用户服务
type RealUserService struct {
	users map[int]*User
}

// NewRealUserService 创建用户服务
func NewRealUserService() *RealUserService {
	return &RealUserService{
		users: map[int]*User{
			1: {ID: 1, Name: "Alice", Role: "admin", Password: "secret"},
			2: {ID: 2, Name: "Bob", Role: "user", Password: "password"},
		},
	}
}

// GetUser 获取用户
func (s *RealUserService) GetUser(id int) (*User, error) {
	user, exists := s.users[id]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// DeleteUser 删除用户
func (s *RealUserService) DeleteUser(id int) error {
	if _, exists := s.users[id]; !exists {
		return fmt.Errorf("user not found")
	}
	delete(s.users, id)
	return nil
}

// UpdateUser 更新用户
func (s *RealUserService) UpdateUser(user *User) error {
	s.users[user.ID] = user
	return nil
}

// UserServiceProxy 用户服务代理（保护代理）
type UserServiceProxy struct {
	service    *RealUserService
	currentUser *User
}

// NewUserServiceProxy 创建用户服务代理
func NewUserServiceProxy(service *RealUserService, currentUser *User) *UserServiceProxy {
	return &UserServiceProxy{
		service:    service,
		currentUser: currentUser,
	}
}

// GetUser 获取用户（需要登录）
func (p *UserServiceProxy) GetUser(id int) (*User, error) {
	if p.currentUser == nil {
		return nil, fmt.Errorf("authentication required")
	}

	user, err := p.service.GetUser(id)
	if err != nil {
		return nil, err
	}

	// 隐藏密码
	safeUser := *user
	safeUser.Password = "***"
	return &safeUser, nil
}

// DeleteUser 删除用户（仅管理员）
func (p *UserServiceProxy) DeleteUser(id int) error {
	if p.currentUser == nil {
		return fmt.Errorf("authentication required")
	}

	if p.currentUser.Role != "admin" {
		return fmt.Errorf("permission denied: admin role required")
	}

	return p.service.DeleteUser(id)
}

// UpdateUser 更新用户（仅本人或管理员）
func (p *UserServiceProxy) UpdateUser(user *User) error {
	if p.currentUser == nil {
		return fmt.Errorf("authentication required")
	}

	if p.currentUser.Role != "admin" && p.currentUser.ID != user.ID {
		return fmt.Errorf("permission denied: can only update own profile")
	}

	return p.service.UpdateUser(user)
}

// ============ 示例 3: 缓存代理 ============

// DataService 数据服务接口
type DataService interface {
	FetchData(key string) (string, error)
}

// RealDataService 真实数据服务
type RealDataService struct{}

// FetchData 获取数据（模拟慢查询）
func (s *RealDataService) FetchData(key string) (string, error) {
	fmt.Printf("[RealDataService] Fetching data for key: %s\n", key)
	time.Sleep(200 * time.Millisecond) // 模拟慢查询
	return fmt.Sprintf("Data for %s", key), nil
}

// CacheProxy 缓存代理
type CacheProxy struct {
	service DataService
	cache   map[string]string
	mu      sync.RWMutex
	hits    int
	misses  int
}

// NewCacheProxy 创建缓存代理
func NewCacheProxy(service DataService) *CacheProxy {
	return &CacheProxy{
		service: service,
		cache:   make(map[string]string),
	}
}

// FetchData 获取数据（带缓存）
func (p *CacheProxy) FetchData(key string) (string, error) {
	// 先查缓存
	p.mu.RLock()
	data, exists := p.cache[key]
	p.mu.RUnlock()

	if exists {
		p.hits++
		fmt.Printf("[CacheProxy] Cache hit for key: %s\n", key)
		return data, nil
	}

	// 缓存未命中，从服务获取
	p.misses++
	data, err := p.service.FetchData(key)
	if err != nil {
		return "", err
	}

	// 存入缓存
	p.mu.Lock()
	p.cache[key] = data
	p.mu.Unlock()

	return data, nil
}

// GetCacheStats 获取缓存统计
func (p *CacheProxy) GetCacheStats() (hits, misses int) {
	return p.hits, p.misses
}

// GetCacheSize 获取缓存大小
func (p *CacheProxy) GetCacheSize() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.cache)
}

// ============ 示例 4: 日志代理 ============

// Calculator 计算器接口
type Calculator interface {
	Add(a, b int) int
	Subtract(a, b int) int
	Multiply(a, b int) int
	Divide(a, b int) (int, error)
}

// RealCalculator 真实计算器
type RealCalculator struct{}

// Add 加法
func (c *RealCalculator) Add(a, b int) int {
	return a + b
}

// Subtract 减法
func (c *RealCalculator) Subtract(a, b int) int {
	return a - b
}

// Multiply 乘法
func (c *RealCalculator) Multiply(a, b int) int {
	return a * b
}

// Divide 除法
func (c *RealCalculator) Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}

// LoggingProxy 日志代理
type LoggingProxy struct {
	calculator Calculator
}

// NewLoggingProxy 创建日志代理
func NewLoggingProxy(calculator Calculator) *LoggingProxy {
	return &LoggingProxy{calculator: calculator}
}

// Add 加法（带日志）
func (p *LoggingProxy) Add(a, b int) int {
	fmt.Printf("[LOG] Add(%d, %d)\n", a, b)
	result := p.calculator.Add(a, b)
	fmt.Printf("[LOG] Result: %d\n", result)
	return result
}

// Subtract 减法（带日志）
func (p *LoggingProxy) Subtract(a, b int) int {
	fmt.Printf("[LOG] Subtract(%d, %d)\n", a, b)
	result := p.calculator.Subtract(a, b)
	fmt.Printf("[LOG] Result: %d\n", result)
	return result
}

// Multiply 乘法（带日志）
func (p *LoggingProxy) Multiply(a, b int) int {
	fmt.Printf("[LOG] Multiply(%d, %d)\n", a, b)
	result := p.calculator.Multiply(a, b)
	fmt.Printf("[LOG] Result: %d\n", result)
	return result
}

// Divide 除法（带日志）
func (p *LoggingProxy) Divide(a, b int) (int, error) {
	fmt.Printf("[LOG] Divide(%d, %d)\n", a, b)
	result, err := p.calculator.Divide(a, b)
	if err != nil {
		fmt.Printf("[LOG] Error: %v\n", err)
		return 0, err
	}
	fmt.Printf("[LOG] Result: %d\n", result)
	return result, nil
}
