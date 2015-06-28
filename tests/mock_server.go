package tests

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"gopkg.in/simversity/gottp.v3"
	"gopkg.in/simversity/gottp.v3/utils"
)

type MockResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   error
}

type MockServer struct {
	server *httptest.Server
}

type MockRequest struct {
	Url    string
	Method string
	Data   map[string]interface{}
}

func (self *MockServer) Close() {
	self.server.Close()
}

func (self *MockServer) Test(request *MockRequest, equality func(response *MockResponse)) {
	var msg = new(MockResponse)

	if request.Method == "" {
		request.Method = "GET"
	}

	req, err := http.NewRequest(
		strings.ToUpper(request.Method),
		self.server.URL+request.Url,
		bytes.NewBuffer(utils.Encoder(request.Data)),
	)

	if err != nil {
		msg.Error = err
		equality(msg)
		return
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errors.New("Redirecting =>" + req.URL.String())
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		status := 0
		if strings.Contains(err.Error(), "Redirecting =>") {
			status = 301
		}

		msg.Status = status
		msg.Error = err
		msg.Message = err.Error()
		equality(msg)
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		msg.Error = err
		equality(msg)
		return
	}

	utils.Decoder(data, msg)
	equality(msg)
	return
}

func NewServer() *MockServer {
	ts := httptest.NewServer(http.HandlerFunc(gottp.GlobalHandler))
	return &MockServer{ts}
}
