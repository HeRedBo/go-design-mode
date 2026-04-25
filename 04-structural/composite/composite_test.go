package composite

import "testing"

func TestFile(t *testing.T) {
	file := NewFile("test.txt", 1024)

	if file.Name() != "test.txt" {
		t.Errorf("Expected name 'test.txt'")
	}
	if file.Size() != 1024 {
		t.Errorf("Expected size 1024")
	}
}

func TestDirectory(t *testing.T) {
	dir := NewDirectory("src")
	dir.Add(NewFile("main.go", 500))
	dir.Add(NewFile("util.go", 300))

	if dir.Name() != "src" {
		t.Errorf("Expected name 'src'")
	}
	if dir.Size() != 800 {
		t.Errorf("Expected size 800, got %d", dir.Size())
	}

	// 测试删除
	removed := dir.Remove("main.go")
	if !removed {
		t.Error("Expected file to be removed")
	}
	if dir.Size() != 300 {
		t.Errorf("Expected size 300 after removal")
	}
}

func TestNestedDirectory(t *testing.T) {
	root := NewDirectory("root")
	src := NewDirectory("src")
	test := NewDirectory("test")

	src.Add(NewFile("main.go", 1000))
	src.Add(NewFile("lib.go", 500))
	test.Add(NewFile("main_test.go", 800))

	src.Add(test)
	root.Add(src)

	// root: 2300 bytes (1000 + 500 + 800)
	if root.Size() != 2300 {
		t.Errorf("Expected total size 2300, got %d", root.Size())
	}
}

func TestOrganization(t *testing.T) {
	cto := NewManager("Alice", "CTO", 150000)

	dev1 := NewDeveloper("Bob", "Senior Dev", 120000)
	dev2 := NewDeveloper("Charlie", "Junior Dev", 80000)

	cto.AddEmployee(dev1)
	cto.AddEmployee(dev2)

	if cto.GetTeamSize() != 2 {
		t.Errorf("Expected team size 2")
	}

	expectedSalary := 150000.0 + 120000.0 + 80000.0
	if cto.GetTotalSalary() != expectedSalary {
		t.Errorf("Expected total salary %.2f, got %.2f",
			expectedSalary, cto.GetTotalSalary())
	}
}

func TestMenuSystem(t *testing.T) {
	fileMenu := NewSubMenu("File")
	fileMenu.AddItem(NewMenuOption("New", func() {}))
	fileMenu.AddItem(NewMenuOption("Open", func() {}))
	fileMenu.AddItem(NewMenuOption("Save", func() {}))

	editMenu := NewSubMenu("Edit")
	editMenu.AddItem(NewMenuOption("Cut", func() {}))
	editMenu.AddItem(NewMenuOption("Copy", func() {}))
	editMenu.AddItem(NewMenuOption("Paste", func() {}))

	mainMenu := NewSubMenu("Main")
	mainMenu.AddItem(fileMenu)
	mainMenu.AddItem(editMenu)

	// 验证结构
	if mainMenu.Title() != "Main" {
		t.Errorf("Expected title 'Main'")
	}
}

func BenchmarkDirectory(b *testing.B) {
	dir := NewDirectory("root")
	for i := 0; i < 100; i++ {
		dir.Add(NewFile("file.go", 1000))
	}

	for i := 0; i < b.N; i++ {
		_ = dir.Size()
	}
}

func BenchmarkOrganization(b *testing.B) {
	ceo := NewManager("CEO", "CEO", 200000)
	for i := 0; i < 10; i++ {
		manager := NewManager("Manager", "Manager", 100000)
		for j := 0; j < 5; j++ {
			manager.AddEmployee(NewDeveloper("Dev", "Developer", 80000))
		}
		ceo.AddEmployee(manager)
	}

	for i := 0; i < b.N; i++ {
		_ = ceo.GetTotalSalary()
	}
}
