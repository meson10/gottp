package gottp

import (
	"bytes"
	"net/http"
	"strings"
	"sync"
	"time"

	"gopkg.in/simversity/gottp.v3/utils"
)

type eachCall struct {
	Host   string
	Method string  `method`
	Url    string  `url`
	Data   utils.Q `data`
}

func eachPipe(request *Request) {
	doRequest(request, &boundUrls)
	return
}

func performPipe(parentReq *Request, async bool) {

	pipeUrls := []string{}

	req := parentReq.Request
	defer Tracer.Notify(getTracerExtra(parentReq))

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

		pipeData := bytes.NewBuffer(utils.Encoder(oneCall.Data))
		pipeHttpRequest, err := http.NewRequest(oneCall.Method, oneCall.Url, pipeData)
		if err != nil {
			e := HttpError{500, oneCall.Url + " had error"}
			parentReq.Raise(e)
			return
		}

		pipeReq := Request{
			Request:    pipeHttpRequest,
			pipeOutput: po,
			pipeIndex:  index,
		}

		wg.Add(1)

		if async {
			go func(call eachCall) {
				defer wg.Done()
				eachPipe(&pipeReq)
			}(oneCall)
		} else {
			func(call eachCall) {
				defer wg.Done()
				eachPipe(&pipeReq)
			}(oneCall)
		}
	}

	defer pipeTimeTrack(time.Now(), req, strings.Join(pipeUrls, ", "))

	wg.Wait()
	close(po)

	<-done
	parentReq.Write(output)
}

type PipeHandler struct {
	BaseHandler
}

type AsyncPipeHandler struct {
	BaseHandler
}

func (self *PipeHandler) Post(req *Request) {
	performPipe(req, false)
}

func (self *AsyncPipeHandler) Post(req *Request) {
	performPipe(req, true)
}
