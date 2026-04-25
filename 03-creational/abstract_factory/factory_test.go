package abstractfactory

import (
	"testing"
)

// ============ Windows UI 工厂测试 ============

func TestWindowsFactory(t *testing.T) {
	factory := &WindowsFactory{}

	// 测试创建按钮
	button := factory.CreateButton()
	if button == nil {
		t.Fatal("Expected button, got nil")
	}

	rendered := button.Render()
	if rendered != "Rendering Windows-style button" {
		t.Errorf("Unexpected render: %s", rendered)
	}

	clicked := button.Click()
	if clicked != "Windows button clicked with ripple effect" {
		t.Errorf("Unexpected click: %s", clicked)
	}

	// 测试创建文本框
	textBox := factory.CreateTextBox()
	if textBox == nil {
		t.Fatal("Expected text box, got nil")
	}

	textBox.SetText("Hello")
	if textBox.GetText() != "Hello" {
		t.Errorf("Expected text 'Hello', got '%s'", textBox.GetText())
	}

	// 测试创建菜单
	menu := factory.CreateMenu()
	if menu == nil {
		t.Fatal("Expected menu, got nil")
	}

	menu.AddItem("File")
	menu.AddItem("Edit")
}

// ============ macOS UI 工厂测试 ============

func TestMacFactory(t *testing.T) {
	factory := &MacFactory{}

	button := factory.CreateButton()
	if button.Render() != "Rendering macOS button with rounded corners" {
		t.Errorf("Unexpected button render")
	}

	textBox := factory.CreateTextBox()
	textBox.SetText("Test")
	if textBox.GetText() != "Test" {
		t.Errorf("Expected text 'Test'")
	}

	menu := factory.CreateMenu()
	if menu.Render() != "Rendering macOS menu with global menu bar" {
		t.Errorf("Unexpected menu render")
	}
}

// ============ Linux UI 工厂测试 ============

func TestLinuxFactory(t *testing.T) {
	factory := &LinuxFactory{}

	button := factory.CreateButton()
	if button.Render() != "Rendering Linux-style button with GTK theme" {
		t.Errorf("Unexpected button render")
	}

	textBox := factory.CreateTextBox()
	if textBox.Render() != "Rendering Linux text box with GTK styling" {
		t.Errorf("Unexpected text box render")
	}

	menu := factory.CreateMenu()
	if menu.Render() != "Rendering Linux menu with GTK theme" {
		t.Errorf("Unexpected menu render")
	}
}

// ============ UI 工厂注册表测试 ============

func TestCreateUIFactory(t *testing.T) {
	// 测试 Windows 工厂
	winFactory, err := CreateUIFactory(Windows)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if winFactory == nil {
		t.Error("Expected Windows factory")
	}

	// 测试 macOS 工厂
	macFactory, err := CreateUIFactory(MacOS)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	button := macFactory.CreateButton()
	if button.Render() == "" {
		t.Error("Expected button render")
	}

	// 测试 Linux 工厂
	linuxFactory, err := CreateUIFactory(Linux)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if linuxFactory.CreateMenu() == nil {
		t.Error("Expected menu")
	}

	// 测试无效类型
	_, err = CreateUIFactory("invalid")
	if err == nil {
		t.Error("Expected error for invalid factory type")
	}
}

// ============ 数据库工厂测试 ============

func TestMySQLFactory(t *testing.T) {
	factory := &MySQLFactory{}

	conn := factory.CreateConnection()
	if conn.Connect() != "Connecting to MySQL database" {
		t.Errorf("Unexpected connection message")
	}

	cmd := factory.CreateCommand()
	result := cmd.Execute("SELECT * FROM users")
	if result != "Executing MySQL query: SELECT * FROM users" {
		t.Errorf("Unexpected query result: %s", result)
	}

	adapter := factory.CreateAdapter()
	if adapter.Adapt() != "Using MySQL adapter with native protocol" {
		t.Errorf("Unexpected adapter message")
	}
}

func TestPostgreSQLFactory(t *testing.T) {
	factory := &PostgreSQLFactory{}

	conn := factory.CreateConnection()
	if conn.Connect() != "Connecting to PostgreSQL database" {
		t.Errorf("Unexpected connection message")
	}

	cmd := factory.CreateCommand()
	result := cmd.Execute("SELECT * FROM products")
	if result != "Executing PostgreSQL query: SELECT * FROM products" {
		t.Errorf("Unexpected query result: %s", result)
	}

	adapter := factory.CreateAdapter()
	if adapter.Adapt() != "Using PostgreSQL adapter with libpq" {
		t.Errorf("Unexpected adapter message")
	}
}

