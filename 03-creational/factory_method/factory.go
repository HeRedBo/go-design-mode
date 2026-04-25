// Package factorymethod 演示工厂方法模式 (Factory Method Pattern)
//
// 工厂方法模式解决的问题：
// 1. 对象创建逻辑复杂,需要封装创建过程
// 2. 客户端不需要知道具体类的名称
// 3. 支持扩展新的产品类型,符合开闭原则
//
// 与简单工厂的区别:
// - 简单工厂:一个工厂方法创建所有产品,违反开闭原则
// - 工厂方法:每个产品有独立的工厂方法,易于扩展
//
// 使用场景:
// - 创建相似但不同类型的对象
// - 需要延迟创建决策到子类
// - 希望提供扩展点让使用者自定义创建逻辑
//
// Go 中的实现特点:
// - Go 没有继承,使用接口和组合
// - 工厂方法通常返回接口类型
// - 利用函数式选项模式增强灵活性
//
// 核心组件:
// 1. Product 接口:定义产品通用行为
// 2. ConcreteProduct:具体产品实现
// 3. Factory 接口:定义工厂方法
// 4. ConcreteFactory:具体工厂实现
package factorymethod

import (
	"fmt"
	"math"
)

// ============ 产品接口和实现 ============

// Shape 形状接口（产品接口）
type Shape interface {
	Draw() string
	Area() float64
	Name() string
}

// Circle 圆形（具体产品）
type Circle struct {
	Radius float64
}

// Draw 实现 Shape 接口
func (c *Circle) Draw() string {
	return fmt.Sprintf("Drawing Circle with radius %.2f", c.Radius)
}

