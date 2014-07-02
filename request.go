package gottp

import (
	"net/http"
	"strconv"
	"strings"
	//"compress/gzip"
	utils "github.com/Simversity/gottp/utils"
)

//*http.Request, rw ResponseWriter

const PAGE_SIZE = 30
const SKIP = 0

type WireSender interface {
	SendOverWire() utils.Q
}

type HttpError struct {
	Status  int
	Message string
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

func (e HttpError) SendOverWire() utils.Q {
	if e.Status == 0 {
		e.Status = 500
	}

	if len(e.Message) == 0 {
		e.Message = "An Internal exception has occured. Please try again after some time"
	}

	return utils.Q{
		"data":    nil,
		"status":  e.Status,
		"message": e.Message,
	}
}

type Request struct {
	Request    *http.Request
	Writer     http.ResponseWriter
	UrlArgs    *map[string]string
	PipeOutput *[]utils.Q
	PipeIndex  int
	params     *utils.Q
}

func (r *Request) makeUrlArgs() {
	original := *r.params
	if r.UrlArgs != nil {
		for key, value := range *r.UrlArgs {
			original[key] = value
		}
	}
}

func (r *Request) GetPaginator() *Paginator {
	p := Paginator{Limit: -1}
	qp := r.GetArguments()
	for key, value := range *qp {
		switch key {
		case "skip":
			p.Skip = makeInt(value, PAGE_SIZE)
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

func (r *Request) makeUrlParams() {
	original := *r.params
	r.Request.ParseForm()
	for key, value := range r.Request.Form {
		length := len(value)
		if length == 1 {
			original[key] = value[0]
		} else {
			original[key] = value
		}
	}
}

func (r *Request) makeBodyParams() {
	//if (r.Request.Method == "PUT" || r.Request.Method == "POST") {
	if r.Request.ContentLength != 0 {
		utils.DecodeStream(r.Request.Body, r.params)
	}
}

func (r *Request) GetArguments() *utils.Q {
	if r.params == nil {
		r.params = &utils.Q{}
		r.makeUrlArgs()
		r.makeBodyParams()
		r.makeUrlParams()
	}

	return r.params
}

func (r *Request) ConvertArguments(f interface{}) {
	utils.Convert(r.GetArguments(), f)
}

func (r *Request) GetArgument(key string) interface{} {
	args := *r.GetArguments()
	return args[key]
}

func (r *Request) ConvertArgument(key string, f interface{}) {
	args := *r.GetArguments()
	val := args[key]
	utils.Convert(val, f)
}

func (r *Request) Write(data interface{}) {
	var piped utils.Q
	if v, ok := data.(WireSender); ok {
		piped = v.SendOverWire()
	} else {
		piped = utils.Q{
			"data":    data,
			"status":  200,
			"message": "",
		}
	}

	if r.PipeOutput != nil {
		(*r.PipeOutput)[r.PipeIndex] = piped
	} else if strings.Contains(r.Request.Header.Get("Accept-Encoding"), "gzip") {
		r.Writer.Write(utils.Encoder(piped))
	} else {
		r.Writer.Write(utils.Encoder(piped))
	}
}

func (r *Request) Raise(e HttpError) {
	r.Write(e)
}
