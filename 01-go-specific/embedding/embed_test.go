package embedding

import (
	"fmt"
	"testing"
)

// TestPerson 测试基础 Person 类型
func TestPerson(t *testing.T) {
	p := Person{
		Name: "Alice",
		Age:  25,
	}

	if p.Name != "Alice" {
		t.Errorf("expected Name 'Alice', got '%s'", p.Name)
	}
	if p.Age != 25 {
		t.Errorf("expected Age 25, got %d", p.Age)
	}

	greeting := p.Greet()
	expected := "Hello, I'm Alice, age 25"
	if greeting != expected {
		t.Errorf("expected greeting '%s', got '%s'", expected, greeting)
	}
}

// TestPersonBirthday 测试生日方法
func TestPersonBirthday(t *testing.T) {
	p := Person{
		Name: "Bob",
		Age:  30,
	}

	p.Birthday()

	if p.Age != 31 {
		t.Errorf("expected Age 31 after birthday, got %d", p.Age)
	}
}

// TestEmployee 测试 Employee 嵌入 Person
func TestEmployee(t *testing.T) {
	emp := CreateEmployee("Charlie", 35, "EMP001", "Engineering", 8000)

	// 验证可以直接访问 Person 的字段
	if emp.Name != "Charlie" {
		t.Errorf("expected Name 'Charlie', got '%s'", emp.Name)
	}
	if emp.Age != 35 {
		t.Errorf("expected Age 35, got %d", emp.Age)
	}

	// 验证 Employee 特有字段
	if emp.EmployeeID != "EMP001" {
		t.Errorf("expected EmployeeID 'EMP001', got '%s'", emp.EmployeeID)
	}
	if emp.Department != "Engineering" {
		t.Errorf("expected Department 'Engineering', got '%s'", emp.Department)
	}

	// 验证可以调用 Person 的方法（方法提升）
	greeting := emp.Greet()
	if greeting == "" {
		t.Error("expected non-empty greeting")
	}
	// 验证 Employee 特有方法
	work := emp.Work()
	expected := "Charlie is working in Engineering department"
	if work != expected {
		t.Errorf("expected work '%s', got '%s'", expected, work)
	}

	// 验证年薪计算
	annual := emp.GetAnnualSalary()
	if annual != 96000 {
		t.Errorf("expected annual salary 96000, got %.2f", annual)
	}
}

// TestStudent 测试 Student 嵌入 Person
func TestStudent(t *testing.T) {
	student := CreateStudent("David", 20, "STU001", "Senior")

	// 验证嵌入的字段
	if student.Name != "David" {
		t.Errorf("expected Name 'David', got '%s'", student.Name)
	}

	// 验证学生特有方法
	study := student.Study()
	expected := "David is studying in grade Senior"
	if study != expected {
		t.Errorf("expected study '%s', got '%s'", expected, study)
	}

	// 测试注册课程
	student.Enroll("Math")
	student.Enroll("Physics")

	if len(student.Courses) != 2 {
		t.Errorf("expected 2 courses, got %d", len(student.Courses))
	}
}

// TestVerboseLogger 测试方法覆盖
func TestVerboseLogger(t *testing.T) {
	// 测试基础 Logger
	logger := &Logger{Prefix: "INFO"}
	logger.Log("test") // 使用 logger

	// 测试 VerboseLogger
	vLogger := &VerboseLogger{
		Logger:    Logger{Prefix: "DEBUG"},
		Timestamp: false,
	}

	// 应该调用覆盖后的方法（不 panic 即可）
	vLogger.Log("test message")

	// 测试带时间戳
	vLogger.Timestamp = true
	vLogger.Log("timed message")
}

// TestDuck 测试多重嵌入
func TestDuck(t *testing.T) {
	duck := &Duck{
		Name: "Donald",
		Flyer: Flyer{
			MaxAltitude: 100.5,
		},
		Swimmer: Swimmer{
			MaxDepth: 5.2,
		},
	}

	// 验证可以直接访问嵌入类型的字段
	if duck.Name != "Donald" {
		t.Errorf("expected Name 'Donald', got '%s'", duck.Name)
	}
	if duck.MaxAltitude != 100.5 {
		t.Errorf("expected MaxAltitude 100.5, got %.2f", duck.MaxAltitude)
	}
	if duck.MaxDepth != 5.2 {
		t.Errorf("expected MaxDepth 5.2, got %.2f", duck.MaxDepth)
	}

	// 验证可以调用嵌入类型的方法
	fly := duck.Fly()
	if fly == "" {
		t.Error("expected non-empty fly message")
	}

	swim := duck.Swim()
	if swim == "" {
		t.Error("expected non-empty swim message")
	}

	// 验证鸭子特有方法
	quack := duck.Quack()
	expected := "Donald says: Quack!"
	if quack != expected {
		t.Errorf("expected quack '%s', got '%s'", expected, quack)
	}
}

