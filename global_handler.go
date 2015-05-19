package gottp

import (
	"net/http"
	"time"
)

func GlobalHandler(w http.ResponseWriter, req *http.Request) {
	defer timeTrack(time.Now(), req)

	request := Request{Writer: w, Request: req}
	defer Tracer.Notify(getTracerExtra(&request))

	for _, url := range boundUrls {
		urlArgs, err := url.MakeUrlArgs(&req.URL.Path)
		if !err {
			request.UrlArgs = urlArgs
			performRequest(url.handler, &request)
			return
		}
	}

	e := HttpError{404, req.URL.Path + " not Found"}
	request.Raise(e)
	return
}
