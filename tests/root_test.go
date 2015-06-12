package tests

import "testing"

func TestRoot(t *testing.T) {
	server := NewServer()
	defer server.Close()

	req := MockRequest{}
	req.Url = "/"

	server.Test(&req, func(msg *MockResponse) {
		expected := 404

		if msg.Status != expected {
			t.Error("Expected Status", expected, "Got:", msg.Status)
		}
	})
}
