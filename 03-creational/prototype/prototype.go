// Package prototype 演示原型模式 (Prototype Pattern)
//
// 原型模式解决的问题：
// 1. 通过复制现有对象来创建新对象，而不是从头创建
// 2. 避免重复初始化代码
// 3. 动态添加/删除产品结构
// 4. 创建复杂对象的成本较高时
//
// 使用场景：
// - 对象创建成本高（需要大量计算、数据库查询等）
// - 需要创建相似对象
// - 避免工厂类层次结构
// - 运行时动态创建对象类型
//
// Go 中的实现特点：
// - 使用 Clone() 方法实现原型接口
// - 深拷贝与浅拷贝的选择
// - 利用 struct 复制和反射
// - 注册表模式管理原型
//
// 核心组件：
// 1. Prototype 接口：定义 Clone 方法
// 2. ConcretePrototype：实现 Clone 方法的具体类
// 3. PrototypeRegistry：原型注册表（可选）
//
// 深拷贝 vs 浅拷贝：
// - 浅拷贝：只复制对象本身，引用类型字段指向相同对象
// - 深拷贝：递归复制所有引用类型字段，完全独立
// - Go 中通常使用深拷贝避免共享状态
package prototype

import (
	"encoding/json"
	"fmt"
	"sync"
)

// ============ 原型接口 ============

// Prototype 原型接口
type Prototype interface {
	Clone() Prototype
	Name() string
}

// ============ 产品：图形对象 ============

// Shape 图形原型
type Shape struct {
	Type     string
	Color    string
	X        float64
	Y        float64
	Metadata map[string]string
}

// Clone 克隆图形（深拷贝）
func (s *Shape) Clone() Prototype {
	// 深拷贝 map
	metadataCopy := make(map[string]string)
	for k, v := range s.Metadata {
		metadataCopy[k] = v
	}

	return &Shape{
		Type:     s.Type,
		Color:    s.Color,
		X:        s.X,
		Y:        s.Y,
		Metadata: metadataCopy,
	}
}

// Name 返回名称
func (s *Shape) Name() string {
	return fmt.Sprintf("Shape(%s, color=%s, pos=%.1f,%.1f)",
		s.Type, s.Color, s.X, s.Y)
}

// ============ 产品：配置对象 ============

// ServerConfig 服务器配置原型
type ServerConfig struct {
	ConfigName  string
	Host        string
	Port        int
	MaxConn     int
	Timeout     int
	Features    map[string]bool
	Extensions  []string
}

// Clone 克隆配置（深拷贝）
func (c *ServerConfig) Clone() Prototype {
	// 深拷贝 map
	featuresCopy := make(map[string]bool)
	for k, v := range c.Features {
		featuresCopy[k] = v
	}

	// 深拷贝 slice
	extensionsCopy := make([]string, len(c.Extensions))
	copy(extensionsCopy, c.Extensions)

	return &ServerConfig{
		ConfigName: c.ConfigName,
		Host:       c.Host,
		Port:       c.Port,
		MaxConn:    c.MaxConn,
		Timeout:    c.Timeout,
		Features:   featuresCopy,
		Extensions: extensionsCopy,
	}
}

// Name 返回名称
func (c *ServerConfig) Name() string {
	return fmt.Sprintf("ServerConfig(%s: %s:%d)", c.ConfigName, c.Host, c.Port)
}

// EnableFeature 启用特性
func (c *ServerConfig) EnableFeature(feature string) {
	c.Features[feature] = true
}

// AddExtension 添加扩展
func (c *ServerConfig) AddExtension(ext string) {
	c.Extensions = append(c.Extensions, ext)
}

// ============ 产品：文档对象 ============

// Document 文档原型
type Document struct {
	Title      string
	Author     string
	Content    string
	Tags       []string
	Properties map[string]interface{}
}

// Clone 克隆文档（使用 JSON 序列化实现深拷贝）
func (d *Document) Clone() Prototype {
	// 使用 JSON 序列化进行深拷贝
	data, _ := json.Marshal(d)
	var clone Document
	json.Unmarshal(data, &clone)

	// 恢复 interface{} 类型
	clone.Properties = make(map[string]interface{})
	for k, v := range d.Properties {
		clone.Properties[k] = v
	}

	return &clone
}

