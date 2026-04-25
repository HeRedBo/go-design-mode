// Package composite 演示组合模式 (Composite Pattern)
//
// 组合模式解决的问题：
// 1. 将对象组合成树形结构表示"部分-整体"的层次结构
// 2. 客户端统一处理单个对象和组合对象
// 3. 忽略组合对象与单个对象的差异
//
// 使用场景：
// - 文件系统（文件和目录）
// - UI 组件（控件和容器）
// - 组织架构（员工和部门）
// - 菜单系统（菜单项和子菜单）
//
// 核心组件：
// 1. Component 组件：定义接口
// 2. Leaf 叶子：单个对象，无子节点
// 3. Composite 组合：容器对象，有子节点
//
// 透明方式 vs 安全方式：
// - 透明方式：Component 声明所有方法，叶子节点的实现无意义
// - 安全方式：Component 只声明共有方法，管理子节点的方法在 Composite 中
package composite

import (
	"fmt"
	"strings"
)

// ============ 文件系统示例 ============

// FileSystemComponent 文件系统组件接口
type FileSystemComponent interface {
	Name() string
	Size() int64
	Display(indent string)
}

// File 文件（叶子节点）
type File struct {
	name string
	size int64
}

// NewFile 创建文件
func NewFile(name string, size int64) *File {
	return &File{name: name, size: size}
}

// Name 返回文件名
func (f *File) Name() string {
	return f.name
}

// Size 返回文件大小
func (f *File) Size() int64 {
	return f.size
}

// Display 显示文件
func (f *File) Display(indent string) {
	fmt.Printf("%s📄 %s (%d bytes)\n", indent, f.name, f.size)
}

// Directory 目录（组合节点）
type Directory struct {
	name     string
	children []FileSystemComponent
}

// NewDirectory 创建目录
func NewDirectory(name string) *Directory {
	return &Directory{
		name:     name,
		children: make([]FileSystemComponent, 0),
	}
}

// Name 返回目录名
func (d *Directory) Name() string {
	return d.name
}

// Size 计算目录总大小
func (d *Directory) Size() int64 {
	var total int64
	for _, child := range d.children {
		total += child.Size()
	}
	return total
}

// Add 添加子组件
func (d *Directory) Add(component FileSystemComponent) {
	d.children = append(d.children, component)
}

// Remove 移除子组件
func (d *Directory) Remove(name string) bool {
	for i, child := range d.children {
		if child.Name() == name {
			d.children = append(d.children[:i], d.children[i+1:]...)
			return true
		}
	}
	return false
}

// GetChild 获取子组件
func (d *Directory) GetChild(name string) (FileSystemComponent, bool) {
	for _, child := range d.children {
		if child.Name() == name {
			return child, true
		}
	}
	return nil, false
}

// Display 显示目录结构
func (d *Directory) Display(indent string) {
	fmt.Printf("%s📁 %s\n", indent, d.name)
	for _, child := range d.children {
		child.Display(indent + "  ")
	}
}

// ============ 组织架构示例 ============

// Employee 员工接口
type Employee interface {
	Name() string
	Role() string
	Salary() float64
	GetLevel() int
	Display(indent string)
}

// Developer 开发者（叶子节点）
type Developer struct {
	name   string
	role   string
	salary float64
}

// NewDeveloper 创建开发者
func NewDeveloper(name, role string, salary float64) *Developer {
	return &Developer{name: name, role: role, salary: salary}
}

// Name 返回姓名
func (d *Developer) Name() string {
	return d.name
}

// Role 返回角色
func (d *Developer) Role() string {
	return d.role
}

// Salary 返回薪水
func (d *Developer) Salary() float64 {
	return d.salary
}

// GetLevel 返回级别
func (d *Developer) GetLevel() int {
	return 2
}

// Display 显示员工信息
func (d *Developer) Display(indent string) {
	fmt.Printf("%s👤 %s - %s ($%.2f)\n", indent, d.name, d.role, d.salary)
}

// Manager 经理（组合节点）
type Manager struct {
	name     string
	role     string
	salary   float64
	team     []Employee
}

