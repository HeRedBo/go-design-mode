// Package bridge 演示桥接模式 (Bridge Pattern)
//
// 桥接模式解决的问题：
// 1. 将抽象与实现分离，使它们可以独立变化
// 2. 避免类爆炸（多维度变化时的组合爆炸）
// 3. 运行时切换实现
//
// 与适配器的区别：
// - 适配器：事后补救，兼容已有接口
// - 桥接：事前设计，分离抽象和实现
//
// 使用场景：
// - UI 组件（抽象：按钮/文本框，实现：Windows/macOS/Linux）
// - 数据库驱动（抽象：ORM，实现：MySQL/PostgreSQL）
// - 消息发送（抽象：通知类型，实现：邮件/短信/推送）
//
// 核心组件：
// 1. Abstraction 抽象：高层接口
// 2. RefinedAbstraction 细化抽象：扩展抽象
// 3. Implementor 实现者：实现接口
// 4. ConcreteImplementor 具体实现：实际实现
package bridge

import "fmt"

// ============ 实现者接口 ============

// MessageImplementor 消息发送实现者接口
type MessageImplementor interface {
	Send(message string, to string) error
	GetType() string
}

// ============ 具体实现 ============

// EmailSender 邮件发送器
type EmailSender struct{}

// Send 发送邮件
func (e *EmailSender) Send(message string, to string) error {
	fmt.Printf("[Email] Sending to %s: %s\n", to, message)
	return nil
}

// GetType 返回类型
func (e *EmailSender) GetType() string {
	return "Email"
}

// SMSSender 短信发送器
type SMSSender struct{}

// Send 发送短信
func (s *SMSSender) Send(message string, to string) error {
	fmt.Printf("[SMS] Sending to %s: %s\n", to, message)
	return nil
}

// GetType 返回类型
func (s *SMSSender) GetType() string {
	return "SMS"
}

// PushSender 推送发送器
type PushSender struct{}

// Send 发送推送
func (p *PushSender) Send(message string, to string) error {
	fmt.Printf("[Push] Sending to %s: %s\n", to, message)
	return nil
}

// GetType 返回类型
func (p *PushSender) GetType() string {
	return "Push"
}

// ============ 抽象层 ============

// Notifier 通知抽象
type Notifier struct {
	sender MessageImplementor
}

// NewNotifier 创建通知器
func NewNotifier(sender MessageImplementor) *Notifier {
	return &Notifier{sender: sender}
}

// Send 发送通知
func (n *Notifier) Send(message string, to string) error {
	return n.sender.Send(message, to)
}

// GetType 返回类型
func (n *Notifier) GetType() string {
	return n.sender.GetType()
}

// ============ 细化抽象 ============

// AlertNotifier 告警通知器
type AlertNotifier struct {
	*Notifier
	priority string
}

// NewAlertNotifier 创建告警通知器
func NewAlertNotifier(sender MessageImplementor, priority string) *AlertNotifier {
	return &AlertNotifier{
		Notifier: NewNotifier(sender),
		priority: priority,
	}
}

// Send 发送告警（添加优先级前缀）
func (a *AlertNotifier) Send(message string, to string) error {
	formatted := fmt.Sprintf("[%s ALERT] %s", a.priority, message)
	return a.sender.Send(formatted, to)
}

// BatchNotifier 批量通知器
type BatchNotifier struct {
	*Notifier
}

// NewBatchNotifier 创建批量通知器
func NewBatchNotifier(sender MessageImplementor) *BatchNotifier {
	return &BatchNotifier{
		Notifier: NewNotifier(sender),
	}
}

// SendBatch 批量发送
func (b *BatchNotifier) SendBatch(message string, recipients []string) error {
	for _, to := range recipients {
		if err := b.sender.Send(message, to); err != nil {
			return err
		}
	}
	return nil
}

// ============ 示例 2: 形状绘制 ============

// Renderer 渲染器接口
type Renderer interface {
	RenderShape(shape string)
	GetType() string
}

// VectorRenderer 矢量渲染器
type VectorRenderer struct{}

// RenderShape 渲染矢量图形
func (v *VectorRenderer) RenderShape(shape string) {
	fmt.Printf("[Vector] Rendering %s with vector graphics\n", shape)
}

// GetType 返回类型
func (v *VectorRenderer) GetType() string {
	return "Vector"
}