// Name 返回名称
func (d *Document) Name() string {
	return fmt.Sprintf("Document(%s by %s)", d.Title, d.Author)
}

// AddTag 添加标签
func (d *Document) AddTag(tag string) {
	d.Tags = append(d.Tags, tag)
}

// ============ 原型注册表 ============

// PrototypeRegistry 原型注册表
type PrototypeRegistry struct {
	prototypes map[string]Prototype
	mu         sync.RWMutex
}

// NewPrototypeRegistry 创建原型注册表
func NewPrototypeRegistry() *PrototypeRegistry {
	return &PrototypeRegistry{
		prototypes: make(map[string]Prototype),
	}
}

// Register 注册原型
func (r *PrototypeRegistry) Register(name string, prototype Prototype) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.prototypes[name] = prototype
}

// Unregister 注销原型
func (r *PrototypeRegistry) Unregister(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.prototypes, name)
}

// Clone 克隆已注册的原型
func (r *PrototypeRegistry) Clone(name string) (Prototype, error) {
	r.mu.RLock()
	prototype, exists := r.prototypes[name]
	r.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("prototype not found: %s", name)
	}

	return prototype.Clone(), nil
}

// ListPrototypes 列出所有注册的原型
func (r *PrototypeRegistry) ListPrototypes() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.prototypes))
	for name := range r.prototypes {
		names = append(names, name)
	}
	return names
}

// HasPrototype 检查原型是否已注册
func (r *PrototypeRegistry) HasPrototype(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.prototypes[name]
	return exists
}

// ============ 泛型原型注册表（Go 1.18+） ============

// GenericPrototypeRegistry 泛型原型注册表
type GenericPrototypeRegistry[T Prototype] struct {
	prototypes map[string]T
	mu         sync.RWMutex
}

// NewGenericPrototypeRegistry 创建泛型原型注册表
func NewGenericPrototypeRegistry[T Prototype]() *GenericPrototypeRegistry[T] {
	return &GenericPrototypeRegistry[T]{
		prototypes: make(map[string]T),
	}
}

// Register 注册原型
func (r *GenericPrototypeRegistry[T]) Register(name string, prototype T) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.prototypes[name] = prototype
}

// Clone 克隆已注册的原型
func (r *GenericPrototypeRegistry[T]) Clone(name string) (T, error) {
	r.mu.RLock()
	prototype, exists := r.prototypes[name]
	r.mu.RUnlock()

	var zero T
	if !exists {
		return zero, fmt.Errorf("prototype not found: %s", name)
	}

	cloned := prototype.Clone()
	if result, ok := cloned.(T); ok {
		return result, nil
	}
	return zero, fmt.Errorf("type mismatch for prototype: %s", name)
}

// ============ 实用工具函数 ============

// DeepCloneShape 深拷贝 Shape
func DeepCloneShape(s *Shape) *Shape {
	if s == nil {
		return nil
	}

	metadataCopy := make(map[string]string)
	for k, v := range s.Metadata {
		metadataCopy[k] = v
	}

	return &Shape{
		Type:     s.Type,
		Color:    s.Color,
		X:        s.X,
		Y:        s.Y,
		Metadata: metadataCopy,
	}
}

// DeepCloneConfig 深拷贝 ServerConfig
func DeepCloneConfig(c *ServerConfig) *ServerConfig {
	if c == nil {
		return nil
	}

	featuresCopy := make(map[string]bool)
	for k, v := range c.Features {
		featuresCopy[k] = v
	}

	extensionsCopy := make([]string, len(c.Extensions))
	copy(extensionsCopy, c.Extensions)

	return &ServerConfig{
		ConfigName: c.ConfigName,
		Host:       c.Host,
		Port:       c.Port,
		MaxConn:    c.MaxConn,
		Timeout:    c.Timeout,
		Features:   featuresCopy,
		Extensions: extensionsCopy,
	}
}

// CloneViaJSON 通过 JSON 序列化实现深拷贝（通用方法）
func CloneViaJSON(src interface{}, dst interface{}) error {
	data, err := json.Marshal(src)
	if err != nil {
		return fmt.Errorf("marshal failed: %w", err)
	}

	if err := json.Unmarshal(data, dst); err != nil {
		return fmt.Errorf("unmarshal failed: %w", err)
	}

	return nil
}
