package handler

import (
	"net/http"
	"net/http/httputil"

	"github.com/gogf/gf/container/gmap"

	"github.com/bzp2010/webvpn/internal/model"
	"github.com/bzp2010/webvpn/internal/processor"
	"github.com/bzp2010/webvpn/internal/utils"
)

var PolicyMap = gmap.NewHashMap(false)
var reverseProxyCache = gmap.NewHashMap(false)

func ProxyHandler(response http.ResponseWriter, request *http.Request) {
	domain := request.Host
	url := request.URL

	policy, ok := PolicyMap.Get(domain).(*model.Policy)
	if !ok {
		_, _ = response.Write([]byte("policy not exist"))
		return
	}

	remote, err := url.Parse(policy.To)
	if err != nil {
		panic(err)
	}

	// reverse proxy cache
	var reverseProxy *httputil.ReverseProxy
	if reverseProxyCache.Get(domain) != nil {
		reverseProxy = reverseProxyCache.Get(domain).(*httputil.ReverseProxy)
	} else {
		reverseProxy = httputil.NewSingleHostReverseProxy(remote)
		reverseProxyCache.Set(domain, reverseProxy)
	}

	// create reverse proxy processor
	p := processor.NewManager()
	p.Use(&processor.ProxyProcessor{Policy: policy, RequestURL: request.URL})

	// modify reverse proxy request
	reverseProxy.Director = func(req *http.Request) {
		p.DoRequestProcess(req)

		utils.Log().Debugf("proxy target url: %s", req.URL.String())
	}

	// modify reverse proxy response
	reverseProxy.ModifyResponse = func(resp *http.Response) error {
		p.DoResponseProcess(resp)
		return nil
	}

	reverseProxy.ServeHTTP(response, request)
}