// NewManager 创建经理
func NewManager(name, role string, salary float64) *Manager {
	return &Manager{
		name:   name,
		role:   role,
		salary: salary,
		team:   make([]Employee, 0),
	}
}

// Name 返回姓名
func (m *Manager) Name() string {
	return m.name
}

// Role 返回角色
func (m *Manager) Role() string {
	return m.role
}

// Salary 返回薪水
func (m *Manager) Salary() float64 {
	return m.salary
}

// GetLevel 返回级别
func (m *Manager) GetLevel() int {
	return 1
}

// AddEmployee 添加员工
func (m *Manager) AddEmployee(emp Employee) {
	m.team = append(m.team, emp)
}

// RemoveEmployee 移除员工
func (m *Manager) RemoveEmployee(name string) bool {
	for i, emp := range m.team {
		if emp.Name() == name {
			m.team = append(m.team[:i], m.team[i+1:]...)
			return true
		}
	}
	return false
}

// GetTeamSize 获取团队规模
func (m *Manager) GetTeamSize() int {
	return len(m.team)
}

// GetTotalSalary 计算总薪水
func (m *Manager) GetTotalSalary() float64 {
	total := m.salary
	for _, emp := range m.team {
		total += emp.Salary()
	}
	return total
}

// Display 显示组织架构
func (m *Manager) Display(indent string) {
	fmt.Printf("%s👔 %s - %s ($%.2f) [Team: %d]\n",
		indent, m.name, m.role, m.salary, len(m.team))
	for _, emp := range m.team {
		emp.Display(indent + "  ")
	}
}

// ============ 菜单系统示例 ============

// MenuItem 菜单项接口
type MenuItem interface {
	Title() string
	Display(indent string)
}

// MenuOption 菜单选项（叶子节点）
type MenuOption struct {
	title    string
	action   func()
}

// NewMenuOption 创建菜单选项
func NewMenuOption(title string, action func()) *MenuOption {
	return &MenuOption{title: title, action: action}
}

// Title 返回标题
func (m *MenuOption) Title() string {
	return m.title
}

// Display 显示菜单项
func (m *MenuOption) Display(indent string) {
	fmt.Printf("%s• %s\n", indent, m.title)
}

// Execute 执行动作
func (m *MenuOption) Execute() {
	if m.action != nil {
		m.action()
	}
}

// SubMenu 子菜单（组合节点）
type SubMenu struct {
	title string
	items []MenuItem
}

// NewSubMenu 创建子菜单
func NewSubMenu(title string) *SubMenu {
	return &SubMenu{
		title: title,
		items: make([]MenuItem, 0),
	}
}

// Title 返回标题
func (s *SubMenu) Title() string {
	return s.title
}

// AddItem 添加菜单项
func (s *SubMenu) AddItem(item MenuItem) {
	s.items = append(s.items, item)
}

// Display 显示菜单
func (s *SubMenu) Display(indent string) {
	fmt.Printf("%s📋 %s\n", indent, s.title)
	for _, item := range s.items {
		item.Display(indent + "  ")
	}
}

// ============ 辅助函数 ============

// PrintTree 打印组件树
func PrintTree(component FileSystemComponent) {
	component.Display("")
}

// CalculateTotalSize 计算总大小
func CalculateTotalSize(component FileSystemComponent) int64 {
	return component.Size()
}

// GetOrganizationChart 获取组织架构图字符串
func GetOrganizationChart(emp Employee) string {
	var builder strings.Builder
	writeEmployee(&builder, emp, "")
	return builder.String()
}

func writeEmployee(builder *strings.Builder, emp Employee, indent string) {
	if mgr, ok := emp.(*Manager); ok {
		builder.WriteString(fmt.Sprintf("%s👔 %s - %s [Team: %d]\n",
			indent, mgr.Name(), mgr.Role(), mgr.GetTeamSize()))
		for _, teamMember := range mgr.team {
			writeEmployee(builder, teamMember, indent+"  ")
		}
	} else {
		builder.WriteString(fmt.Sprintf("%s👤 %s - %s\n",
			indent, emp.Name(), emp.Role()))
	}
}
