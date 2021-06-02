package core

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gogf/gf/errors/gerror"
	"github.com/spf13/viper"
)

type Options struct {
	Public bool
	Admin  bool
}

type WebVPN struct {
	options   *Options
	publicMux *chi.Mux
	adminMux  *chi.Mux
}

var server *WebVPN

func Server(o *Options) (*WebVPN, error) {
	if server != nil {
		return server, nil
	}
	return newServer(o)
}

func newServer(o *Options) (*WebVPN, error) {
	Log().Infof("create webvpn server: public: %t  admin: %t", o.Public, o.Admin)

	// create server
	s := &WebVPN{options: o}

	// init server mux
	if o.Public {
		s.publicMux = initPublicServer()
	}
	if o.Admin {
		s.adminMux = initAdminServer()
	}

	server = s

	return s, nil
}

func initPublicServer() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/*", Handler)

	return r
}

func initAdminServer() *chi.Mux {
	return nil
}

func (s *WebVPN) Start() error {
	if s.options.Public {
		hostAddr := viper.GetString("serve.public.host") + ":" + viper.GetString("serve.public.port")
		err := http.ListenAndServe(hostAddr, s.publicMux)
		if err != nil {
			return gerror.Newf("public server start failed: %s", err.Error())
		}
	}

	if s.options.Admin {
		hostAddr := viper.GetString("serve.admin.host") + ":" + viper.GetString("serve.admin.port")
		err := http.ListenAndServe(hostAddr, s.adminMux)
		if err != nil {
			return gerror.Newf("admin server start failed: %s", err.Error())
		}
	}

	return nil
}
