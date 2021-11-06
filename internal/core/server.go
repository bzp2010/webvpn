/*
 * Copyright (C) 2021
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package core

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	m "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/hostrouter"
	"github.com/spf13/viper"

	"github.com/bzp2010/webvpn/internal/handler"
	"github.com/bzp2010/webvpn/internal/middleware"
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
	rootRouter := chi.NewRouter()
	authenticationRouter := chi.NewRouter()
	proxyRouter := chi.NewRouter()
	hostRouter := hostrouter.New()

	// setup middleware
	authentication := &middleware.Authentication{}
	proxyRouter.Use(authentication.Authenticate, m.Logger)

	// setup authentication route
	authenticationRouter.Get("/authentication/login", handler.Login)
	hostRouter.Map(viper.GetString("authentication.service_url"), authenticationRouter)

	// setup reverse proxy route
	proxyRouter.Get("/*", handler.ProxyHandler)
	proxyRouter.Get("/.webvpn/callback", handler.Login)
	hostRouter.Map("*", proxyRouter)

	rootRouter.Mount("/", hostRouter)

	return rootRouter
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
