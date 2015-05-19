package gottp

import (
	"net/http"
	"time"
)

func UrlHandler(w http.ResponseWriter, req *http.Request) {
	performUrls(w, req)
	return
}

func performUrls(w http.ResponseWriter, req *http.Request) {
	defer timeTrack(time.Now(), req)

	p := Request{Writer: w, Request: req, UrlArgs: nil}
	defer Tracer.Notify(getTracerExtra(&p))
	performRequest(new(allUrls), &p)
	return
}
