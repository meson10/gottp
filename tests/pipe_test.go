package tests

import (
	"testing"

	"gopkg.in/simversity/gottp.v3/utils"
)

func TestPipeNotAllowed(t *testing.T) {
	server := NewServer()
	defer server.Close()

	req := MockRequest{}
	req.Url = "/pipe"

	server.Test(&req, func(msg *MockResponse) {
		expected := 501

		if msg.Status != expected {
			t.Error("Expected Status", expected, "Got:", msg.Status)
		}
	})
}

func TestPipe(t *testing.T) {
	server := NewServer()
	defer server.Close()

	req := MockRequest{}
	req.Url = "/pipe"
	req.Method = "post"
	req.Data = map[string]interface{}{
		"stack": []map[string]interface{}{
			map[string]interface{}{
				"url":    "/urls",
				"method": "GET",
				"data":   map[string]string{},
			},
		},
	}

	server.Test(&req, func(msg *MockResponse) {
		pipeResponse := []MockResponse{}
		utils.Convert(&msg.Data, &pipeResponse)

		if len(pipeResponse) != 1 {
			t.Error("Expected Length of output 1")
		}

		if pipeResponse[0].Status != 200 {
			t.Error("Expected 200 from Server")
		}

	})
}

func TestAsyncPipe(t *testing.T) {
	server := NewServer()
	defer server.Close()

	req := MockRequest{}
	req.Url = "/async_pipe"
	req.Method = "post"
	req.Data = map[string]interface{}{
		"stack": []map[string]interface{}{
			map[string]interface{}{
				"url":    "/urls",
				"method": "GET",
				"data":   map[string]string{},
			},
		},
	}

	server.Test(&req, func(msg *MockResponse) {
		pipeResponse := []MockResponse{}
		utils.Convert(&msg.Data, &pipeResponse)

		if len(pipeResponse) != 1 {
			t.Error("Expected Length of output 1")
		}

		if pipeResponse[0].Status != 200 {
			t.Error("Expected 200 from Server")
		}

	})
}

func TestMixPipe(t *testing.T) {
	server := NewServer()
	defer server.Close()

	req := MockRequest{}
	req.Url = "/pipe"
	req.Method = "post"
	req.Data = map[string]interface{}{
		"stack": []map[string]interface{}{
			map[string]interface{}{
				"url":    "/urls",
				"method": "GET",
				"data":   map[string]string{},
			},
			map[string]interface{}{
				"url":    "/urls",
				"method": "POST",
				"data":   map[string]string{},
			},
			map[string]interface{}{
				"url":    "/urlsd",
				"method": "GET",
				"data":   map[string]string{},
			},
		},
	}

	server.Test(&req, func(msg *MockResponse) {
		pipeResponse := []MockResponse{}
		utils.Convert(&msg.Data, &pipeResponse)

		if len(pipeResponse) != 3 {
			t.Error("Expected Length of output 1")
		}

		if pipeResponse[0].Status != 200 {
			t.Error("Expected 200 from Server")
		}

		if pipeResponse[1].Status != 501 {
			t.Error("Post pipe call on /urls should not be allowed")
		}

		if pipeResponse[2].Status != 404 {
			t.Error("/urlsd should not be found")
		}
	})
}

func TestAsyncMixPipe(t *testing.T) {
	server := NewServer()
	defer server.Close()

	req := MockRequest{}
	req.Url = "/async_pipe"
	req.Method = "post"
	req.Data = map[string]interface{}{
		"stack": []map[string]interface{}{
			map[string]interface{}{
				"url":    "/urls",
				"method": "GET",
				"data":   map[string]string{},
			},
			map[string]interface{}{
				"url":    "/urls",
				"method": "POST",
				"data":   map[string]string{},
			},
			map[string]interface{}{
				"url":    "/urlsd",
				"method": "GET",
				"data":   map[string]string{},
			},
		},
	}

	server.Test(&req, func(msg *MockResponse) {
		pipeResponse := []MockResponse{}
		utils.Convert(&msg.Data, &pipeResponse)

		if len(pipeResponse) != 3 {
			t.Error("Expected Length of output 1")
		}

		if pipeResponse[0].Status != 200 {
			t.Error("Expected 200 from Server")
		}

		if pipeResponse[1].Status != 501 {
			t.Error("Post pipe call on /urls should not be allowed")
		}

		if pipeResponse[2].Status != 404 {
			t.Error("/urlsd should not be found")
		}
	})
}
