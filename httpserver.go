package gottp

import (
	"bytes"
	"gottp/conf"
	"log"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"
	"utils"
)

const ERROR = "An Internal Error has occured while processing this Request." +
	"If you believe this is an Error, please contact support team."

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
		conf.Settings.EmailFrom,
		conf.Settings.ErrorTo,
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
			e := HttpError{500, ERROR}
			req.Raise(e)
		} else {
			http.Error(req.Writer, ERROR, 500)
		}
	}
}

func getUrls(req *Request, urls []*Url) {
	ret := map[string]string{}
	for _, url := range urls {
		ret[url.name] = url.url
	}
	req.Write(ret)
}

type eachCall struct {
	Host   string
	Method string                 `method`
	Url    string                 `url`
	Data   map[string]interface{} `data`
}

func eachPipe(req *Request, call eachCall, urls []*Url) {
	for _, url := range urls {
		urlArgs, err := url.MakeUrlArgs(&call.Url)
		if err {
			continue
		}

		pipeData := bytes.NewBuffer(utils.Encoder(call.Data))
		pipeReq, err2 := http.NewRequest(call.Method, call.Url, pipeData)

		if err2 != nil {
			log.Println(err)
		} else {
			pipeReq.Host = call.Host
			req.Request = pipeReq
			req.UrlArgs = urlArgs
			url.handler(req)
		}
		break
	}
}

func performPipe(w http.ResponseWriter, req *http.Request, urls []*Url, async bool) {

	pipeUrls := []string{}

	parentReq := Request{Writer: w, Request: req, UrlArgs: nil}
	defer Exception(&parentReq)

	if req.Method != "POST" {
		parentReq.Write("Pipe/Async-Pipe only supports POST request.")
		return
	}

	var calls []eachCall
	parentReq.ConvertArgument("stack", &calls)

	po := make([]utils.Q, len(calls))

	var wg sync.WaitGroup
	for index, oneCall := range calls {
		oneCall.Host = req.Host
		pipeUrls = append(pipeUrls, oneCall.Url)
		pipeReq := Request{PipeOutput: &po, PipeIndex: index}

		if async {
			wg.Add(1)
			go func(call eachCall) {
				defer Exception(&parentReq)
				defer wg.Done()
				eachPipe(&pipeReq, call, urls)
			}(oneCall)
		} else {
			eachPipe(&pipeReq, oneCall, urls)
		}
	}

	defer pipeTimeTrack(time.Now(), req, strings.Join(pipeUrls, ", "))

	if async {
		wg.Wait()
	}

	parentReq.Write(po)
}

func pipeTimeTrack(start time.Time, req *http.Request, pipeUrls string) {
	elapsed := time.Since(start)
	defer log.Printf("[%s] %s %s [%s] %s\n", req.Method, req.URL, req.RemoteAddr, pipeUrls, elapsed)
}

func timeTrack(start time.Time, req *http.Request) {
	elapsed := time.Since(start)
	log.Printf("[%s] %s %s %s\n", req.Method, req.URL, req.RemoteAddr, elapsed)
}

func BindHandlers(urls []*Url) {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	http.HandleFunc("/async-pipe", func(w http.ResponseWriter, req *http.Request) {
		performPipe(w, req, urls, true)
	})

	http.HandleFunc("/pipe", func(w http.ResponseWriter, req *http.Request) {
		performPipe(w, req, urls, false)
	})

	http.HandleFunc("/urls", func(w http.ResponseWriter, req *http.Request) {
		defer timeTrack(time.Now(), req)
		p := Request{Writer: w, Request: req, UrlArgs: nil}
		defer Exception(&p)
		getUrls(&p, urls)
		return
	})

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		defer timeTrack(time.Now(), req)

		for _, url := range urls {
			urlArgs, err := url.MakeUrlArgs(&req.URL.Path)
			if !err {
				p := Request{Writer: w, Request: req, UrlArgs: urlArgs}
				defer Exception(&p)
				url.handler(&p)
				return
			}
		}

		http.NotFound(w, req)
	})
}
