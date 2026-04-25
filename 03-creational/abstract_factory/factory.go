// Package abstractfactory 演示抽象工厂模式 (Abstract Factory Pattern)
//
// 抽象工厂模式解决的问题:
// 1. 创建一系列相关或相互依赖的对象
// 2. 客户端不需要知道具体类的名称
// 3. 保证产品族的兼容性
//
// 与工厂方法的区别:
// - 工厂方法:一个工厂创建一个产品
// - 抽象工厂:一个工厂创建多个相关产品(产品族)
//
// 使用场景:
// - UI 组件库(不同操作系统的按钮、文本框、菜单)
// - 数据库驱动(不同数据库的连接、命令、适配器)
// - 主题系统(不同主题的窗口、按钮、颜色)
//
// Go 中的实现特点:
// - 使用接口定义工厂和产品
// - 利用组合而非继承
// - 返回接口类型,隐藏实现细节
//
// 核心组件:
// 1. AbstractFactory 接口:定义创建产品族的方法
// 2. AbstractProduct 接口:定义产品通用行为
// 3. ConcreteFactory:具体工厂实现,创建特定产品族
// 4. ConcreteProduct:具体产品实现
package abstractfactory

import "fmt"

// ============ 产品族 1: UI 组件 ============

// Button 按钮接口
type Button interface {
	Render() string
	Click() string
}

// TextBox 文本框接口
type TextBox interface {
	Render() string
	GetText() string
	SetText(text string)
}

// Menu 菜单接口
type Menu interface {
	Render() string
	AddItem(item string)
}

// ============ Windows 产品族 ============

// WindowsButton Windows 按钮
type WindowsButton struct{}

// Render 渲染按钮
func (b *WindowsButton) Render() string {
	return "Rendering Windows-style button"
}

// Click 点击按钮
func (b *WindowsButton) Click() string {
	return "Windows button clicked with ripple effect"
}

// WindowsTextBox Windows 文本框
type WindowsTextBox struct {
	text string
}

// Render 渲染文本框
func (t *WindowsTextBox) Render() string {
	return "Rendering Windows-style text box"
}

// GetText 获取文本
func (t *WindowsTextBox) GetText() string {
	return t.text
}

// SetText 设置文本
func (t *WindowsTextBox) SetText(text string) {
	t.text = text
}

// WindowsMenu Windows 菜单
type WindowsMenu struct {
	items []string
}

// Render 渲染菜单
func (m *WindowsMenu) Render() string {
	return "Rendering Windows-style menu with native look"
}

// AddItem 添加菜单项
func (m *WindowsMenu) AddItem(item string) {
	m.items = append(m.items, item)
}

// ============ macOS 产品族 ============

// MacButton macOS 按钮
type MacButton struct{}

// Render 渲染按钮
func (b *MacButton) Render() string {
	return "Rendering macOS button with rounded corners"
}

// Click 点击按钮
func (b *MacButton) Click() string {
	return "Mac button clicked with smooth animation"
}

// MacTextBox macOS 文本框
type MacTextBox struct {
	text string
}

// Render 渲染文本框
func (t *MacTextBox) Render() string {
	return "Rendering macOS text box with blur effect"
}

// GetText 获取文本
func (t *MacTextBox) GetText() string {
	return t.text
}

// SetText 设置文本
func (t *MacTextBox) SetText(text string) {
	t.text = text
}

// MacMenu macOS 菜单
type MacMenu struct {
	items []string
}

// Render 渲染菜单
func (m *MacMenu) Render() string {
	return "Rendering macOS menu with global menu bar"
}

// AddItem 添加菜单项
func (m *MacMenu) AddItem(item string) {
	m.items = append(m.items, item)
}

// ============ Linux 产品族 ============

// LinuxButton Linux 按钮
type LinuxButton struct{}

// Render 渲染按钮
func (b *LinuxButton) Render() string {
	return "Rendering Linux-style button with GTK theme"
}

// Click 点击按钮
func (b *LinuxButton) Click() string {
	return "Linux button clicked with classic feedback"
}

// LinuxTextBox Linux 文本框
type LinuxTextBox struct {
	text string
}

// Render 渲染文本框
func (t *LinuxTextBox) Render() string {
	return "Rendering Linux text box with GTK styling"
}

// GetText 获取文本
func (t *LinuxTextBox) GetText() string {
	return t.text
}

// SetText 设置文本
func (t *LinuxTextBox) SetText(text string) {
	t.text = text
}

// LinuxMenu Linux 菜单
type LinuxMenu struct {
	items []string
}

// Render 渲染菜单
func (m *LinuxMenu) Render() string {
	return "Rendering Linux menu with GTK theme"
}

// AddItem 添加菜单项
func (m *LinuxMenu) AddItem(item string) {
	m.items = append(m.items, item)
}

// ============ 抽象工厂接口 ============

// UIFactory UI 工厂接口（抽象工厂）
type UIFactory interface {
	CreateButton() Button
	CreateTextBox() TextBox
	CreateMenu() Menu
}

// ============ 具体工厂实现 ============

// WindowsFactory Windows UI 工厂
type WindowsFactory struct{}

// CreateButton 创建 Windows 按钮
func (f *WindowsFactory) CreateButton() Button {
	return &WindowsButton{}
}

// CreateTextBox 创建 Windows 文本框
func (f *WindowsFactory) CreateTextBox() TextBox {
	return &WindowsTextBox{}
}

// CreateMenu 创建 Windows 菜单
func (f *WindowsFactory) CreateMenu() Menu {
	return &WindowsMenu{items: make([]string, 0)}
}

