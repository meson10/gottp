package gottp

import "net/http"

//Handler is an interface that implements all the methods of an HTTP Request.
type Handler interface {
	Get(request *Request)
	Put(request *Request)
	Post(request *Request)
	Delete(request *Request)
	Head(request *Request)
	Options(request *Request)
	Patch(request *Request)
}

//BaseHandler is a type that has a name and a pattern to match.
//This type will implement the required methods of Handler interface.
type BaseHandler struct {
	name    string
	pattern string
}

func notImplemented(request *Request) {
	e := HttpError{http.StatusNotImplemented, "Method not Implemented"}
	request.Raise(e)
}

func (b *BaseHandler) Get(request *Request) {
	notImplemented(request)
}

func (b *BaseHandler) Put(request *Request) {
	notImplemented(request)
}

func (b *BaseHandler) Post(request *Request) {
	notImplemented(request)
}

func (b *BaseHandler) Delete(request *Request) {
	notImplemented(request)
}

func (b *BaseHandler) Head(request *Request) {
	notImplemented(request)
}

func (b *BaseHandler) Options(request *Request) {
	notImplemented(request)
}

func (b *BaseHandler) Patch(request *Request) {
	notImplemented(request)
}
