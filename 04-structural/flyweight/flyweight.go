// Package flyweight 演示享元模式 (Flyweight Pattern)
//
// 享元模式解决的问题：
// 1. 通过共享技术有效地支持大量细粒度对象
// 2. 减少内存占用，避免重复创建相同对象
// 3. 提高性能
//
// 使用场景：
// - 文本编辑器中的字符对象
// - 游戏地图中的图块
// - 连接池、线程池
// - 缓存系统
//
// 核心概念：
// - 内在状态（共享）：可以共享的状态
// - 外在状态（不共享）：每个对象独立的状态
//
// 核心组件：
// 1. Flyweight 享元接口：定义共享方法
// 2. ConcreteFlyweight 具体享元：实现共享
// 3. FlyweightFactory 享元工厂：管理和创建享元对象
package flyweight

import (
	"fmt"
	"sync"
)

// ============ 示例 1: 字符享元（文本编辑器） ============

// CharacterFlyweight 字符享元接口
type CharacterFlyweight interface {
	Display(position int)
}

// Character 字符享元
type Character struct {
	symbol rune // 内在状态（共享）
}

// NewCharacter 创建字符享元
func NewCharacter(symbol rune) *Character {
	return &Character{symbol: symbol}
}

// Display 显示字符
func (c *Character) Display(position int) {
	fmt.Printf("Character '%c' at position %d\n", c.symbol, position)
}

// CharacterFactory 字符工厂
type CharacterFactory struct {
	flyweights map[rune]*Character
	mu         sync.RWMutex
}

// NewCharacterFactory 创建字符工厂
func NewCharacterFactory() *CharacterFactory {
	return &CharacterFactory{
		flyweights: make(map[rune]*Character),
	}
}

// GetCharacter 获取字符享元
func (f *CharacterFactory) GetCharacter(symbol rune) *Character {
	f.mu.Lock()
	defer f.mu.Unlock()

	if char, exists := f.flyweights[symbol]; exists {
		return char
	}

	char := NewCharacter(symbol)
	f.flyweights[symbol] = char
	return char
}

// GetCount 获取已创建的享元数量
func (f *CharacterFactory) GetCount() int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return len(f.flyweights)
}

// Document 文档（使用享元）
type Document struct {
	characters []struct {
		char     *Character
		position int
	}
	factory *CharacterFactory
}

// NewDocument 创建文档
func NewDocument(factory *CharacterFactory) *Document {
	return &Document{
		factory: factory,
	}
}

// AddCharacter 添加字符
func (d *Document) AddCharacter(symbol rune, position int) {
	char := d.factory.GetCharacter(symbol)
	d.characters = append(d.characters, struct {
		char     *Character
		position int
	}{char: char, position: position})
}

// Display 显示文档
func (d *Document) Display() {
	for _, c := range d.characters {
		c.char.Display(c.position)
	}
}

// ============ 示例 2: 游戏地图图块享元 ============

// TileType 图块类型
type TileType string

const (
	Grass    TileType = "grass"
	Water    TileType = "water"
	Sand     TileType = "sand"
	Forest   TileType = "forest"
	Mountain TileType = "mountain"
)

// TileFlyweight 图块享元接口
type TileFlyweight interface {
	Render(x, y int)
	GetType() TileType
}

// Tile 图块享元
type Tile struct {
	tileType   TileType
	texture    string
	color      string
	isPassable bool
}

// NewTile 创建图块享元
func NewTile(tileType TileType, texture string, color string, isPassable bool) *Tile {
	return &Tile{
		tileType:   tileType,
		texture:    texture,
		color:      color,
		isPassable: isPassable,
	}
}

// Render 渲染图块
func (t *Tile) Render(x, y int) {
	fmt.Printf("[%s] at (%d,%d) - color: %s, passable: %v\n",
		t.tileType, x, y, t.color, t.isPassable)
}

// GetType 获取类型
func (t *Tile) GetType() TileType {
	return t.tileType
}

