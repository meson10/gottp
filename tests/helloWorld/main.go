// +build !appengine

package main

import (
	"log"

	"gopkg.in/simversity/gottp.v0"
)

func sysInit() {
	<-(gottp.SysInitChan) //Buffered Channel to receive the server upstart boolean
	log.Println("System is ready to Serve")
}

func main() {
	go sysInit()
	gottp.BindHandlers(urls) //Urls is a slice of all the registered urls.
	gottp.MakeServer(&settings)
}
