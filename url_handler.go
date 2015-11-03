package gottp

// UrlHandler reponds with all the url exposed by the api.
type UrlHandler struct {
	BaseHandler
}

var allUrlsMap = make(map[string]interface{})

// Get implements the method for UrlHandler.
// It would be an advantage to have methods
// implemented on urls in the response.
func (self *UrlHandler) Get(req *Request) {
	if len(allUrlsMap) == 0 {
		for _, url := range boundUrls {
			temp := make(map[string]interface{})
			temp["url"] = url.url
			temp["methodsImplemented"] = url.methodsImplemented
			allUrlsMap[url.name] = temp
		}
	}
	req.Write(allUrlsMap)
	return
}
