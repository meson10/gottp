[![Build Status](https://travis-ci.org/meson10/gottp.svg?branch=master)](https://travis-ci.org/meson10/gottp)

gottp
=====

Gottp is not a regular front-end server to do user-facing CSS powered websites. It was designed using backend servers in mind and offers a variety of features like:

* Background Workers
* Call Aggregation using Non-Blocking or Blocking Pipes. [1]
* Optionally Listens on Unix Domain socket.
* In-built error traceback emails.
* Optional data compression using zlib/gzip.
* Pagination support.
* Automatic listing of all exposed URLs

[1] Much like Batch requests in Facebook Graph API (https://developers.facebook.com/docs/graph-api/making-multiple-requests)


Installation
=============

Installation is as easy as:

```go
go get gopkg.in/meson10/gottp.v3
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

```go
type Configurer interface {
    MakeConfig(string)
        GetGottpConfig() *GottpSettings
}
```

A minimalist extended configuration looks like:

```go
import "gopkg.in/meson10/gottp.v3/conf"

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

```go
package main

import "gopkg.in/meson10/gottp.v3"

func init(){
    gottp.NewUrl("hello", "/hello/\\w{3,5}/?$", new(handlers.HelloMessage)),
}
```

This would match all urls that are like "/hello/world" or "/hello/greet" but NOT /hello/123 and not /hello/yester


handlers.go
-----------

A sample handler looks like:

```go
package handlers

import "gopkg.in/meson10/gottp.v3"

type HelloMessage struct {
  gottp.BaseHandler
}

func (self *HelloMessage) Get(req *gottp.Request) {
    req.Write("hello world")
}
```

main.go
-------

A sample main.go looks like:

```go
package main

import (
    "log"
    "gopkg.in/meson10/gottp.v3"
)

func main() {
  gottp.MakeServer(&settings)
}

```

Build & Run
-----------

```go
go install test && test
```

Point your browser to http://127.0.0.1:8005/hello/check

Should give you a JSON output:

```json
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

```go
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

```go
type Url struct {
    name    string //shortname of the url
    url     string //provided regular pattern
    handler func(r *Request) //ReuqestHandler
    pattern *regexp.Regexp //Compiled Regular Expression
}
```

URLs can be constructed using gottp.NewUrl which accepts shortname, regular expression & Handler interface implementor respectively.


Request Handler
---------------

A request handler must implement the Handler Interface which must expose the following methods

```go
type Handler interface {
	Get(request *Request)
	Put(request *Request)
	Post(request *Request)
	Delete(request *Request)
	Head(request *Request)
	Options(request *Request)
	Patch(request *Request)
}

```

gottp.BaseHandler has most common HTTP requests as method sets. So, If a struct uses BaseHandler as an embedded type it is sufficient to qualify as a Handler. To expose a new HTTP method for a URL type, implement the HTTP method set in the handler struct. See the Example below:

```go
type helloMessage struct {
	gottp.BaseHandler
}

func (self *helloMessage) Get(req *gottp.Request) {
	req.Write("hello world - GET")
}

func (self *helloMessage) Post(req *gottp.Request) {
	req.Write("hello world - POST")
}

```


Request
-------

Request exposes a few method structs:
```go
GetArguments() *utils.Q
```
_
Returns a map of all arguments passed to the request. This includes Body arguments in case of a PUT/POST request, url GET arguments and named arguments captured in URL regular expression. You can call this over as it handles caching internally.
_
```go
GetArgument(key string) interface{}
ConvertArguments(dest interface{})
ConvertArgument(key string, dest interface{})
```
These methods come quite handy when you have to deal with the arguments. For a request with GET arguments that look like:

```go
?abc=world&param=1&check=0
```

You can either fetch individual arguments like:

```go
abcVar, _ := req.GetArgument("abc").(string)
log.Println(abcVar)
```

Or initialize a struct to convert all the arguments:

```go
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

With backends powered by Gottp, consumers can simply access http://{host}:{port}/urls to fetch a json of all URLs exposed by the application. URLs are returned as a map of key: pattern where the key is the human-readable-identifier provided at the time of constructing URLS and the pattern is the URL regular expression pattern.

This helps in preventing hard coding of endpoints where consumers can fetch the URLs at application start time and reconstruct URLs using the indentifiers.

Sample Output:

```json
{
  "data": {
    "hello": "/hello/\\w{3,5}/?$"
  },
  "message": "",
  "status": 200
}
```

Pipes
=====

All gottp based servers allow clubbing of requests to be executed together. This is a very handy feature that is used for clubbing calls that need to processed one after the other. It can save the network handshaking costs as the handshake is only performed once and this does not effect server performance as all gottp requests are performed as goroutines.

Sequential Pipes
----------------

Example usage of such calls is:

1. Create a Comment
2. Mark the parent notification as read
3. Send out notification to the parent autho
4. Churn the new activity feed as per edge rank.

Using Pipes, call 1, 2, 3 & 4 can be combined into a single request to the server.

All the requests would sequentially evaluated and data would be returned in the same order as requested.

To submit a PIPE request issue a POST call on /pipe (available by default).

Async Pipes
-----------

Async pipes behave the same way, except the requests are executed in parallel
as go-routines. Async pipes are as fast as the slowest call in your request
stack.

Despite parallel execution they return the data in the same order as requested.

To submit an Async PIPE request issue a POST call on /async-pipe (available by default)


Pipe Request Object
-------------------

Requests submitted as PIPEs should again be a valid JSON dump of


```json
{
    "stack": [
        {"url": "/url1", "method": "POST", "data": {"sample": "json"}},
        {"url": "/url2", "method": "GET", "data": {"get": "argument"}},
        ...
    ]
}
```

Error Reporting
===============

Gottp can send error tracebacks if a request failed.
This can be enabled by setting EmailDummy as false in the gottp section of cfg.

A sample traceback email looks like this:

![Sample Traceback](http://test.simversity.com.s3.amazonaws.com/original_image/142052802943405598970493/emailtraceback.png)

TODO: Expose a way to use custom email templates.

Background Workers
==================

Gottp now supports background workers. Background workers allow parallel background execution of long running processes.

There are a few advantage of using gottp managed background worker over simply spawning a goroutine:

 * Gottp Workers are autoamtically recovered and re-spawned in case of Panic. You also get a fancy error traceback provided you have configured the email tracebacks, discussed earlier.
 * Gottp Workers do not quit abruptly when main thread receives an interrupt. This is beneficial to process cleanups and other vital exit routines.
 * Gottp Workers are provided with a default timeout of 10 seconds and a worker that seems to be taking forever is then terminated.
 * Gottp gracefully handles all process interrupts and signals the background worker accordingly.

NOTE: At this moment, gottp only supports the option of one background worker.

Usage:

```go
import "gopkg.in/meson10/gottp.v3"

func main() {
	gottp.RunWorker(func(exitChan bool) {
		log.Println("Do some background work here");
	})
}
```
