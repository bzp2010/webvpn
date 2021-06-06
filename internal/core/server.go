package core

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"

	"github.com/bzp2010/webvpn/internal/handler"
	"github.com/bzp2010/webvpn/internal/utils"
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
	utils.Log().Infof("create webvpn server: public: %t admin: %t", o.Public, o.Admin)

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
	r.Get("/*", handler.ProxyHandler)

	return r
}

func initAdminServer() *chi.Mux {
	return nil
}

func (s *WebVPN) Start() {
	if s.options.Public {
		hostAddr := viper.GetString("serve.public.host") + ":" + viper.GetString("serve.public.port")
		go func() {
			err := http.ListenAndServe(hostAddr, s.publicMux)
			if err != nil {
				utils.Log().Errorf("webvpn server: public start failed: %s", err.Error())
				panic("webvpn server: public start failed: " + err.Error())
			}
		}()
		utils.Log().Infof("webvpn server: public run on %s", hostAddr)
	}

	if s.options.Admin {
		hostAddr := viper.GetString("serve.admin.host") + ":" + viper.GetString("serve.admin.port")
		go func() {
			err := http.ListenAndServe(hostAddr, s.adminMux)
			if err != nil {
				utils.Log().Errorf("webvpn server: admin start failed: %s", err.Error())
				panic("webvpn server: admin start failed: " + err.Error())
			}
		}()
		utils.Log().Infof("webvpn server: admin run on %s", hostAddr)
	}
}
