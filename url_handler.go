package gottp

// UrlHandler reponds with all the url exposed by the api.
type UrlHandler struct {
	BaseHandler
}

var allUrlsMap = map[string]string{}

// Get implements the method for UrlHandler.
// It would be an advantage to have methods
// implemented on urls in the response.
func (self *UrlHandler) Get(req *Request) {
	if len(allUrlsMap) == 0 {
		for _, url := range boundUrls {
			allUrlsMap[url.name] = url.url
		}
	}

	req.Write(allUrlsMap)
	return
}
