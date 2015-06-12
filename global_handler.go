package gottp

import (
	"net/http"
	"time"
)

func GlobalHandler(w http.ResponseWriter, req *http.Request) {
	defer timeTrack(time.Now(), req)
	request := Request{Writer: w, Request: req}
	doRequest(&request, &boundUrls)
	return
}
