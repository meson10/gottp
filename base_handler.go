package gottp

import "net/http"

type Handler interface {
	Get(request *Request)
	Put(request *Request)
	Post(request *Request)
	Delete(request *Request)
	Head(request *Request)
	Options(request *Request)
	Patch(request *Request)
}

type BaseHandler struct {
	name    string
	pattern string
}

func notImplemented(request *Request) {
	e := HttpError{http.StatusNotImplemented, "Method not Implemented"}
	request.Raise(e)
}

func (self *BaseHandler) Get(request *Request) {
	notImplemented(request)
}

func (self *BaseHandler) Put(request *Request) {
	notImplemented(request)
}

func (self *BaseHandler) Post(request *Request) {
	notImplemented(request)
}

func (self *BaseHandler) Delete(request *Request) {
	notImplemented(request)
}

func (self *BaseHandler) Head(request *Request) {
	notImplemented(request)
}

func (self *BaseHandler) Options(request *Request) {
	notImplemented(request)
}

func (self *BaseHandler) Patch(request *Request) {
	notImplemented(request)
}
