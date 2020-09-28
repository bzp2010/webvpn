package processor

import (
	"github.com/bzp2010/webvpn/utils"
	"net/http"
	"strings"
)

type StaticResourceProcessor struct {
	next ProcessorInterface
}

func (h *StaticResourceProcessor) RequestProcess(r *http.Request) *http.Request {
	return r
}

func (h *StaticResourceProcessor) ResponseProcess(r *http.Response) *http.Response {
	if contentType := r.Header.Get("Content-Type"); contentType != "" {
		contentType = strings.Split(contentType, ";")[0]
		if contentType == "application/javascript" || contentType == "text/css" {
			utils.Logger.Info("静态资源处理器调用")
		}
	}

	if h.next != nil {
		h.next.ResponseProcess(r)
	}

	return r
}

func (h *StaticResourceProcessor) SetNext(p ProcessorInterface) ProcessorInterface {
	h.next = p
	return p
}
