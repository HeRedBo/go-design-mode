package facade

import "testing"

func TestHomeTheaterFacade(t *testing.T) {
	facade := NewHomeTheaterFacade()
	
	// 测试观看电影
	facade.WatchMovie("Inception")
	
	// 测试结束电影
	facade.EndMovie()
}

func TestOrderFacade(t *testing.T) {
	facade := NewOrderFacade()
	
	orderID, err := facade.PlaceOrder(
		"PROD-001",
		2,
		99.99,
		"1234-5678-9012-3456",
		"123 Main St",
		"customer@example.com",
	)
	
	if err != nil {
		t.Errorf("Order failed: %v", err)
	}
	
	if orderID == "" {
		t.Error("Expected order ID")
	}
}

func TestDatabaseFacade(t *testing.T) {
	facade := NewDatabaseFacade()
	
	// 测试简单查询
	results, err := facade.ExecuteQuery(
		"localhost:5432/mydb",
		"SELECT * FROM users",
	)
	
	if err != nil {
		t.Errorf("Query failed: %v", err)
	}
	
	if len(results) == 0 {
		t.Error("Expected results")
	}
	
	// 测试事务
	queries := []string{
		"INSERT INTO users (name) VALUES ('Alice')",
		"UPDATE users SET age = 30 WHERE name = 'Alice'",
	}
	
	err = facade.ExecuteTransaction("localhost:5432/mydb", queries)
	if err != nil {
		t.Errorf("Transaction failed: %v", err)
	}
}

func BenchmarkHomeTheater(b *testing.B) {
	facade := NewHomeTheaterFacade()
	
	for i := 0; i < b.N; i++ {
		facade.WatchMovie("Test Movie")
		facade.EndMovie()
	}
}

func BenchmarkOrderPlacement(b *testing.B) {
	facade := NewOrderFacade()
	
	for i := 0; i < b.N; i++ {
		_, _ = facade.PlaceOrder(
			"PROD-001",
			1,
			50.0,
			"1234-5678-9012-3456",
			"123 Main St",
			"test@example.com",
		)
	}
}
