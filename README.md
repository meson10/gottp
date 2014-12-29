[![Build Status](https://travis-ci.org/Simversity/gottp.svg?branch=master)](https://travis-ci.org/Simversity/gottp)

gottp
=====

Gottp is not a regular front-end server to do user-facing CSS powered websites. It was designed using backend servers in mind and offers a variety of features like:

* Call Aggregation using Non-Blocking or Blocking Pipes. [1]
* Optionally Listens on Unix Domain socket.
* In-built error traceback emails.
* Optional data compression using zlib/gzip.
* Pagination support.

[1] Much like Batch requests in Facebook Graph API (https://developers.facebook.com/docs/graph-api/making-multiple-requests)


Installation
=============

Installation is as easy as:

```
go get github.com/Simversity/gottp
```

Getting Started
===============

*A sample application named helloWorld is available inside the tests directory of your checkout*

To start building a web service using gottp just create a new project with the following structure.

  * conf.go -> Configuration
  * main.go -> main Server engine
  * urls.go -> Register urls & corresponding handlers.
  * handlers.go -> handlers processing the request.

conf.go
-------

This section is applicable only when you need to provide application based settings alongside those used by gottp.

Configuration must implement the Configurer interface.

Configurer requires two method Sets:

* MakeConfig(string) which accepts the path of the .cfg file provided as an command line argument.
* GetGottpConfig() which must return the Settings 

```
type Configurer interface {
	MakeConfig(string)
	GetGottpConfig() *GottpSettings
}
```

A minimalist extended configuration looks like:

```
import (
	"github.com/Simversity/gottp/conf"
)

type config struct {
    Custom struct {
        VarOne string
        VarTwo string
    }
	Gottp conf.GottpSettings
}

func (self *config) MakeConfig(configPath string) {
	if configPath != "" {
		conf.MakeConfig(configPath, self)
	}
}

func (self *config) GetGottpConfig() *conf.GottpSettings {
	return &self.Gottp
}

var settings config
```

urls.go
-------

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

This would match all urls that are like "/hello/world" or "/hello/greet" but NOT /hello/123 and not /hello/yester


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

main.go
-------

A sample main.go looks like:

```
package main

import (
    "log"
	"github.com/Simversity/gottp"
)

func main() {
	gottp.BindHandlers(urls) //Urls is a slice of all the registered urls.
	gottp.MakeServer(&settings)
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

Configuration
-------------

Gottp allows you to provide .cfg files via the command line which is as easy as ./binary -config=path_to_cfg

Default gottp settings is a struct called GottpSettings

```
type GottpSettings struct {
	EmailHost     string //SMTP Host to send server Tracebacks.
	EmailPort     string //SMTP Port
	EmailUsername string //Username or Password to connect with SMTP
	EmailPassword string
	EmailSender   string   //Sender Name to be used for traceback.
	EmailFrom     string   //Verified sender email address like errors@example.com
	ErrorTo       []string //List of recipients for tracebacks.
	EmailDummy    bool     //Set to True, if Tracebacks should not be sent.
	Listen        string   //Address to Listen on default: 127.0.0.1:8005
}
```

.cfg sample
-----------

You can provide a simple configuration file in .cfg format to load the settings which would look like this:

```
[custom]
VarOne="one"
VarTwo="two"

[gottp]
listen="/tmp/custom.sock"
EmailHost="email-smtp.us-east-1.amazonaws.com"
EmailPort="587"
EmailPassword="TDVAGCWCTCTWCTCQ&&*!!*!*!*/NeURB5"
EmailUsername="HelloWorldSample"
EmailSender="My Gottp Server"
EmailFrom="errors@example.com"
ErrorTo="dev@example.com"
EmailDummy=false
```

URLs
----

Urls are of type gottp.Url

```
type Url struct {
	name    string //shortname of the url
	url     string //provided regular pattern
	handler func(r *Request) //ReuqestHandler
	pattern *regexp.Regexp //Compiled Regular Expression
}
```

URLs can be constructed using gottp.NewUrl which accepts shortname, regular expression & request Handler respectively.


Request Handler
---------------

A request handler supplied to URLs can be a function or a closure that accepts gottp.Request as argument.

Handler is responsible for writing the request data using request.Write(). Do note that gottp is purely JSON based so data is written with headers application/json.

This can be improved with time, I intend to support other content-types by exposing extensible middlewares.

Request exposes a few method structs:

GetArguments() *utils.Q

_Returns a map of all arguments passed to the request. This includes Body arguments in case of a PUT/POST request, url GET arguments and named arguments captured in URL regular expression. You can call this over as it handles caching internally._

GetArgument(key string) interface{}

ConvertArguments(dest interface{})
ConvertArgument(key string, dest interface{})

These methods come quite handy when you have to deal with the arguments. For a request with GET arguments that look like:

```
?abc=world&param=1&check=0
```

You can either fetch individual arguments like:

```
abcVar, _ := req.GetArgument("abc").(string)
log.Println(abcVar)

```

Or initialize a struct to convert all the arguments:


```
type params struct {
    abc string
    param int
    check int
}

getArgs := params{}
req.ConvertArguments(&getArgs)
log.Println(getArgs.abc)
```

Available URLs
==============

You can visit http://{host}:{port}/urls to access the URLs exposed by your application.

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
Pipes & Async Pipes
===================

To be documented

Error Reporting
===============

To be documented