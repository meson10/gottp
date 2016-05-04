package gottp

import (
	"compress/gzip"
	"compress/zlib"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gopkg.in/simversity/gottp.v3/utils"
)

//*http.Request, rw ResponseWriter

const PAGE_SIZE = 30
const SKIP = 0
const serverUA = "Gottp Server"
const ERROR = `Oops! An Internal Error occured while performing that action.
Please try again later`

type WireSender interface {
	SendOverWire() utils.Q
}

type Paginator struct {
	Skip  int
	Limit int
	Wlt   string
	Wgt   string
	Wkey  string
	Ids   []string
}

func makeInt(val interface{}, fallback int) int {
	switch val.(type) {
	case int:
		return val.(int)
	case string:
		_val := val.(string)
		ret, err := strconv.Atoi(_val)
		if err == nil {
			return ret
		}
	case float64:
		return int(val.(float64))
	}
	return fallback
}

func makeString(val interface{}) (ret string) {
	switch val.(type) {
	case string:
		ret = val.(string)
	case int:
		val := val.(int)
		ret = strconv.Itoa(val)
	}
	return
}

// Request type is a gottp wrapper of the incoming request and its 
// response. It has data structures for url arguments or parameters.
type Request struct {
	Request    *http.Request
	Writer     http.ResponseWriter
	UrlArgs    *map[string]string
	pipeOutput chan<- *utils.Q
	pipeIndex  int
	params     *utils.Q
}

// makeUrlArgs returns map of url arguments.
func (r *Request) makeUrlArgs(params utils.Q) {
	if r.UrlArgs != nil {
		for key, value := range *r.UrlArgs {
			params[key] = value
		}
	}
}

// GetPaginator returns pointer to Paginator type with appropriate values 
// assigned.
func (r *Request) GetPaginator() *Paginator {
	p := Paginator{Limit: -1}
	qp := r.GetArguments()
	for key, value := range *qp {
		switch key {
		case "skip":
			p.Skip = makeInt(value, SKIP)
		case "limit":
			p.Limit = makeInt(value, PAGE_SIZE)
		case "wlt":
			p.Wlt = makeString(value)
		case "wgt":
			p.Wgt = makeString(value)
		case "wkey":
			p.Wkey = makeString(value)
		case "ids":
			ids, ok := value.([]string)
			if !ok {
				id := value.(string)
				p.Ids = []string{id}
			} else {
				p.Ids = ids
			}
		}
	}

	if p.Limit < 0 {
		p.Limit = PAGE_SIZE
	}

	return &p
}

// makeUrlParams returns map of url parameters.
func (r *Request) makeUrlParams(params utils.Q) {
	r.Request.ParseForm()
	for key, value := range r.Request.Form {
		length := len(value)
		if length == 1 {
			params[key] = value[0]
		} else {
			params[key] = value
		}
	}
}

// makeBodyParams returns map of body parameters.
func (r *Request) makeBodyParams(params utils.Q) {
	//if (r.Request.Method == "PUT" || r.Request.Method == "POST") {
	if r.Request.ContentLength != 0 {
		utils.DecodeStream(r.Request.Body, &params)
	}
}

// GetArguments returns pointer to request's param which
// is a map of all the parameters/arguments in the url/body.
func (r *Request) GetArguments() *utils.Q {
	if r.params == nil {
		params := utils.Q{}
		r.makeUrlArgs(params)
		r.makeBodyParams(params)
		r.makeUrlParams(params)
		r.params = &params
	}

	return r.params
}

// ConvertArguments converts the request arguments to correspondig golang data structures
// and stores it into f.
func (r *Request) ConvertArguments(f interface{}) {
	utils.Convert(r.GetArguments(), f)
}

// GetArgument returns the argument value matching the key.
func (r *Request) GetArgument(key string) interface{} {
	args := *r.GetArguments()
	return args[key]
}

// ConvertArgument converts the value of argument matching key to correspondig golang data
// structures and stores it into f.
func (r *Request) ConvertArgument(key string, f interface{}) {
	args := *r.GetArguments()
	val := args[key]
	utils.Convert(val, f)
}

// Finish adds three headers, and returns the corresponding JSON of data.
func (r *Request) Finish(data interface{}) []byte {
	r.Writer.Header().Set("Server", serverUA)
	r.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	r.Writer.Header().Set("Content-Type", "application/json")
	return utils.Encoder(data)
}

// Redirect redirects the user to given url.
func (r *Request) Redirect(url string, status int) {
	log.Println("Redirecting to", url)
	http.Redirect(r.Writer, r.Request, url, status)
	return
}

// Write writes and finishes the request with apporpiate response.
func (r *Request) Write(data interface{}) {
	var piped utils.Q

	if v, ok := data.(WireSender); ok {
		piped = v.SendOverWire()
	} else {
		piped = utils.Q{
			"data":    data,
			"status":  http.StatusOK,
			// why is there a need for message field in 
			// the response when theres no way to change it?
			"message": "",
		}
	}

	if r.pipeOutput != nil {
		piped["index"] = r.pipeIndex
		r.pipeOutput <- &piped
	} else if strings.Contains(
		r.Request.Header.Get("Accept-Encoding"), "deflate",
	) {
		r.Writer.Header().Set("Content-Encoding", "deflate")

		gz := zlib.NewWriter(r.Writer)
		defer gz.Close()
		gz.Write(r.Finish(piped))

	} else if strings.Contains(
		r.Request.Header.Get("Accept-Encoding"), "gzip",
	) {
		r.Writer.Header().Set("Content-Encoding", "gzip")

		gz := gzip.NewWriter(r.Writer)
		defer gz.Close()
		gz.Write(r.Finish(piped))

	} else {
		r.Writer.Write(r.Finish(piped))
	}
}

// Raise sends HttpError as response 
func (r *Request) Raise(e HttpError) {
	r.Write(e)
}

// performRequest takes a handler and a pointer to request and 
// calls the appropriate handler defined by the user.
func performRequest(handler Handler, p *Request) {
	method := (*p).Request.Method

	switch method {
	case "GET":
		handler.Get(p)
	case "POST":
		handler.Post(p)
	case "PUT":
		handler.Put(p)
	case "DELETE":
		handler.Delete(p)
	case "HEAD":
		handler.Head(p)
	case "OPTIONS":
		handler.Options(p)
	case "PATCH":
		handler.Patch(p)
	default:
		log.Println("Unsupported method", method)
	}
}

func doRequest(request *Request, availableUrls *[]*Url) {
	requestUrl := request.Request.URL.Path
	defer Tracer.Notify(getTracerExtra(request))

	for _, url := range *availableUrls {
		urlArgs, err := url.MakeUrlArgs(&requestUrl)
		if !err {
			request.UrlArgs = urlArgs
			performRequest(url, request)
			return
		}
	}

	e := HttpError{404, requestUrl + " did not match any exposed urls."}
	request.Raise(e)
	return
}
