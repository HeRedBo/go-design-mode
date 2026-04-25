package prototype

import (
	"testing"
)

// ============ Shape 原型测试 ============

func TestShapeClone(t *testing.T) {
	original := &Shape{
		Type:  "circle",
		Color: "red",
		X:     10.0,
		Y:     20.0,
		Metadata: map[string]string{
			"author": "test",
			"version": "1.0",
		},
	}

	// 克隆
	cloned := original.Clone().(*Shape)

	// 验证值相同
	if cloned.Type != original.Type {
		t.Errorf("Type mismatch: expected %s, got %s", original.Type, cloned.Type)
	}
	if cloned.Color != original.Color {
		t.Errorf("Color mismatch: expected %s, got %s", original.Color, cloned.Color)
	}
	if cloned.X != original.X {
		t.Errorf("X mismatch: expected %f, got %f", original.X, cloned.X)
	}
	if cloned.Y != original.Y {
		t.Errorf("Y mismatch: expected %f, got %f", original.Y, cloned.Y)
	}

	// 验证是深拷贝（修改克隆对象不影响原对象）
	cloned.Metadata["author"] = "modified"
	if original.Metadata["author"] == "modified" {
		t.Error("Shallow copy detected: modifying cloned metadata affected original")
	}

	// 验证是不同的对象
	if cloned == original {
		t.Error("Clone returned the same object, not a copy")
	}
}

func TestShapeName(t *testing.T) {
	shape := &Shape{
		Type:  "rectangle",
		Color: "blue",
		X:     5.5,
		Y:     10.5,
	}

	expected := "Shape(rectangle, color=blue, pos=5.5,10.5)"
	if shape.Name() != expected {
		t.Errorf("Expected name '%s', got '%s'", expected, shape.Name())
	}
}

// ============ ServerConfig 原型测试 ============

func TestServerConfigClone(t *testing.T) {
	original := &ServerConfig{
		ConfigName: "production",
		Host:       "localhost",
		Port:       8080,
		MaxConn:    1000,
		Timeout:    30,
		Features: map[string]bool{
			"logging": true,
			"metrics": true,
			"cache":   false,
		},
		Extensions: []string{"ext1", "ext2", "ext3"},
	}

	// 克隆
	cloned := original.Clone().(*ServerConfig)

	// 验证值相同
	if cloned.ConfigName != original.ConfigName {
		t.Errorf("ConfigName mismatch")
	}
	if cloned.Port != original.Port {
		t.Errorf("Port mismatch")
	}
	if cloned.MaxConn != original.MaxConn {
		t.Errorf("MaxConn mismatch")
	}

	// 验证深拷贝 - map
	cloned.Features["logging"] = false
	if original.Features["logging"] == false {
		t.Error("Shallow copy detected: modifying cloned Features affected original")
	}

	// 验证深拷贝 - slice
	cloned.Extensions[0] = "modified"
	if original.Extensions[0] == "modified" {
		t.Error("Shallow copy detected: modifying cloned Extensions affected original")
	}

	// 验证是不同的对象
	if cloned == original {
		t.Error("Clone returned the same object")
	}
}

func TestServerConfigMethods(t *testing.T) {
	config := &ServerConfig{
		ConfigName: "test",
		Host:       "localhost",
		Port:       3000,
		Features:   make(map[string]bool),
		Extensions: []string{},
	}

	// 测试启用特性
	config.EnableFeature("logging")
	if !config.Features["logging"] {
		t.Error("Expected logging feature to be enabled")
	}

	// 测试添加扩展
	config.AddExtension("ext1")
	if len(config.Extensions) != 1 || config.Extensions[0] != "ext1" {
		t.Error("Expected extension to be added")
	}
}

func TestServerConfigName(t *testing.T) {
	config := &ServerConfig{
		ConfigName: "staging",
		Host:       "staging.example.com",
		Port:       8443,
	}

	expected := "ServerConfig(staging: staging.example.com:8443)"
	if config.Name() != expected {
		t.Errorf("Expected name '%s', got '%s'", expected, config.Name())
	}
}

// ============ Document 原型测试 ============

func TestDocumentClone(t *testing.T) {
	original := &Document{
		Title:   "Test Document",
		Author:  "John Doe",
		Content: "Hello World",
		Tags:    []string{"test", "demo"},
		Properties: map[string]interface{}{
			"version": 1,
			"draft":   true,
		},
	}

	// 克隆
	cloned := original.Clone().(*Document)

	// 验证值相同
	if cloned.Title != original.Title {
		t.Errorf("Title mismatch")
	}
	if cloned.Author != original.Author {
		t.Errorf("Author mismatch")
	}
	if cloned.Content != original.Content {
		t.Errorf("Content mismatch")
	}

	// 验证是不同的对象
	if cloned == original {
		t.Error("Clone returned the same object")
	}
}

func TestDocumentAddTag(t *testing.T) {
	doc := &Document{
		Title: "Test",
		Tags:  []string{},
	}

	doc.AddTag("tag1")
	doc.AddTag("tag2")

	if len(doc.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(doc.Tags))
	}
}

func TestDocumentName(t *testing.T) {
	doc := &Document{
		Title:  "My Document",
		Author: "Jane",
	}

	expected := "Document(My Document by Jane)"
	if doc.Name() != expected {
		t.Errorf("Expected name '%s', got '%s'", expected, doc.Name())
	}
}

// ============ 原型注册表测试 ============

