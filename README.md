gottp
=====

A tiny golang web framework

Installation
=============
```
go get github.com/Simversity/gottp
```

Getting Started
===============

*A sample application named helloWorld is available inside the tests directory of your checkout*

To start building a web service using gottp just create a new project with the following structure.

  * conf.go -> Configuration
  * main.go -> main Server engine
  * urls.go -> Register acceptable urls & handlers.
  * handlers/hello.go -> Register handlers that process the request.


conf.go
-------

A minimalist configuration looks like:

```
package main

import "github.com/Simversity/gottp"

type config struct {
	Gottp gottp.SettingsMap
}

func (self *config) MakeConfig(configPath string) {
	gottp.Settings = self.Gottp
}

var settings config
```

main.go
-------

A sample main.go would look like:

```
package main

import (
    "log"
	"github.com/Simversity/gottp"
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
```

urls.go
-------

Urls are of type gottp.Url

```
type Url struct {
	name    string //shortname of the url
	url     string //provided regular pattern
	handler func(r *Request) //ReuqestHandler 
	pattern *regexp.Regexp //Compiled Regular Expression
}
```

A sample urls.go looks like:

```
package main

import (
	"github.com/Simversity/gottp"
)

var urls = []*gottp.Url{
    gottp.NewUrl("hello", "/hello/\\w{3,5}/?$", handlers.HelloMessage),
}
```

This would match all urls that are like "/hello/world" or "/hello/greet"

handlers.go
-----------

A sample handler looks like:

```
package handlers

import (
	"github.com/Simversity/gottp"
)

func HelloMessage(req *gottp.Request) {
    req.Write("hello world")
}
```

Build & Run
-----------

```
go install test && test
```

Point your browser to http://127.0.0.1:8005/hello/check

Should give you a JSON output:

```
{
    "data": "hello world",
    "message": "",
    "status": 200
}
```

Available URLs
==============

You can visit http://127.0.0.1:8005/urls to access the URLs exposed by your application.

Sample Output:

```
{
    "data": {
        "hello": "/hello/\\w{3,5}/?$"
    },
    "message": "",
    "status": 200
}
```

Command Line Options & CFG files
================================

To be documented

Pipes & Async Pipes
===================

To be documented

Error Reporting
===============

To be documented

Accessing Request Arguments
===========================

To be documented
=======
... Working on Documentation, Should be available soon ...
