// Package embedding 演示 Go 语言的嵌入组合模式 (Embedding Pattern)
//
// Go 没有传统的继承机制，而是使用组合和嵌入
// 嵌入的特点：
// 1. 匿名嵌入字段，可以直接访问嵌入类型的字段和方法
// 2. 方法提升（Method Promotion）：外部类型可以直接调用嵌入类型的方法
// 3. 可以重写（覆盖）嵌入类型的方法
// 4. 支持嵌入接口和结构体
//
// 与继承的区别：
// - 继承是"is-a"关系，嵌入是"has-a"关系
// - Go 偏向于组合而非继承
// - 嵌入更灵活，避免继承层次过深的问题
//
// 使用场景：
// - 代码复用
// - 扩展已有类型的功能
// - 实现接口组合
// - 构建装饰器/包装器
//
// 最佳实践：
// 1. 优先使用组合而非嵌入
// 2. 嵌入应该是"is implemented in terms of"关系
// 3. 避免过度嵌入导致代码混乱
// 4. 明确嵌入的目的是代码复用还是接口实现
package embedding

import (
	"fmt"
	"time"
)

// ============ 结构体嵌入 ============

// Person 基础类型
type Person struct {
	Name string
	Age  int
}

// Greet 打招呼方法
func (p *Person) Greet() string {
	return fmt.Sprintf("Hello, I'm %s, age %d", p.Name, p.Age)
}

// Birthday 生日方法
func (p *Person) Birthday() {
	p.Age++
}

// Employee 嵌入 Person（组合）
type Employee struct {
	Person      // 匿名嵌入，可以直接访问 Person 的字段和方法
	EmployeeID  string
	Department  string
	Salary      float64
}

// Work 员工特有的方法
func (e *Employee) Work() string {
	return fmt.Sprintf("%s is working in %s department", e.Name, e.Department)
}

// GetAnnualSalary 获取年薪
func (e *Employee) GetAnnualSalary() float64 {
	return e.Salary * 12
}

// Student 嵌入 Person
type Student struct {
	Person     // 匿名嵌入
	StudentID  string
	Grade      string
	Courses    []string
}

// Study 学生特有的方法
func (s *Student) Study() string {
	return fmt.Sprintf("%s is studying in grade %s", s.Name, s.Grade)
}

// Enroll 注册课程
func (s *Student) Enroll(course string) {
	s.Courses = append(s.Courses, course)
}

// ============ 接口嵌入 ============

// Reader 读取接口
type Reader interface {
	Read() string
}

// Writer 写入接口
type Writer interface {
	Write(data string) error
}

// ReadWriter 嵌入 Reader 和 Writer 接口
type ReadWriter interface {
	Reader
	Writer
	Close() error
}

// ============ 方法覆盖 ============

// Logger 基础日志器
type Logger struct {
	Prefix string
}

// Log 基础日志方法
func (l *Logger) Log(msg string) {
	fmt.Printf("[%s] %s\n", l.Prefix, msg)
}

// VerboseLogger 详细日志器，覆盖 Log 方法
type VerboseLogger struct {
	Logger        // 嵌入 Logger
	Timestamp     bool
}

// Log 覆盖父类的 Log 方法
func (vl *VerboseLogger) Log(msg string) {
	if vl.Timestamp {
		ts := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("[%s] [%s] %s\n", vl.Prefix, ts, msg)
	} else {
		vl.Logger.Log(msg) // 调用嵌入类型的方法
	}
}

// ============ 多重嵌入 ============

// Flyer 飞行能力
type Flyer struct {
	MaxAltitude float64
}

// Fly 飞行方法
func (f *Flyer) Fly() string {
	return fmt.Sprintf("Flying at altitude %.2f meters", f.MaxAltitude)
}

// Swimmer 游泳能力
type Swimmer struct {
	MaxDepth float64
}

// Swim 游泳方法
func (s *Swimmer) Swim() string {
	return fmt.Sprintf("Swimming at depth %.2f meters", s.MaxDepth)
}

// Duck 鸭子，嵌入多种能力
type Duck struct {
	Name    string
	Flyer   // 嵌入飞行能力
	Swimmer // 嵌入游泳能力
}

// Quack 鸭子特有的方法
func (d *Duck) Quack() string {
	return fmt.Sprintf("%s says: Quack!", d.Name)
}

// ============ 实际应用示例 ============

// HTTPHandler 基础 HTTP 处理器
type HTTPHandler struct {
	BasePath string
}

// Handle 处理请求
func (h *HTTPHandler) Handle(path string) string {
	return fmt.Sprintf("Handling %s%s", h.BasePath, path)
}

// LoggingHandler 带日志的处理器
type LoggingHandler struct {
	HTTPHandler // 嵌入基础处理器
	LogEnabled  bool
}

// Handle 重写方法，添加日志
func (lh *LoggingHandler) Handle(path string) string {
	if lh.LogEnabled {
		fmt.Printf("[LOG] Handling request: %s\n", path)
	}
	return lh.HTTPHandler.Handle(path) // 调用嵌入类型的方法
}

// AuthHandler 带认证的处理器
type AuthHandler struct {
	HTTPHandler      // 嵌入基础处理器
	RequireAuth      bool
	AllowedUsers     []string
}

// Handle 重写方法，添加认证
func (ah *AuthHandler) Handle(path string) string {
	if ah.RequireAuth {
		fmt.Printf("[AUTH] Checking authentication for: %s\n", path)
		// 认证逻辑...
	}
	return ah.HTTPHandler.Handle(path)
}

// ============ 工具函数 ============

// CreateEmployee 创建员工的便捷函数
func CreateEmployee(name string, age int, id, dept string, salary float64) *Employee {
	return &Employee{
		Person: Person{
			Name: name,
			Age:  age,
		},
		EmployeeID: id,
		Department: dept,
		Salary:     salary,
	}
}

// CreateStudent 创建学生的便捷函数
func CreateStudent(name string, age int, id, grade string) *Student {
	return &Student{
		Person: Person{
			Name: name,
			Age:  age,
		},
		StudentID: id,
		Grade:     grade,
		Courses:   make([]string, 0),
	}
}
