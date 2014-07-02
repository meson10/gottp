package main

import (
	"github.com/Simversity/gottp"
)

func helloMessage(req *gottp.Request) {
	req.Write("hello world")
}
