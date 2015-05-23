package tests

import "testing"

func TestUrls(t *testing.T) {
	server := NewServer()
	defer server.Close()

	server.Test("/urls", "get", func(msg *MockResponse) {
		if msg.Status != 200 {
			t.Error("Expected Status return", msg.Status)
		}
	})
}