func TestCreateDatabaseFactory(t *testing.T) {
	// 测试 MySQL 工厂
	mysqlFactory, err := CreateDatabaseFactory(MySQL)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	conn := mysqlFactory.CreateConnection()
	if conn.Close() != "Closing MySQL connection" {
		t.Errorf("Unexpected close message")
	}

	// 测试 PostgreSQL 工厂
	psqlFactory, err := CreateDatabaseFactory(PostgreSQL)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	conn = psqlFactory.CreateConnection()
	if conn.Close() != "Closing PostgreSQL connection" {
		t.Errorf("Unexpected close message")
	}

	// 测试无效类型
	_, err = CreateDatabaseFactory("oracle")
	if err == nil {
		t.Error("Expected error for unsupported database type")
	}
}

// ============ 产品族兼容性测试 ============

func TestUIFactoryProductCompatibility(t *testing.T) {
	// 确保同一工厂创建的产品属于同一产品族
	factories := map[string]UIFactory{
		"Windows": &WindowsFactory{},
		"macOS":   &MacFactory{},
		"Linux":   &LinuxFactory{},
	}

	for name, factory := range factories {
		button := factory.CreateButton()
		textBox := factory.CreateTextBox()
		menu := factory.CreateMenu()

		// 验证所有产品都能正常工作
		if button.Render() == "" {
			t.Errorf("%s: button render failed", name)
		}
		if textBox.Render() == "" {
			t.Errorf("%s: text box render failed", name)
		}
		if menu.Render() == "" {
			t.Errorf("%s: menu render failed", name)
		}
	}
}

func TestDatabaseFactoryProductCompatibility(t *testing.T) {
	// 确保同一工厂创建的产品属于同一产品族
	factories := map[string]DatabaseFactory{
		"MySQL":      &MySQLFactory{},
		"PostgreSQL": &PostgreSQLFactory{},
	}

	for name, factory := range factories {
		conn := factory.CreateConnection()
		cmd := factory.CreateCommand()
		adapter := factory.CreateAdapter()

		// 验证所有产品都能正常工作
		if conn.Connect() == "" {
			t.Errorf("%s: connection failed", name)
		}
		if cmd.Execute("SELECT 1") == "" {
			t.Errorf("%s: command failed", name)
		}
		if adapter.Adapt() == "" {
			t.Errorf("%s: adapter failed", name)
		}
	}
}

// ============ 多态测试 ============

func TestUIFactoryPolymorphism(t *testing.T) {
	// 使用接口实现多态
	var factories []UIFactory
	factories = append(factories, &WindowsFactory{})
	factories = append(factories, &MacFactory{})
	factories = append(factories, &LinuxFactory{})

	for _, factory := range factories {
		button := factory.CreateButton()
		if button == nil {
			t.Error("Expected button")
		}

		textBox := factory.CreateTextBox()
		if textBox == nil {
			t.Error("Expected text box")
		}
	}
}

func TestDatabaseFactoryPolymorphism(t *testing.T) {
	// 使用接口实现多态
	var factories []DatabaseFactory
	factories = append(factories, &MySQLFactory{})
	factories = append(factories, &PostgreSQLFactory{})

	for _, factory := range factories {
		conn := factory.CreateConnection()
		if conn == nil {
			t.Error("Expected connection")
		}

		cmd := factory.CreateCommand()
		if cmd == nil {
			t.Error("Expected command")
		}
	}
}

// ============ 基准测试 ============

func BenchmarkWindowsFactory(b *testing.B) {
	factory := &WindowsFactory{}
	for i := 0; i < b.N; i++ {
		_ = factory.CreateButton()
		_ = factory.CreateTextBox()
		_ = factory.CreateMenu()
	}
}

func BenchmarkCreateUIFactory(b *testing.B) {
	for i := 0; i < b.N; i++ {
		factory, _ := CreateUIFactory(Windows)
		_ = factory.CreateButton()
	}
}

func BenchmarkMySQLFactory(b *testing.B) {
	factory := &MySQLFactory{}
	for i := 0; i < b.N; i++ {
		_ = factory.CreateConnection()
		_ = factory.CreateCommand()
		_ = factory.CreateAdapter()
	}
}

func BenchmarkCreateDatabaseFactory(b *testing.B) {
	for i := 0; i < b.N; i++ {
		factory, _ := CreateDatabaseFactory(MySQL)
		_ = factory.CreateConnection()
	}
}
