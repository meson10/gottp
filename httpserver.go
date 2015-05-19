package gottp

import (
	"log"
	"net/http"
)

const serverUA = "Gottp Server"
const ERROR = `Oops! An Internal Error occured while performing that action.
Please try again later`

func performRequest(handler Handler, p *Request) {
	method := (*p).Request.Method

	switch method {
	case "GET":
		handler.Get(p)
	case "POST":
		handler.Post(p)
	case "PUT":
		handler.Put(p)
	case "DELETE":
		handler.Delete(p)
	case "HEAD":
		handler.Head(p)
	case "OPTIONS":
		handler.Options(p)
	case "PATCH":
		handler.Patch(p)
	default:
		log.Println("Unsupported method", method)
	}
}

func bindHandlers() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	http.HandleFunc("/async-pipe", AsyncPipeHandler)
	http.HandleFunc("/async-pipe/", AsyncPipeHandler)
	http.HandleFunc("/pipe", PipeHandler)
	http.HandleFunc("/pipe/", PipeHandler)
	http.HandleFunc("/urls", UrlHandler)
	http.HandleFunc("/urls/", UrlHandler)
	http.HandleFunc("/", GlobalHandler)
}
