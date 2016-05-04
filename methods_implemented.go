package gottp

import (
	"encoding/json"
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
	client := gorequest.New()
	for _, method := range methodsSupported {
		switch method {
		case "GET":
			resp, _, err := client.Get(url).End()
			if err == nil {
				if resp.StatusCode == http.StatusOK {
					var data map[string]int
					json.NewDecoder(resp.Body).Decode(&data)
					if data["status"] != 501 && data["status"] != 0 {
						u.methodsImplemented = append(u.methodsImplemented, method)
					}
				}
			}
		case "POST":
			resp, _, err := client.Post(url).End()
			if err == nil {
				if resp.StatusCode == http.StatusOK {
					var data map[string]int
					json.NewDecoder(resp.Body).Decode(&data)
					if data["status"] != 501 && data["status"] != 0 {
						u.methodsImplemented = append(u.methodsImplemented, method)
					}
				}
			}
		case "PUT":
			resp, _, err := client.Put(url).End()
			if err == nil {
				if resp.StatusCode == http.StatusOK {
					var data map[string]int
					json.NewDecoder(resp.Body).Decode(&data)
					if data["status"] != 501 && data["status"] != 0 {
						u.methodsImplemented = append(u.methodsImplemented, method)
					}
				}
			}
		case "DELETE":
			resp, _, err := client.Delete(url).End()
			if err == nil {
				if resp.StatusCode == http.StatusOK {
					var data map[string]int
					json.NewDecoder(resp.Body).Decode(&data)
					if data["status"] != 501 && data["status"] != 0 {
						u.methodsImplemented = append(u.methodsImplemented, method)
					}
				}
			}
		case "HEAD":
			resp, _, err := client.Head(url).End()
			if err == nil {
				if resp.StatusCode == http.StatusOK {
					var data map[string]int
					json.NewDecoder(resp.Body).Decode(&data)
					if data["status"] != 501 && data["status"] != 0 {
						u.methodsImplemented = append(u.methodsImplemented, method)
					}
				}
			}
		case "PATCH":
			resp, _, err := client.Patch(url).End()
			if err == nil {
				if resp.StatusCode == http.StatusOK {
					var data map[string]int
					json.NewDecoder(resp.Body).Decode(&data)
					if data["status"] != 501 && data["status"] != 0 {
						u.methodsImplemented = append(u.methodsImplemented, method)
					}
				}
			}
		}
	}
	return
}
