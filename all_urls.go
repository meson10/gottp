package gottp

type allUrls struct {
	BaseHandler
}

func (self *allUrls) Get(req *Request) {
	ret := map[string]string{}
	for _, url := range boundUrls {
		ret[url.name] = url.url
	}

	req.Write(ret)
}