// MacFactory macOS UI 工厂
type MacFactory struct{}

// CreateButton 创建 macOS 按钮
func (f *MacFactory) CreateButton() Button {
	return &MacButton{}
}

// CreateTextBox 创建 macOS 文本框
func (f *MacFactory) CreateTextBox() TextBox {
	return &MacTextBox{}
}

// CreateMenu 创建 macOS 菜单
func (f *MacFactory) CreateMenu() Menu {
	return &MacMenu{items: make([]string, 0)}
}

// LinuxFactory Linux UI 工厂
type LinuxFactory struct{}

// CreateButton 创建 Linux 按钮
func (f *LinuxFactory) CreateButton() Button {
	return &LinuxButton{}
}

// CreateTextBox 创建 Linux 文本框
func (f *LinuxFactory) CreateTextBox() TextBox {
	return &LinuxTextBox{}
}

// CreateMenu 创建 Linux 菜单
func (f *LinuxFactory) CreateMenu() Menu {
	return &LinuxMenu{items: make([]string, 0)}
}

// ============ 工厂注册表（便于动态创建） ============

// UIFactoryType UI 工厂类型
type UIFactoryType string

const (
	Windows UIFactoryType = "windows"
	MacOS   UIFactoryType = "macos"
	Linux   UIFactoryType = "linux"
)

// CreateUIFactory 根据类型创建 UI 工厂
func CreateUIFactory(factoryType UIFactoryType) (UIFactory, error) {
	switch factoryType {
	case Windows:
		return &WindowsFactory{}, nil
	case MacOS:
		return &MacFactory{}, nil
	case Linux:
		return &LinuxFactory{}, nil
	default:
		return nil, fmt.Errorf("unknown UI factory type: %s", factoryType)
	}
}

// ============ 实际应用示例：数据库驱动 ============

// Connection 数据库连接接口
type Connection interface {
	Connect() string
	Close() string
}

// Command 数据库命令接口
type Command interface {
	Execute(query string) string
}

// Adapter 数据库适配器接口
type Adapter interface {
	Adapt() string
}

// ============ MySQL 产品族 ============

// MySQLConnection MySQL 连接
type MySQLConnection struct{}

// Connect 连接数据库
func (c *MySQLConnection) Connect() string {
	return "Connecting to MySQL database"
}

// Close 关闭连接
func (c *MySQLConnection) Close() string {
	return "Closing MySQL connection"
}

// MySQLCommand MySQL 命令
type MySQLCommand struct{}

// Execute 执行命令
func (c *MySQLCommand) Execute(query string) string {
	return fmt.Sprintf("Executing MySQL query: %s", query)
}

// MySQLAdapter MySQL 适配器
type MySQLAdapter struct{}

// Adapt 适配
func (a *MySQLAdapter) Adapt() string {
	return "Using MySQL adapter with native protocol"
}

// ============ PostgreSQL 产品族 ============

// PostgreSQLConnection PostgreSQL 连接
type PostgreSQLConnection struct{}

// Connect 连接数据库
func (c *PostgreSQLConnection) Connect() string {
	return "Connecting to PostgreSQL database"
}

// Close 关闭连接
func (c *PostgreSQLConnection) Close() string {
	return "Closing PostgreSQL connection"
}

// PostgreSQLCommand PostgreSQL 命令
type PostgreSQLCommand struct{}

// Execute 执行命令
func (c *PostgreSQLCommand) Execute(query string) string {
	return fmt.Sprintf("Executing PostgreSQL query: %s", query)
}

// PostgreSQLAdapter PostgreSQL 适配器
type PostgreSQLAdapter struct{}

// Adapt 适配
func (a *PostgreSQLAdapter) Adapt() string {
	return "Using PostgreSQL adapter with libpq"
}

// ============ 数据库工厂接口和实现 ============

// DatabaseFactory 数据库工厂接口
type DatabaseFactory interface {
	CreateConnection() Connection
	CreateCommand() Command
	CreateAdapter() Adapter
}

// MySQLFactory MySQL 数据库工厂
type MySQLFactory struct{}

// CreateConnection 创建 MySQL 连接
func (f *MySQLFactory) CreateConnection() Connection {
	return &MySQLConnection{}
}

// CreateCommand 创建 MySQL 命令
func (f *MySQLFactory) CreateCommand() Command {
	return &MySQLCommand{}
}

// CreateAdapter 创建 MySQL 适配器
func (f *MySQLFactory) CreateAdapter() Adapter {
	return &MySQLAdapter{}
}

// PostgreSQLFactory PostgreSQL 数据库工厂
type PostgreSQLFactory struct{}

// CreateConnection 创建 PostgreSQL 连接
func (f *PostgreSQLFactory) CreateConnection() Connection {
	return &PostgreSQLConnection{}
}

// CreateCommand 创建 PostgreSQL 命令
func (f *PostgreSQLFactory) CreateCommand() Command {
	return &PostgreSQLCommand{}
}

// CreateAdapter 创建 PostgreSQL 适配器
func (f *PostgreSQLFactory) CreateAdapter() Adapter {
	return &PostgreSQLAdapter{}
}

// DatabaseType 数据库类型
type DatabaseType string

const (
	MySQL      DatabaseType = "mysql"
	PostgreSQL DatabaseType = "postgresql"
)

// CreateDatabaseFactory 创建数据库工厂
func CreateDatabaseFactory(dbType DatabaseType) (DatabaseFactory, error) {
	switch dbType {
	case MySQL:
		return &MySQLFactory{}, nil
	case PostgreSQL:
		return &PostgreSQLFactory{}, nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}