// TestLoggingHandler 测试 HTTP Handler 嵌入
func TestLoggingHandler(t *testing.T) {
	handler := &LoggingHandler{
		HTTPHandler: HTTPHandler{
			BasePath: "/api",
		},
		LogEnabled: true,
	}

	result := handler.Handle("/users")
	if result == "" {
		t.Error("expected non-empty result")
	}
}

// TestAuthHandler 测试认证处理器
func TestAuthHandler(t *testing.T) {
	handler := &AuthHandler{
		HTTPHandler: HTTPHandler{
			BasePath: "/admin",
		},
		RequireAuth:  true,
		AllowedUsers: []string{"admin", "user1"},
	}

	result := handler.Handle("/dashboard")
	if result == "" {
		t.Error("expected non-empty result")
	}
}

// TestMethodPromotion 测试方法提升
func TestMethodPromotion(t *testing.T) {
	emp := CreateEmployee("Eve", 28, "EMP002", "HR", 6000)

	// 验证 Birthday 方法被提升
	originalAge := emp.Age
	emp.Birthday()
	if emp.Age != originalAge+1 {
		t.Errorf("expected age to increase from %d to %d", originalAge, originalAge+1)
	}

	// 验证 Greet 方法被提升
	greeting := emp.Greet()
	if greeting == "" {
		t.Error("expected greeting from promoted method")
	}
}

// TestInterfaceEmbedding 测试接口嵌入
func TestInterfaceEmbedding(t *testing.T) {
	// 创建一个实现 ReadWriter 接口的类型
	rw := &mockReadWriter{}

	// 验证可以实现 Reader 接口
	var r Reader = rw
	if r.Read() != "data" {
		t.Error("expected Read to return 'data'")
	}

	// 验证可以实现 Writer 接口
	var w Writer = rw
	err := w.Write("test")
	if err != nil {
		t.Errorf("expected no error from Write, got %v", err)
	}

	// 验证可以实现 ReadWriter 接口
	var rwInterface ReadWriter = rw
	err = rwInterface.Close()
	if err != nil {
		t.Errorf("expected no error from Close, got %v", err)
	}
}

// mockReadWriter 模拟实现 ReadWriter 接口
type mockReadWriter struct{}

func (m *mockReadWriter) Read() string {
	return "data"
}

func (m *mockReadWriter) Write(data string) error {
	return nil
}

func (m *mockReadWriter) Close() error {
	return nil
}

// TestEmbeddingVsComposition 测试嵌入vs组合的区别
func TestEmbeddingVsComposition(t *testing.T) {
	// 嵌入：可以直接访问字段和方法
	emp := CreateEmployee("Frank", 40, "EMP003", "IT", 10000)

	// 可以直接访问 Name（通过方法提升）
	if emp.Name != "Frank" {
		t.Errorf("expected Name 'Frank', got '%s'", emp.Name)
	}

	// 可以直接调用 Greet（通过方法提升）
	greeting := emp.Greet()
	if greeting == "" {
		t.Error("expected non-empty greeting")
	}
}

// ExampleEmployee 演示 Employee 的使用
func ExampleEmployee() {
	emp := CreateEmployee("Alice", 30, "EMP001", "Engineering", 8000)

	fmt.Println("Name:", emp.Name)
	fmt.Println("Greeting:", emp.Greet())
	fmt.Println("Work:", emp.Work())
	fmt.Println("Annual Salary:", emp.GetAnnualSalary())
	// Output:
	// Name: Alice
	// Greeting: Hello, I'm Alice, age 30
	// Work: Alice is working in Engineering department
	// Annual Salary: 96000
}

// ExampleDuck 演示多重嵌入
func ExampleDuck() {
	duck := &Duck{
		Name: "Donald",
		Flyer: Flyer{
			MaxAltitude: 100,
		},
		Swimmer: Swimmer{
			MaxDepth: 5,
		},
	}

	fmt.Println(duck.Quack())
	fmt.Println(duck.Fly())
	fmt.Println(duck.Swim())
	// Output:
	// Donald says: Quack!
	// Flying at altitude 100.00 meters
	// Swimming at depth 5.00 meters
}
