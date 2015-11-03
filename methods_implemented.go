package gottp

import (
	"fmt"
	"net/http"

	"github.com/parnurzeal/gorequest"
)

func methondImplemented(u *Url) {
	u.methodsImplemented = append(u.methodsImplemented, "OPTIONS")
	if u.name == "async_pipe" || u.name == "urls" || u.name == "pipe" {
		return
	}
	methodsSupported := []string{"GET", "PUT", "POST", "DELETE", "HEAD", "OPTIONS", "PATCH"}
	url := "http://" + Address() + u.url
	fmt.Println(url)
	client := gorequest.New()
	for _, method := range methodsSupported {
		switch method {
		case "GET":
			resp, _, _ := client.Get(url).End()
			if resp.StatusCode != http.StatusNotImplemented {
				u.methodsImplemented = append(u.methodsImplemented, method)
			}
		case "POST":
			resp, _, _ := client.Post(url).End()
			if resp.StatusCode != http.StatusNotImplemented {
				u.methodsImplemented = append(u.methodsImplemented, method)
			}
		case "PUT":
			resp, _, _ := client.Put(url).End()
			if resp.StatusCode != http.StatusNotImplemented {
				u.methodsImplemented = append(u.methodsImplemented, method)
			}
		case "DELETE":
			resp, _, _ := client.Delete(url).End()
			if resp.StatusCode != http.StatusNotImplemented {
				u.methodsImplemented = append(u.methodsImplemented, method)
			}
		case "HEAD":
			resp, _, _ := client.Head(url).End()
			if resp.StatusCode != http.StatusNotImplemented {
				u.methodsImplemented = append(u.methodsImplemented, method)
			}
		case "PATCH":
			resp, _, _ := client.Patch(url).End()
			if resp.StatusCode != http.StatusNotImplemented {
				u.methodsImplemented = append(u.methodsImplemented, method)
			}
		}
	}
	return
}
