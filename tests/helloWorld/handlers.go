package main

import (
	"gopkg.in/simversity/gottp.v0"
)

func helloMessage(req *gottp.Request) {
	req.Write("hello world")
}