// TileFactory 图块工厂
type TileFactory struct {
	flyweights map[TileType]*Tile
	mu         sync.RWMutex
}

// NewTileFactory 创建图块工厂
func NewTileFactory() *TileFactory {
	factory := &TileFactory{
		flyweights: make(map[TileType]*Tile),
	}

	// 预创建所有图块类型
	factory.flyweights[Grass] = NewTile(Grass, "grass.png", "green", true)
	factory.flyweights[Water] = NewTile(Water, "water.png", "blue", false)
	factory.flyweights[Sand] = NewTile(Sand, "sand.png", "yellow", true)
	factory.flyweights[Forest] = NewTile(Forest, "forest.png", "darkgreen", true)
	factory.flyweights[Mountain] = NewTile(Mountain, "mountain.png", "gray", false)

	return factory
}

// GetTile 获取图块享元
func (f *TileFactory) GetTile(tileType TileType) *Tile {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.flyweights[tileType]
}

// GameMap 游戏地图
type GameMap struct {
	tiles [][]*Tile
	width, height int
	factory *TileFactory
}

// NewGameMap 创建游戏地图
func NewGameMap(factory *TileFactory, width, height int) *GameMap {
	tiles := make([][]*Tile, height)
	for i := range tiles {
		tiles[i] = make([]*Tile, width)
	}
	return &GameMap{
		tiles:   tiles,
		width:   width,
		height:  height,
		factory: factory,
	}
}

// SetTile 设置图块
func (m *GameMap) SetTile(x, y int, tileType TileType) {
	if x >= 0 && x < m.width && y >= 0 && y < m.height {
		m.tiles[y][x] = m.factory.GetTile(tileType)
	}
}

// Render 渲染地图
func (m *GameMap) Render() {
	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			if m.tiles[y] != nil && m.tiles[y][x] != nil {
				m.tiles[y][x].Render(x, y)
			}
		}
	}
}

// ============ 示例 3: 颜色享元 ============

// Color 颜色享元
type Color struct {
	red, green, blue int
	hex              string
}

// NewColor 创建颜色
func NewColor(red, green, blue int) *Color {
	return &Color{
		red:   red,
		green: green,
		blue:  blue,
		hex:   fmt.Sprintf("#%02X%02X%02X", red, green, blue),
	}
}

// Apply 应用颜色
func (c *Color) Apply(element string) {
	fmt.Printf("Applying color %s to %s\n", c.hex, element)
}

// ColorFactory 颜色工厂
type ColorFactory struct {
	flyweights map[string]*Color
	mu         sync.RWMutex
}

// NewColorFactory 创建颜色工厂
func NewColorFactory() *ColorFactory {
	return &ColorFactory{
		flyweights: make(map[string]*Color),
	}
}

// GetColor 获取颜色享元
func (f *ColorFactory) GetColor(red, green, blue int) *Color {
	f.mu.Lock()
	defer f.mu.Unlock()

	key := fmt.Sprintf("%d-%d-%d", red, green, blue)
	if color, exists := f.flyweights[key]; exists {
		return color
	}

	color := NewColor(red, green, blue)
	f.flyweights[key] = color
	return color
}

// GetCount 获取已创建的颜色数量
func (f *ColorFactory) GetCount() int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return len(f.flyweights)
}

// ============ 辅助函数 ============

// GetMemorySavings 计算内存节省
func GetMemorySavings(totalObjects, uniqueObjects int, objectSize int) string {
	withoutFlyweight := totalObjects * objectSize
	withFlyweight := uniqueObjects*objectSize + (totalObjects-uniqueObjects)*4 // 指针大小
	saved := withoutFlyweight - withFlyweight

	return fmt.Sprintf("Without flyweight: %d bytes, With flyweight: %d bytes, Saved: %d bytes (%.1f%%)",
		withoutFlyweight, withFlyweight, saved,
		float64(saved)/float64(withoutFlyweight)*100)
}
