package tests

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"

	"gopkg.in/simversity/gottp.v2/utils"
	"gopkg.in/simversity/gottp.v3"
)

type MockResponse struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type MockServer struct {
	server *httptest.Server
}

func (self *MockServer) Close() {
	self.server.Close()
}

func (self *MockServer) Test(url, method string, equality func(response *MockResponse)) {
	url = self.server.URL + url

	res, err := http.Get(url)
	if err != nil {
		log.Panic(err)
	}

	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
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
