package tests

import (
	"testing"

	"gopkg.in/simversity/gottp.v3/utils"
)

func TestUrls(t *testing.T) {
	server := NewServer()
	defer server.Close()

	req := MockRequest{}
	req.Url = "/urls"

	server.Test(&req, func(msg *MockResponse) {
		if msg.Status != 200 {
			t.Error("Expected Status return", msg.Status)
		}

		urlResponse := map[string]string{}
		utils.Convert(&msg.Data, &urlResponse)

		if _, ok := urlResponse["async_pipe"]; !ok {
			t.Error("Invalid response from Server")
		}
	})
}
