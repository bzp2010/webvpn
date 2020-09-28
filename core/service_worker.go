package core

import (
	_ "github.com/bzp2010/webvpn/statik"
	"github.com/rakyll/statik/fs"
	"net/http"
)

func ServiceWorkerHandler(w http.ResponseWriter, r *http.Request) {
	vfs, _ := fs.New()

	fileServer := http.FileServer(vfs)
	fileServer.ServeHTTP(w, r)
}
