package gottp

type UrlHandler struct {
	BaseHandler
}

var allUrlsMap = map[string]string{}

func (self *UrlHandler) Get(req *Request) {
	if len(allUrlsMap) == 0 {
		for _, url := range boundUrls {
			allUrlsMap[url.name] = url.url
		}
	}

	req.Write(allUrlsMap)
	return
}
