package gottp

import (
	"bytes"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	utils "github.com/Simversity/gottp/utils"
)

const serverUA = "Gottp Server"
const ERROR = `Oops! An Internal Error occured while performing that action.
Please try again later`

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

	done := make(chan bool)
	output := make([]*utils.Q, len(calls))
	po := make(chan *utils.Q, len(calls))

	var wg sync.WaitGroup

	go func(sink <-chan *utils.Q) {
		for {
			v, more := <-sink
			if more {
				if index, ok := (*v)["index"].(int); ok {
					output[index] = v
				}
			} else {
				done <- true
				return
			}
		}

	}(po)

	for index, oneCall := range calls {
		oneCall.Host = req.Host
		pipeUrls = append(pipeUrls, oneCall.Url)
		pipeReq := Request{pipeOutput: po, pipeIndex: index}

		wg.Add(1)

		if async {
			go func(call eachCall) {
				defer Exception(&pipeReq)
				defer wg.Done()
				eachPipe(&pipeReq, call, urls)
			}(oneCall)
		} else {
			func(call eachCall) {
				defer Exception(&pipeReq)
				defer wg.Done()
				eachPipe(&pipeReq, call, urls)
			}(oneCall)
		}
	}

	defer pipeTimeTrack(time.Now(), req, strings.Join(pipeUrls, ", "))

	wg.Wait()
	close(po)

	<-done
	parentReq.Write(output)
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
