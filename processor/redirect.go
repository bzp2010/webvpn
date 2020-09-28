package processor

import (
	"github.com/bzp2010/webvpn/model"
	"net/http"
	"net/url"
	"strings"
)

type RedirectProcessor struct {
	next    ProcessorInterface
	Service *model.Service
	Host    string
}

func (h *RedirectProcessor) RequestProcess(r *http.Request) *http.Request {
	return r
}

func (h *RedirectProcessor) ResponseProcess(r *http.Response) *http.Response {
	// replace target URL when backend return HTTP 301 or 302
	if location := r.Header.Get("Location"); location != "" {
		// not redirect out domain resource
		locationURL, _ := url.Parse(location)
		if locationURL.Host != h.Host {
			return r
		}

		// modify redirect url
		newURLStr := ""
		if strings.HasPrefix(location, h.Service.Url) {
			newURLStr = strings.Replace(location, h.Service.Url, "/"+h.Service.Name+"/", 1)
		} else {
			newURLStr = "/" + h.Service.Name + location
		}

		r.Header.Set("Location", newURLStr)
		//utils.Logger.Info("重定向处理器调用", zap.String("old", location),zap.String("new_url", r.Header.Get("Location")))
	}

	if h.next != nil {
		h.next.ResponseProcess(r)
	}

	return r
}

func (h *RedirectProcessor) SetNext(p ProcessorInterface) ProcessorInterface {
	h.next = p
	return p
}
