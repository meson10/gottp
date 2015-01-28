package main

import (
	"gopkg.in/simversity/gottp.v0"
)

var urls = []*gottp.Url{
	gottp.NewUrl("hello", "/hello/\\w{3,5}/?$", helloMessage),
}
