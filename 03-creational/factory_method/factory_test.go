package factorymethod

import (
	"testing"
)

// ============ 形状工厂测试 ============

func TestCircleFactory(t *testing.T) {
	factory := &CircleFactory{Radius: 5.0}
	shape := factory.CreateShape()

	if shape.Name() != "Circle" {
		t.Errorf("Expected name 'Circle', got '%s'", shape.Name())
	}

	expectedArea := 78.53981633974483
	if shape.Area() != expectedArea {
		t.Errorf("Expected area %f, got %f", expectedArea, shape.Area())
	}

	expectedDraw := "Drawing Circle with radius 5.00"
	if shape.Draw() != expectedDraw {
		t.Errorf("Expected draw '%s', got '%s'", expectedDraw, shape.Draw())
	}
}

func TestRectangleFactory(t *testing.T) {
	factory := &RectangleFactory{Width: 4.0, Height: 5.0}
	shape := factory.CreateShape()

	if shape.Name() != "Rectangle" {
		t.Errorf("Expected name 'Rectangle', got '%s'", shape.Name())
	}

	if shape.Area() != 20.0 {
		t.Errorf("Expected area 20.0, got %f", shape.Area())
	}

	expectedDraw := "Drawing Rectangle 4.00x5.00"
	if shape.Draw() != expectedDraw {
		t.Errorf("Expected draw '%s', got '%s'", expectedDraw, shape.Draw())
	}
}

func TestTriangleFactory(t *testing.T) {
	factory := &TriangleFactory{Base: 6.0, Height: 4.0}
	shape := factory.CreateShape()

	if shape.Name() != "Triangle" {
		t.Errorf("Expected name 'Triangle', got '%s'", shape.Name())
	}

	if shape.Area() != 12.0 {
		t.Errorf("Expected area 12.0, got %f", shape.Area())
	}

	expectedDraw := "Drawing Triangle base=6.00, height=4.00"
	if shape.Draw() != expectedDraw {
		t.Errorf("Expected draw '%s', got '%s'", expectedDraw, shape.Draw())
	}
}

func TestCreateShape(t *testing.T) {
	// 测试创建圆形
	circle, err := CreateShape(ShapeCircle, ShapeConfig{Radius: 3.0})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if circle.Name() != "Circle" {
		t.Errorf("Expected 'Circle', got '%s'", circle.Name())
	}

	// 测试创建矩形
	rect, err := CreateShape(ShapeRectangle, ShapeConfig{Width: 3.0, Height: 4.0})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if rect.Area() != 12.0 {
		t.Errorf("Expected area 12.0, got %f", rect.Area())
	}

	// 测试创建三角形
	tri, err := CreateShape(ShapeTriangle, ShapeConfig{Base: 5.0, Height: 6.0})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if tri.Area() != 15.0 {
		t.Errorf("Expected area 15.0, got %f", tri.Area())
	}

	// 测试无效类型
	_, err = CreateShape("invalid", ShapeConfig{})
	if err == nil {
		t.Error("Expected error for invalid shape type")
	}

	// 测试无效参数
	_, err = CreateShape(ShapeCircle, ShapeConfig{Radius: -1.0})
	if err == nil {
		t.Error("Expected error for negative radius")
	}
}

// ============ 文档工厂测试 ============

func TestPDFDocumentFactory(t *testing.T) {
	factory := &PDFDocumentFactory{}
	doc := factory.CreateDocument("test.pdf")

	if doc.GetType() != "PDF" {
		t.Errorf("Expected type 'PDF', got '%s'", doc.GetType())
	}

	expected := "Opening PDF document: test.pdf"
	if doc.Open() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, doc.Open())
	}
}

func TestWordDocumentFactory(t *testing.T) {
	factory := &WordDocumentFactory{}
	doc := factory.CreateDocument("report.docx")

	if doc.GetType() != "Word" {
		t.Errorf("Expected type 'Word', got '%s'", doc.GetType())
	}

	expected := "Opening Word document: report.docx"
	if doc.Open() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, doc.Open())
	}
}

func TestExcelDocumentFactory(t *testing.T) {
	factory := &ExcelDocumentFactory{}
	doc := factory.CreateDocument("data.xlsx")

	if doc.GetType() != "Excel" {
		t.Errorf("Expected type 'Excel', got '%s'", doc.GetType())
	}

	expected := "Opening Excel document: data.xlsx"
	if doc.Open() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, doc.Open())
	}
}

func TestCreateDocument(t *testing.T) {
	// 测试创建 PDF
	pdf, err := CreateDocument(DocPDF, "file.pdf")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if pdf.GetType() != "PDF" {
		t.Errorf("Expected 'PDF', got '%s'", pdf.GetType())
	}

	// 测试创建 Word
	word, err := CreateDocument(DocWord, "file.docx")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if word.GetType() != "Word" {
		t.Errorf("Expected 'Word', got '%s'", word.GetType())
	}

	// 测试创建 Excel
	excel, err := CreateDocument(DocExcel, "file.xlsx")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if excel.GetType() != "Excel" {
		t.Errorf("Expected 'Excel', got '%s'", excel.GetType())
	}

	// 测试无效类型
	_, err = CreateDocument("invalid", "file.txt")
	if err == nil {
		t.Error("Expected error for invalid document type")
	}
}

// ============ 多态测试 ============

func TestShapePolymorphism(t *testing.T) {
	// 使用接口实现多态
	shapes := []Shape{
		&Circle{Radius: 5.0},
		&Rectangle{Width: 4.0, Height: 6.0},
		&Triangle{Base: 3.0, Height: 8.0},
	}

	expectedAreas := []float64{
		78.53981633974483, // Circle
		24.0,              // Rectangle
		12.0,              // Triangle
	}

	for i, shape := range shapes {
		if shape.Area() != expectedAreas[i] {
			t.Errorf("Shape %d: expected area %f, got %f",
				i, expectedAreas[i], shape.Area())
		}
	}
}

func TestDocumentPolymorphism(t *testing.T) {
	// 使用接口实现多态
	docs := []Document{
		&PDFDocument{Filename: "a.pdf"},
		&WordDocument{Filename: "b.docx"},
		&ExcelDocument{Filename: "c.xlsx"},
	}

	expectedTypes := []string{"PDF", "Word", "Excel"}

	for i, doc := range docs {
		if doc.GetType() != expectedTypes[i] {
			t.Errorf("Doc %d: expected type '%s', got '%s'",
				i, expectedTypes[i], doc.GetType())
		}
	}
}

// ============ 基准测试 ============

func BenchmarkCircleFactory(b *testing.B) {
	factory := &CircleFactory{Radius: 5.0}
	for i := 0; i < b.N; i++ {
		_ = factory.CreateShape()
	}
}

func BenchmarkCreateShape(b *testing.B) {
	config := ShapeConfig{Radius: 5.0}
	for i := 0; i < b.N; i++ {
		_, _ = CreateShape(ShapeCircle, config)
	}
}

func BenchmarkPDFDocumentFactory(b *testing.B) {
	factory := &PDFDocumentFactory{}
	for i := 0; i < b.N; i++ {
		_ = factory.CreateDocument("test.pdf")
	}
}
