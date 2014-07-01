package gottp

import "regexp"

type Url struct {
	name    string
	url     string
	handler func(r *Request)
	pattern *regexp.Regexp
}

func NewUrl(name string, pattern string, handler func(r *Request)) *Url {
	compiled_pattern, err := regexp.Compile(pattern)
	if err != nil {
		panic(err)
	}
	return &Url{name: name, handler: handler, pattern: compiled_pattern, url: pattern}
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
		err = false
	} else if len(named_groups) > 0 {
		err = true
	}

	return &data, err
}
