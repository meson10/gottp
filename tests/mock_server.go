package tests

import (
	"bytes"
	"io/ioutil"
	"log"
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
	if request.Method == "" {
		request.Method = "GET"
	}

	req, err := http.NewRequest(
		strings.ToUpper(request.Method),
		self.server.URL+request.Url,
		bytes.NewBuffer(utils.Encoder(request.Data)),
	)

	if err != nil {
		log.Panic(err)
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Panic(err)
	}

	var msg = new(MockResponse)
	utils.Decoder(data, msg)

	equality(msg)
}

func NewServer() *MockServer {
	ts := httptest.NewServer(http.HandlerFunc(gottp.GlobalHandler))
	return &MockServer{ts}
}
