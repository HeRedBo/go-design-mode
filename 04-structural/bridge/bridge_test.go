package bridge

import "testing"

func TestEmailNotifier(t *testing.T) {
	sender := &EmailSender{}
	notifier := NewNotifier(sender)

	if notifier.GetType() != "Email" {
		t.Errorf("Expected type 'Email'")
	}

	err := notifier.Send("Hello", "user@example.com")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestAlertNotifier(t *testing.T) {
	sender := &SMSSender{}
	alert := NewAlertNotifier(sender, "CRITICAL")

	err := alert.Send("Server down", "admin@example.com")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestBatchNotifier(t *testing.T) {
	sender := &PushSender{}
	batch := NewBatchNotifier(sender)

	recipients := []string{"user1", "user2", "user3"}
	err := batch.SendBatch("Update available", recipients)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestCircleShape(t *testing.T) {
	vectorRenderer := &VectorRenderer{}
	circle := NewCircle(vectorRenderer, 5.0)

	circle.Draw()
}

func TestRectangleShape(t *testing.T) {
	rasterRenderer := &RasterRenderer{}
	rect := NewRectangle(rasterRenderer, 10.0, 20.0)

	rect.Draw()
}

func TestDatabase(t *testing.T) {
	mysqlDriver := &MySQLDriver{}
	db := NewDatabase(mysqlDriver, "mydb")

	err := db.Connect("localhost", 3306)
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}

	_, err = db.Query("SELECT * FROM users")
	if err != nil {
		t.Errorf("Query failed: %v", err)
	}

	err = db.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}
}

func BenchmarkNotifier(b *testing.B) {
	sender := &EmailSender{}
	notifier := NewNotifier(sender)

	for i := 0; i < b.N; i++ {
		_ = notifier.Send("test", "user@example.com")
	}
}

func BenchmarkShape(b *testing.B) {
	renderer := &VectorRenderer{}
	circle := NewCircle(renderer, 5.0)

	for i := 0; i < b.N; i++ {
		circle.Draw()
	}
}
