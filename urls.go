package gottp

import (
	"log"
	"net/http"
	"regexp"
)

type Url struct {
	name    string
	url     string
	handler Handler
	pattern *regexp.Regexp
}

var boundUrls = []*Url{}

func NewUrl(name string, pattern string, handler Handler) {
	compiled_pattern, err := regexp.Compile(pattern)
	if err != nil {
		panic(err)
	}

	url := Url{
		name:    name,
		handler: handler,
		pattern: compiled_pattern,
		url:     pattern,
	}

	boundUrls = append(boundUrls, &url)
}

func (u Url) MakeUrlArgs(url *string) (*map[string]string, bool) {
	matches := u.pattern.FindStringSubmatch(*url)
	named_groups := u.pattern.SubexpNames()
	data := map[string]string{}

	var err bool

	if len(matches) > 0 {
		for ix, key := range named_groups {
			if len(key) > 0 {
				data[key] = matches[ix]
			}
		}
	} else if len(named_groups) > 0 {
		err = true
	}

	return &data, err
}

func bindGlobalHandler() {
	http.HandleFunc("/", GlobalHandler)
}

func bindHandlers() {
	NewUrl("async_pipe", "^/async_pipe/?$", new(AsyncPipeHandler))
	NewUrl("pipe", "^/pipe/?$", new(PipeHandler))
	NewUrl("urls", "^/urls/?$", new(UrlHandler))
}

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	bindGlobalHandler()
	bindHandlers()
}
