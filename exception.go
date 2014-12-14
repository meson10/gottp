package gottp

import (
	"log"
	"net/http"
	"runtime"
	"strings"

	utils "github.com/Simversity/gottp/utils"
)

func make_tag(name string, value string) string {
	return "<b>" + name + "</b>" + ": " + value
}

func requestStack(req *Request) string {
	var referrer string
	if len(req.Request.Header["Referer"]) > 0 {
		referrer = req.Request.Header["Referer"][0]
	}

	params := string(utils.Encoder(req.GetArguments()))

	return strings.Join([]string{
		make_tag("Host", req.Request.Host),
		make_tag("Method", req.Request.Method),
		make_tag("Protocol", req.Request.Proto),
		make_tag("RemoteIP", req.Request.RemoteAddr),
		make_tag("URI", req.Request.URL.Path),
		make_tag("Referer", referrer),
		make_tag("Arguments", params),
	}, "<br/>")
}

func sendException(err interface{}, stack, buffer string) {
	var exc_message string

	switch err.(type) {
	case string:
		_err, ok := err.(string)
		if ok == true {
			exc_message = _err
		}
	case interface{}:
		_err, ok := err.(error)
		if ok == true {
			exc_message = _err.Error()
		}
	}

	log.Println(err, buffer)
	connection := MakeConn()
	connection.SenderName += " Exception"

	go connection.SendEmail(Message{
		settings.Gottp.EmailFrom,
		settings.Gottp.ErrorTo,
		exc_message,
		ErrorTemplate(stack, buffer),
	})

}

func Exception(req *Request) {
	if err := recover(); err != nil {
		const size = 4096
		buf := make([]byte, size)
		buf = buf[:runtime.Stack(buf, false)]

		buffer := string(buf)
		if req == nil {
			sendException(err, "Error was raised in a goroutine", buffer)
			return
		}

		var is_json bool
		contentType := req.Request.Header["Content-Type"]
		if len(contentType) > 0 && strings.Index(contentType[0], "json") > -1 {
			is_json = true
		}

		defer sendException(err, requestStack(req), buffer)

		if is_json {
			e := HttpError{500, err.(string)}
			req.Raise(e)
		} else {
			http.Error(req.Writer, err.(string), 500)
		}
	}
}
