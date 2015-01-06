package main

import (
	"github.com/Simversity/gottp"
)

var urls = []*gottp.Url{
	gottp.NewUrl("hello", "/hello/\\w{3,5}/?$", helloMessage),
}
