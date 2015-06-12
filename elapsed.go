package gottp

import (
	"log"
	"net/http"
	"time"
)

func timeTrack(start time.Time, req *http.Request) {
	elapsed := time.Since(start)
	log.Printf("[%s] %s %s %s\n", req.Method, req.URL, req.RemoteAddr, elapsed)
}

func pipeTimeTrack(start time.Time, req *http.Request, pipeUrls string) {
	elapsed := time.Since(start)
	defer log.Printf("[%s] %s %s [%s] %s\n", req.Method, req.URL, req.RemoteAddr, pipeUrls, elapsed)
}
