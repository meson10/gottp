package gottp

import (
	"net/http"
	"time"
)

// GlobalHandler handles all requests for "/".
func GlobalHandler(w http.ResponseWriter, req *http.Request) {
	defer timeTrack(time.Now(), req)
	request := Request{Writer: w, Request: req}
	doRequest(&request, &boundUrls)
	return
}
