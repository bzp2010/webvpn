//go:generate statik -src=./static -f
package main

import (
	"github.com/bzp2010/webvpn/core"
	"github.com/bzp2010/webvpn/utils"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	// init logger
	utils.InitLogger()

	//
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/sw.js", core.ServiceWorkerHandler)
	r.Get("/{service}/*", core.RequestHandler)
	r.Post("/{service}/*", core.RequestHandler)
	r.Head("/{service}/*", core.RequestHandler)
	r.Options("/{service}/*", core.RequestHandler)
	err := http.ListenAndServe(":18085", r)
	if err != nil {
		utils.Logger.Error("WebVPN Server Starting Failed", zap.Error(err))
	}
}
