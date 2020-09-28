package processor

import (
	"bytes"
	"github.com/bzp2010/webvpn/model"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type HTMLProcessor struct {
	next       ProcessorInterface
	Service    *model.Service
	RequestURI string
}

func (h *HTMLProcessor) RequestProcess(r *http.Request) *http.Request {
	return r
}

func (h *HTMLProcessor) ResponseProcess(r *http.Response) *http.Response {
	if contentType := r.Header.Get("Content-Type"); contentType != "" {
		contentType = strings.Split(contentType, ";")[0]
		if contentType == "text/html" {

			// change resource link in html
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return r
			}

			// replace absolute resource import
			body = bytes.ReplaceAll(body, []byte("href=\"/"), []byte("href=\"/"+h.Service.Name+"/"))
			body = bytes.ReplaceAll(body, []byte("src=\"/"), []byte("src=\"/"+h.Service.Name+"/"))

			// replace form action
			body = bytes.ReplaceAll(body, []byte("action=\"/"), []byte("action=\"/"+h.Service.Name+"/"))

			// insert ServiceWorker loader
			attendHTML := "<script>window.WebVPNSWConfig = {path: \"" + h.RequestURI + "\", service: {name: \"" + h.Service.Name + "\"}};\n(function () {if ('serviceWorker' in navigator) {console.log('sw激活调用');navigator.serviceWorker.register('/sw.js');}})();\nnavigator.serviceWorker.controller.postMessage(window.WebVPNSWConfig);</script>"
			body = bytes.ReplaceAll(body, []byte("</body>"), []byte(attendHTML+"</body>"))

			r.Body = ioutil.NopCloser(bytes.NewReader(body))
			r.ContentLength = int64(len(body))
			r.Header.Set("Content-Length", strconv.Itoa(len(body)))
		}
	}

	if h.next != nil {
		h.next.ResponseProcess(r)
	}

	return r
}

func (h *HTMLProcessor) SetNext(p ProcessorInterface) ProcessorInterface {
	h.next = p
	return p
}
