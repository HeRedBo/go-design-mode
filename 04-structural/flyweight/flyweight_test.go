package flyweight

import "testing"

func TestCharacterFactory(t *testing.T) {
	factory := NewCharacterFactory()

	// 创建多个相同字符
	char1 := factory.GetCharacter('a')
	char2 := factory.GetCharacter('a')
	char3 := factory.GetCharacter('b')

	// 验证享元共享
	if char1 != char2 {
		t.Error("Expected same flyweight for same character")
	}

	if char1 == char3 {
		t.Error("Expected different flyweights for different characters")
	}

	// 验证数量
	if factory.GetCount() != 2 {
		t.Errorf("Expected 2 flyweights, got %d", factory.GetCount())
	}
}

func TestDocument(t *testing.T) {
	factory := NewCharacterFactory()
	doc := NewDocument(factory)

	// 添加字符
	text := "hello"
	for i, char := range text {
		doc.AddCharacter(char, i)
	}

	// 验证享元共享（'l' 出现2次，应该共享）
	if factory.GetCount() != 4 { // h, e, l, o
		t.Errorf("Expected 4 flyweights, got %d", factory.GetCount())
	}
}

func TestTileFactory(t *testing.T) {
	factory := NewTileFactory()

	// 获取相同图块
	tile1 := factory.GetTile(Grass)
	tile2 := factory.GetTile(Grass)
	tile3 := factory.GetTile(Water)

	if tile1 != tile2 {
		t.Error("Expected same tile flyweight")
	}

	if tile1 == tile3 {
		t.Error("Expected different tile flyweights")
	}

	// 测试地图
	gameMap := NewGameMap(factory, 3, 3)
	gameMap.SetTile(0, 0, Grass)
	gameMap.SetTile(1, 0, Water)
	gameMap.SetTile(2, 0, Grass)

	// 验证内存节省
	savings := GetMemorySavings(1000000, 5, 100)
	if savings == "" {
		t.Error("Expected memory savings info")
	}
}

func TestColorFactory(t *testing.T) {
	factory := NewColorFactory()

	// 创建相同颜色
	color1 := factory.GetColor(255, 0, 0)
	color2 := factory.GetColor(255, 0, 0)
	color3 := factory.GetColor(0, 255, 0)

	if color1 != color2 {
		t.Error("Expected same color flyweight")
	}

	if color1 == color3 {
		t.Error("Expected different color flyweights")
	}

	if factory.GetCount() != 2 {
		t.Errorf("Expected 2 colors, got %d", factory.GetCount())
	}
}

func BenchmarkDocumentCreation(b *testing.B) {
	factory := NewCharacterFactory()

	for i := 0; i < b.N; i++ {
		doc := NewDocument(factory)
		text := "hello world"
		for j, char := range text {
			doc.AddCharacter(char, j)
		}
	}
}

func BenchmarkTileRendering(b *testing.B) {
	factory := NewTileFactory()
	gameMap := NewGameMap(factory, 100, 100)

	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			if (x+y)%2 == 0 {
				gameMap.SetTile(x, y, Grass)
			} else {
				gameMap.SetTile(x, y, Water)
			}
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gameMap.Render()
	}
}
