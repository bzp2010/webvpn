package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
)

type Options struct {
	Public bool
	Admin  bool
}

func NewServer(o *Options) error {
	if o.Public {
		return initPublicServer()
	}

	if o.Admin {
		return initAdminServer()
	}

	return nil
}

func initPublicServer() error {
	r := chi.NewRouter()
	r.Get("/*", Handler)

	hostAddr := viper.GetString("serve.public.host") + ":" + viper.GetString("serve.public.port")

	return http.ListenAndServe(hostAddr, r)
}

func initAdminServer() error {
	return nil
}