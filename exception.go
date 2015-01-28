package gottp

import (
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"

	utils "gopkg.in/simversity/gottp.v1/utils"
)

const layout = "Jan 2, 2006 at 3:04pm (MST)"

func sendException(err interface{}, stack *ErrorStack) {
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

	log.Println(err, stack.Traceback)

	connection := MakeConn()
	connection.SenderName += " Exception"

	go connection.SendEmail(Message{
		settings.Gottp.EmailFrom,
		settings.Gottp.ErrorTo,
		exc_message,
		ErrorTemplate(stack),
	})

}

type ErrorStack struct {
	Host      string
	Method    string
	Protocol  string
	RemoteIP  string
	URI       string
	Referer   string
	Arguments string
	Traceback string
	Timestamp string
}

func getStack(req *Request, traceback string) *ErrorStack {
	if req == nil {
		return &ErrorStack{
			Traceback: traceback,
			Timestamp: time.Now().Format(layout),
		}
	}

	var referrer string
	if len(req.Request.Header["Referer"]) > 0 {
		referrer = req.Request.Header["Referer"][0]
	}

	params := string(utils.Encoder(req.GetArguments()))

	return &ErrorStack{
		Host:      req.Request.Host,
		Method:    req.Request.Method,
		Protocol:  req.Request.Proto,
		RemoteIP:  req.Request.RemoteAddr,
		URI:       req.Request.URL.Path,
		Referer:   referrer,
		Arguments: params,
		Traceback: traceback,
		Timestamp: time.Now().Format(layout),
	}
}

func Exception(req *Request) {
	if err := recover(); err != nil {
		const size = 4096

		buf := make([]byte, size)
		buf = buf[:runtime.Stack(buf, false)]
		buffer := string(buf)

		stack := getStack(req, buffer)
		defer sendException(err, stack)

		if req == nil {
			return
		}

		var is_json bool
		contentType := req.Request.Header["Content-Type"]
		if len(contentType) > 0 && strings.Index(contentType[0], "json") > -1 {
			is_json = true
		}

		errorCode := http.StatusInternalServerError
		if is_json {
			e := HttpError{errorCode, err.(string)}
			req.Raise(e)
		} else {
			http.Error(req.Writer, err.(string), errorCode)
		}
	}
}