// Area 计算面积
func (c *Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Name 返回名称
func (c *Circle) Name() string {
	return "Circle"
}

// Rectangle 矩形（具体产品）
type Rectangle struct {
	Width, Height float64
}

// Draw 实现 Shape 接口
func (r *Rectangle) Draw() string {
	return fmt.Sprintf("Drawing Rectangle %.2fx%.2f", r.Width, r.Height)
}

// Area 计算面积
func (r *Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Name 返回名称
func (r *Rectangle) Name() string {
	return "Rectangle"
}

// Triangle 三角形（具体产品）
type Triangle struct {
	Base, Height float64
}

// Draw 实现 Shape 接口
func (t *Triangle) Draw() string {
	return fmt.Sprintf("Drawing Triangle base=%.2f, height=%.2f", t.Base, t.Height)
}

// Area 计算面积
func (t *Triangle) Area() float64 {
	return 0.5 * t.Base * t.Height
}

// Name 返回名称
func (t *Triangle) Name() string {
	return "Triangle"
}

// ============ 工厂接口和实现 ============

// ShapeFactory 形状工厂接口
type ShapeFactory interface {
	CreateShape() Shape
}

// CircleFactory 圆形工厂
type CircleFactory struct {
	Radius float64
}

// CreateShape 实现 ShapeFactory 接口
func (f *CircleFactory) CreateShape() Shape {
	return &Circle{Radius: f.Radius}
}

// RectangleFactory 矩形工厂
type RectangleFactory struct {
	Width, Height float64
}

// CreateShape 实现 ShapeFactory 接口
func (f *RectangleFactory) CreateShape() Shape {
	return &Rectangle{Width: f.Width, Height: f.Height}
}

// TriangleFactory 三角形工厂
type TriangleFactory struct {
	Base, Height float64
}

// CreateShape 实现 ShapeFactory 接口
func (f *TriangleFactory) CreateShape() Shape {
	return &Triangle{Base: f.Base, Height: f.Height}
}

// ============ 通用工厂函数（Go 风格） ============

// ShapeType 形状类型
type ShapeType string

const (
	ShapeCircle    ShapeType = "circle"
	ShapeRectangle ShapeType = "rectangle"
	ShapeTriangle  ShapeType = "triangle"
)

// ShapeConfig 形状配置
type ShapeConfig struct {
	Radius  float64 // 圆形半径
	Width   float64 // 矩形宽度
	Height  float64 // 矩形高度/三角形高度
	Base    float64 // 三角形底边
}

// CreateShape 通用工厂函数（根据类型创建）
// Go 语言更倾向于使用函数而非工厂类
func CreateShape(shapeType ShapeType, config ShapeConfig) (Shape, error) {
	switch shapeType {
	case ShapeCircle:
		if config.Radius <= 0 {
			return nil, fmt.Errorf("invalid radius: %f", config.Radius)
		}
		return &Circle{Radius: config.Radius}, nil

	case ShapeRectangle:
		if config.Width <= 0 || config.Height <= 0 {
			return nil, fmt.Errorf("invalid dimensions: %.2fx%.2f", config.Width, config.Height)
		}
		return &Rectangle{Width: config.Width, Height: config.Height}, nil

	case ShapeTriangle:
		if config.Base <= 0 || config.Height <= 0 {
			return nil, fmt.Errorf("invalid dimensions: base=%.2f, height=%.2f", config.Base, config.Height)
		}
		return &Triangle{Base: config.Base, Height: config.Height}, nil

	default:
		return nil, fmt.Errorf("unknown shape type: %s", shapeType)
	}
}

// ============ 实际应用示例 ============

// Document 文档接口
type Document interface {
	Open() string
	Save() string
	GetType() string
}

// PDFDocument PDF 文档
type PDFDocument struct {
	Filename string
}

// Open 打开文档
func (p *PDFDocument) Open() string {
	return fmt.Sprintf("Opening PDF document: %s", p.Filename)
}

// Save 保存文档
func (p *PDFDocument) Save() string {
	return fmt.Sprintf("Saving PDF document: %s", p.Filename)
}

// GetType 获取类型
func (p *PDFDocument) GetType() string {
	return "PDF"
}

// WordDocument Word 文档
type WordDocument struct {
	Filename string
}

// Open 打开文档
func (w *WordDocument) Open() string {
	return fmt.Sprintf("Opening Word document: %s", w.Filename)
}

// Save 保存文档
func (w *WordDocument) Save() string {
	return fmt.Sprintf("Saving Word document: %s", w.Filename)
}

// GetType 获取类型
func (w *WordDocument) GetType() string {
	return "Word"
}

// ExcelDocument Excel 文档
type ExcelDocument struct {
	Filename string
}

// Open 打开文档
func (e *ExcelDocument) Open() string {
	return fmt.Sprintf("Opening Excel document: %s", e.Filename)
}

// Save 保存文档
func (e *ExcelDocument) Save() string {
	return fmt.Sprintf("Saving Excel document: %s", e.Filename)
}

// GetType 获取类型
func (e *ExcelDocument) GetType() string {
	return "Excel"
}

// DocumentFactory 文档工厂接口
type DocumentFactory interface {
	CreateDocument(filename string) Document
}

// PDFDocumentFactory PDF 文档工厂
type PDFDocumentFactory struct{}

// CreateDocument 创建 PDF 文档
func (f *PDFDocumentFactory) CreateDocument(filename string) Document {
	return &PDFDocument{Filename: filename}
}

// WordDocumentFactory Word 文档工厂
type WordDocumentFactory struct{}

// CreateDocument 创建 Word 文档
func (f *WordDocumentFactory) CreateDocument(filename string) Document {
	return &WordDocument{Filename: filename}
}

// ExcelDocumentFactory Excel 文档工厂
type ExcelDocumentFactory struct{}

// CreateDocument 创建 Excel 文档
func (f *ExcelDocumentFactory) CreateDocument(filename string) Document {
	return &ExcelDocument{Filename: filename}
}

// DocumentType 文档类型
type DocumentType string

const (
	DocPDF   DocumentType = "pdf"
	DocWord  DocumentType = "word"
	DocExcel DocumentType = "excel"
)

// CreateDocument 通用文档工厂函数
func CreateDocument(docType DocumentType, filename string) (Document, error) {
	switch docType {
	case DocPDF:
		return &PDFDocument{Filename: filename}, nil
	case DocWord:
		return &WordDocument{Filename: filename}, nil
	case DocExcel:
		return &ExcelDocument{Filename: filename}, nil
	default:
		return nil, fmt.Errorf("unsupported document type: %s", docType)
	}
}