// RasterRenderer 光栅渲染器
type RasterRenderer struct{}

// RenderShape 渲染光栅图形
func (r *RasterRenderer) RenderShape(shape string) {
	fmt.Printf("[Raster] Rendering %s with raster graphics\n", shape)
}

// GetType 返回类型
func (r *RasterRenderer) GetType() string {
	return "Raster"
}

// Shape 形状抽象
type Shape struct {
	renderer Renderer
	name     string
}

// NewShape 创建形状
func NewShape(renderer Renderer, name string) *Shape {
	return &Shape{
		renderer: renderer,
		name:     name,
	}
}

// Draw 绘制形状
func (s *Shape) Draw() {
	s.renderer.RenderShape(s.name)
}

// Circle 圆形
type Circle struct {
	*Shape
	radius float64
}

// NewCircle 创建圆形
func NewCircle(renderer Renderer, radius float64) *Circle {
	return &Circle{
		Shape:  NewShape(renderer, "Circle"),
		radius: radius,
	}
}

// Draw 绘制圆形
func (c *Circle) Draw() {
	fmt.Printf("Circle (radius: %.2f): ", c.radius)
	c.renderer.RenderShape(c.name)
}

// Rectangle 矩形
type Rectangle struct {
	*Shape
	width, height float64
}

// NewRectangle 创建矩形
func NewRectangle(renderer Renderer, width, height float64) *Rectangle {
	return &Rectangle{
		Shape:  NewShape(renderer, "Rectangle"),
		width:  width,
		height: height,
	}
}

// Draw 绘制矩形
func (r *Rectangle) Draw() {
	fmt.Printf("Rectangle (%.2fx%.2f): ", r.width, r.height)
	r.renderer.RenderShape(r.name)
}

// ============ 示例 3: 数据库驱动 ============

// DatabaseDriver 数据库驱动接口
type DatabaseDriver interface {
	Connect(host string, port int) error
	Query(sql string) ([]map[string]interface{}, error)
	Close() error
	GetType() string
}

// MySQLDriver MySQL 驱动
type MySQLDriver struct{}

// Connect 连接 MySQL
func (m *MySQLDriver) Connect(host string, port int) error {
	fmt.Printf("[MySQL] Connecting to %s:%d\n", host, port)
	return nil
}

// Query 执行查询
func (m *MySQLDriver) Query(sql string) ([]map[string]interface{}, error) {
	fmt.Printf("[MySQL] Executing: %s\n", sql)
	return []map[string]interface{}{{"id": 1, "name": "test"}}, nil
}

// Close 关闭连接
func (m *MySQLDriver) Close() error {
	fmt.Println("[MySQL] Closing connection")
	return nil
}

// GetType 返回类型
func (m *MySQLDriver) GetType() string {
	return "MySQL"
}

// PostgreSQLDriver PostgreSQL 驱动
type PostgreSQLDriver struct{}

// Connect 连接 PostgreSQL
func (p *PostgreSQLDriver) Connect(host string, port int) error {
	fmt.Printf("[PostgreSQL] Connecting to %s:%d\n", host, port)
	return nil
}

// Query 执行查询
func (p *PostgreSQLDriver) Query(sql string) ([]map[string]interface{}, error) {
	fmt.Printf("[PostgreSQL] Executing: %s\n", sql)
	return []map[string]interface{}{{"id": 1, "name": "test"}}, nil
}

// Close 关闭连接
func (p *PostgreSQLDriver) Close() error {
	fmt.Println("[PostgreSQL] Closing connection")
	return nil
}

// GetType 返回类型
func (p *PostgreSQLDriver) GetType() string {
	return "PostgreSQL"
}

// Database 数据库抽象
type Database struct {
	driver DatabaseDriver
	name   string
}

// NewDatabase 创建数据库
func NewDatabase(driver DatabaseDriver, name string) *Database {
	return &Database{
		driver: driver,
		name:   name,
	}
}

// Connect 连接数据库
func (d *Database) Connect(host string, port int) error {
	fmt.Printf("Database '%s' (%s): ", d.name, d.driver.GetType())
	return d.driver.Connect(host, port)
}

// Query 执行查询
func (d *Database) Query(sql string) ([]map[string]interface{}, error) {
	return d.driver.Query(sql)
}

// Close 关闭数据库
func (d *Database) Close() error {
	return d.driver.Close()
}
