package core

import (
	"github.com/bzp2010/webvpn/model"
	"github.com/bzp2010/webvpn/processor"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"go.uber.org/zap"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func RequestHandler(r *ghttp.Request) {
	if r.RequestURI == "/favicon.ico" {
		r.Response.Write("Request path Error")
		return
	}

	// get service name from path
	serviceName := r.GetString("service")

	// split real path and query
	pathWithQuery := strings.Replace(r.RequestURI, "/"+serviceName, "", 1)

	// get service info
	service := model.GetServiceByName(serviceName)
	rawURL := service.Url + strings.Trim(pathWithQuery, "/")

	// create reverse proxy
	remote, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}

	g.Log().Info(r.Host, zap.String("test", r.RequestURI))

	handler := &processor.RedirectProcessor{Service: service, Host: r.Host}
	handler.SetNext(new(processor.StaticResourceProcessor))
	handler.SetNext(&processor.HTMLProcessor{RequestURI: r.RequestURI, Service: service})

	proxy := httputil.NewSingleHostReverseProxy(remote)

	// modify reverse proxy request
	proxy.Director = func(req *http.Request) {
		// set the request host header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = remote.Path
		req.URL.RawQuery = remote.RawQuery

		// close gzip
		req.Header.Set("Accept-Encoding", "")

		req = handler.RequestProcess(req)
	}

	// modify reverse proxy response
	proxy.ModifyResponse = func(rep *http.Response) error {
		rep = handler.ResponseProcess(rep)
		//err = rep.Body.Close()
		//if err != nil {
		//	return err
		//}

		return nil
	}

	r.URL.Path = ""

	proxy.ServeHTTP(r.Response.ResponseWriter, r.Request)
}
