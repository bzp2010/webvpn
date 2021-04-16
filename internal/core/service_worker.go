package core

import (
	_ "github.com/bzp2010/webvpn/statik"
	"github.com/gogf/gf/net/ghttp"
	"github.com/rakyll/statik/fs"
	"net/http"
)

func ServiceWorkerHandler(r *ghttp.Request) {
	vfs, _ := fs.New()

	fileServer := http.FileServer(vfs)
	fileServer.ServeHTTP(r.Response.ResponseWriter, r.Request)
}
