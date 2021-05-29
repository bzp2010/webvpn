package server

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Options struct {
	Public bool
	Admin  bool
}

var publicMux *chi.Mux
var database *gorm.DB

func NewServer(o *Options) error {
	// initialize database
	err := initDatabase()
	if err != nil {
		return err
	}

	// start server
	g.Log().Infof("Run WebVPN server: public: %t  admin: %t", o.Public, o.Admin)
	if o.Public {
		err := initPublicServer()
		if err != nil {
			return gerror.Newf("public server start failed: %s", err.Error())
		}
	}

	if o.Admin {
		err := initAdminServer()
		if err != nil {
			return gerror.Newf("admin server start failed: %s", err.Error())
		}
	}

	return nil
}

func initDatabase() error {
	g.Log().Info("connecting to database")

	dsnURL := viper.GetString("dsn")
	dsn := strings.Split(dsnURL, "://")

	if dsn[0] != "mysql" {
		return gerror.Newf("unsupported database software: %s", dsn[0])
	}

	var err error
	switch dsn[0] {
	case "mysql":
	default:
		database, err = gorm.Open(mysql.Open(dsn[1]), &gorm.Config{})
		break
	}

	if err != nil {
		return err
	}

	g.Log().Infof("connected to database: dsn=%s://%s", dsn[0], regexp.MustCompile(`(:).*(@)`).ReplaceAllString(dsn[1], ":****@"))

	return nil
}

func initPublicServer() error {
	publicMux = chi.NewRouter()
	publicMux.Get("/*", Handler)

	hostAddr := viper.GetString("serve.public.host") + ":" + viper.GetString("serve.public.port")

	return http.ListenAndServe(hostAddr, publicMux)
}

func initAdminServer() error {
	return nil
}