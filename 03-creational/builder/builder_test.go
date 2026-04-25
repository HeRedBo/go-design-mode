package builder

import (
	"strings"
	"testing"
)

// ============ HTTP 请求建造者测试 ============

func TestHTTPRequestBuilder(t *testing.T) {
	req, err := NewRequestBuilder().
		WithMethod("POST").
		WithURL("https://api.example.com/users").
		WithHeader("Accept", "application/json").
		WithBody(`{"name":"test"}`).
		WithTimeout(60).
		WithRetries(3).
		WithAuth("token123").
		WithContentType("application/json").
		Build()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if req.method != "POST" {
		t.Errorf("Expected method POST, got %s", req.method)
	}

	if req.url != "https://api.example.com/users" {
		t.Errorf("Expected URL 'https://api.example.com/users', got %s", req.url)
	}

	if req.timeout != 60 {
		t.Errorf("Expected timeout 60, got %d", req.timeout)
	}

	if req.retries != 3 {
		t.Errorf("Expected retries 3, got %d", req.retries)
	}

	if req.authToken != "token123" {
		t.Errorf("Expected auth token 'token123', got %s", req.authToken)
	}

	str := req.String()
	if !strings.Contains(str, "POST") {
		t.Errorf("Expected string to contain 'POST'")
	}
	if !strings.Contains(str, "application/json") {
		t.Errorf("Expected string to contain 'application/json'")
	}
}

func TestHTTPRequestBuilderValidation(t *testing.T) {
	// 测试缺少 URL 的情况
	_, err := NewRequestBuilder().
		WithMethod("GET").
		Build()

	if err == nil {
		t.Error("Expected error for missing URL")
	}

	if !strings.Contains(err.Error(), "URL") {
		t.Errorf("Expected error message to contain 'URL', got: %s", err.Error())
	}
}

func TestHTTPRequestBuilderDefaults(t *testing.T) {
	req, err := NewRequestBuilder().
		WithURL("https://example.com").
		Build()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if req.method != "GET" {
		t.Errorf("Expected default method GET, got %s", req.method)
	}

	if req.timeout != 30 {
		t.Errorf("Expected default timeout 30, got %d", req.timeout)
	}
}

// ============ Director 测试 ============

func TestRequestDirector(t *testing.T) {
	builder := NewRequestBuilder()
	director := NewRequestDirector(builder)

	// 测试构建 API 请求
	apiReq, err := director.BuildAPIRequest("https://api.example.com", "secret-token")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if apiReq.method != "POST" {
		t.Errorf("Expected POST method, got %s", apiReq.method)
	}

	if apiReq.timeout != 60 {
		t.Errorf("Expected timeout 60, got %d", apiReq.timeout)
	}

	if apiReq.retries != 3 {
		t.Errorf("Expected retries 3, got %d", apiReq.retries)
	}

	// 测试构建简单 GET 请求
	getReq, err := director.BuildSimpleGet("https://example.com")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if getReq.method != "GET" {
		t.Errorf("Expected GET method, got %s", getReq.method)
	}

	if getReq.timeout != 30 {
		t.Errorf("Expected timeout 30, got %d", getReq.timeout)
	}
}

// ============ SQL 查询建造者测试 ============

func TestSQLQueryBuilder(t *testing.T) {
	query := NewSQLQueryBuilder().
		Select("id", "name", "email").
		From("users").
		Where("age > 18").
		Where("status = 'active'").
		OrderBy("name ASC").
		Limit(10).
		Offset(20).
		Build()

	expected := "SELECT id, name, email FROM users WHERE age > 18 AND status = 'active' ORDER BY name ASC LIMIT 10 OFFSET 20"
	if query.String() != expected {
		t.Errorf("Expected SQL:\n%s\nGot:\n%s", expected, query.String())
	}
}

func TestSQLQueryBuilderSimple(t *testing.T) {
	query := NewSQLQueryBuilder().
		Select("*").
		From("products").
		Where("price > 100").
		Build()

	expected := "SELECT * FROM products WHERE price > 100"
	if query.String() != expected {
		t.Errorf("Expected SQL:\n%s\nGot:\n%s", expected, query.String())
	}
}

func TestSQLQueryBuilderWithJoin(t *testing.T) {
	query := NewSQLQueryBuilder().
		Select("u.name", "o.total").
		From("users u").
		Join("INNER JOIN orders o ON u.id = o.user_id").
		Where("o.total > 100").
		OrderBy("o.total DESC").
		Limit(5).
		Build()

	expected := "SELECT u.name, o.total FROM users u INNER JOIN orders o ON u.id = o.user_id WHERE o.total > 100 ORDER BY o.total DESC LIMIT 5"
	if query.String() != expected {
		t.Errorf("Expected SQL:\n%s\nGot:\n%s", expected, query.String())
	}
}

func TestSQLQueryBuilderNoColumns(t *testing.T) {
	query := NewSQLQueryBuilder().
		From("users").
		Build()

	expected := "SELECT * FROM users"
	if query.String() != expected {
		t.Errorf("Expected SQL:\n%s\nGot:\n%s", expected, query.String())
	}
}

// ============ HTML 文档建造者测试 ============

func TestHTMLBuilder(t *testing.T) {
	doc := NewHTMLBuilder().
		WithTitle("My Page").
		AddStyle("body { margin: 0; }").
		AddScript("console.log('loaded');").
		AddHeadElement(`<meta charset="utf-8">`).
		AddBodyElement(`<h1>Hello World</h1>`).
		AddBodyElement(`<p>Welcome to my page</p>`).
		Build()

	html := doc.String()

	if !strings.Contains(html, "<title>My Page</title>") {
		t.Error("Expected HTML to contain title")
	}

	if !strings.Contains(html, "<style>body { margin: 0; }</style>") {
		t.Error("Expected HTML to contain style")
	}

	if !strings.Contains(html, "<script>console.log('loaded');</script>") {
		t.Error("Expected HTML to contain script")
	}

	if !strings.Contains(html, "<h1>Hello World</h1>") {
		t.Error("Expected HTML to contain h1")
	}

	if !strings.Contains(html, "<p>Welcome to my page</p>") {
		t.Error("Expected HTML to contain p")
	}

	if !strings.Contains(html, "<!DOCTYPE html>") {
		t.Error("Expected HTML to contain DOCTYPE")
	}
}

func TestHTMLBuilderMinimal(t *testing.T) {
	doc := NewHTMLBuilder().
		WithTitle("Simple Page").
		Build()

	html := doc.String()

	if !strings.Contains(html, "<title>Simple Page</title>") {
		t.Error("Expected HTML to contain title")
	}
}

// ============ 基准测试 ============

func BenchmarkHTTPRequestBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NewRequestBuilder().
			WithURL("https://api.example.com").
			WithMethod("POST").
			WithBody(`{"test":"data"}`).
			Build()
	}
}

func BenchmarkSQLQueryBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewSQLQueryBuilder().
			Select("id", "name").
			From("users").
			Where("active = true").
			OrderBy("id").
			Build()
	}
}

func BenchmarkHTMLBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewHTMLBuilder().
			WithTitle("Test Page").
			AddBodyElement("<p>Test</p>").
			Build()
	}
}

func BenchmarkRequestDirector(b *testing.B) {
	builder := NewRequestBuilder()
	director := NewRequestDirector(builder)

	for i := 0; i < b.N; i++ {
		_, _ = director.BuildAPIRequest("https://api.example.com", "token")
	}
}
