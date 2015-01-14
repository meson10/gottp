package gottp

type allUrls struct {
	BaseHandler
}

var allUrlsMap = map[string]string{}

func (self *allUrls) Get(req *Request) {
	if len(allUrlsMap) == 0 {
		for _, url := range boundUrls {
			allUrlsMap[url.name] = url.url
		}
	}

	req.Write(allUrlsMap)
}
