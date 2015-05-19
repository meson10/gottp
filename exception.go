package gottp

import (
	"bytes"
	"html/template"

	utils "gopkg.in/simversity/gottp.v3/utils"
)

const errorExtra = `
Host: {{.Host}}
Method: {{.Method}}
Protocol: {{.Protocol}}
RemoteIP: {{.RemoteIP}}
URI: {{.URI}}
Referer: {{.Referer}}
Arguments: {{.Arguments}}
`

type ErrorStack struct {
	Host      string
	Method    string
	Protocol  string
	RemoteIP  string
	URI       string
	Referer   string
	Arguments string
}

func getTracerExtra(req *Request) func(reason string) string {
	return func(reason string) string {
		e := HttpError{500, "Internal Server Error: " + reason}
		req.Raise(e)

		if req == nil {
			return ""
		}

		var referrer string
		if len(req.Request.Header["Referer"]) > 0 {
			referrer = req.Request.Header["Referer"][0]
		}

		params := string(utils.Encoder(req.GetArguments()))

		stack := ErrorStack{
			Host:      req.Request.Host,
			Method:    req.Request.Method,
			Protocol:  req.Request.Proto,
			RemoteIP:  req.Request.RemoteAddr,
			URI:       req.Request.URL.Path,
			Referer:   referrer,
			Arguments: params,
		}

		var doc bytes.Buffer

		t := template.Must(template.New("error_email").Parse(errorExtra))
		t.Execute(&doc, stack)

		return doc.String()
	}
}