func TestPrototypeRegistry(t *testing.T) {
	registry := NewPrototypeRegistry()

	// 注册原型
	shapePrototype := &Shape{
		Type:  "circle",
		Color: "red",
		X:     0,
		Y:     0,
	}
	registry.Register("red_circle", shapePrototype)

	configPrototype := &ServerConfig{
		ConfigName: "default",
		Host:       "localhost",
		Port:       8080,
		Features:   make(map[string]bool),
		Extensions: []string{},
	}
	registry.Register("default_config", configPrototype)

	// 克隆原型
	clonedShape, err := registry.Clone("red_circle")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	clonedConfig, err := registry.Clone("default_config")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// 验证克隆
	if clonedShape.(*Shape).Type != "circle" {
		t.Error("Cloned shape type mismatch")
	}
	if clonedConfig.(*ServerConfig).Port != 8080 {
		t.Error("Cloned config port mismatch")
	}

	// 修改克隆对象不影响原型
	clonedShape.(*Shape).Color = "blue"
	if shapePrototype.Color == "blue" {
		t.Error("Modifying clone affected original prototype")
	}
}

func TestPrototypeRegistryList(t *testing.T) {
	registry := NewPrototypeRegistry()

	registry.Register("shape1", &Shape{Type: "circle"})
	registry.Register("shape2", &Shape{Type: "rectangle"})
	registry.Register("config1", &ServerConfig{ConfigName: "test"})

	names := registry.ListPrototypes()
	if len(names) != 3 {
		t.Errorf("Expected 3 prototypes, got %d", len(names))
	}
}

func TestPrototypeRegistryUnregister(t *testing.T) {
	registry := NewPrototypeRegistry()

	registry.Register("temp", &Shape{Type: "circle"})
	if !registry.HasPrototype("temp") {
		t.Error("Expected prototype to exist")
	}

	registry.Unregister("temp")
	if registry.HasPrototype("temp") {
		t.Error("Expected prototype to be removed")
	}
}

func TestPrototypeRegistryCloneNotFound(t *testing.T) {
	registry := NewPrototypeRegistry()

	_, err := registry.Clone("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent prototype")
	}
}

// ============ 泛型原型注册表测试 ============

func TestGenericPrototypeRegistry(t *testing.T) {
	registry := NewGenericPrototypeRegistry[*Shape]()

	registry.Register("circle", &Shape{
		Type:  "circle",
		Color: "green",
		X:     5.0,
		Y:     10.0,
	})

	// 克隆
	cloned, err := registry.Clone("circle")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if cloned.Type != "circle" {
		t.Errorf("Expected type 'circle', got '%s'", cloned.Type)
	}

	if cloned.Color != "green" {
		t.Errorf("Expected color 'green', got '%s'", cloned.Color)
	}

	// 验证是深拷贝
	cloned.Color = "blue"
	original, _ := registry.Clone("circle")
	if original.Color == "blue" {
		t.Error("Modifying clone affected original")
	}
}

// ============ 工具函数测试 ============

func TestDeepCloneShape(t *testing.T) {
	original := &Shape{
		Type:  "triangle",
		Color: "yellow",
		X:     1.0,
		Y:     2.0,
		Metadata: map[string]string{
			"key": "value",
		},
	}

	cloned := DeepCloneShape(original)

	if cloned.Type != original.Type {
		t.Errorf("Type mismatch")
	}

	// 验证深拷贝
	cloned.Metadata["key"] = "modified"
	if original.Metadata["key"] == "modified" {
		t.Error("Shallow copy detected")
	}
}

func TestDeepCloneShapeNil(t *testing.T) {
	cloned := DeepCloneShape(nil)
	if cloned != nil {
		t.Error("Expected nil for nil input")
	}
}

func TestDeepCloneConfigNil(t *testing.T) {
	cloned := DeepCloneConfig(nil)
	if cloned != nil {
		t.Error("Expected nil for nil input")
	}
}

// ============ 基准测试 ============

func BenchmarkShapeClone(b *testing.B) {
	shape := &Shape{
		Type:     "circle",
		Color:    "red",
		X:        10.0,
		Y:        20.0,
		Metadata: map[string]string{"key": "value"},
	}

	for i := 0; i < b.N; i++ {
		_ = shape.Clone()
	}
}

func BenchmarkServerConfigClone(b *testing.B) {
	config := &ServerConfig{
		ConfigName: "test",
		Host:       "localhost",
		Port:       8080,
		MaxConn:    1000,
		Timeout:    30,
		Features:   map[string]bool{"a": true, "b": true},
		Extensions: []string{"ext1", "ext2"},
	}

	for i := 0; i < b.N; i++ {
		_ = config.Clone()
	}
}

func BenchmarkDocumentClone(b *testing.B) {
	doc := &Document{
		Title:      "Test",
		Author:     "Author",
		Content:    "Content",
		Tags:       []string{"tag1", "tag2"},
		Properties: map[string]interface{}{"key": "value"},
	}

	for i := 0; i < b.N; i++ {
		_ = doc.Clone()
	}
}

func BenchmarkRegistryClone(b *testing.B) {
	registry := NewPrototypeRegistry()
	registry.Register("shape", &Shape{
		Type:  "circle",
		Color: "red",
	})

	for i := 0; i < b.N; i++ {
		_, _ = registry.Clone("shape")
	}
}

func BenchmarkDeepCloneShape(b *testing.B) {
	shape := &Shape{
		Type:     "circle",
		Color:    "red",
		X:        10.0,
		Y:        20.0,
		Metadata: map[string]string{"key": "value"},
	}

	for i := 0; i < b.N; i++ {
		_ = DeepCloneShape(shape)
	}
}
