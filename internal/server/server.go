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

type server struct {
	options   *Options
	publicMux *chi.Mux
	adminMux  *chi.Mux
	DB        *gorm.DB
}

func NewServer(o *Options) (*server, error) {
	g.Log().Infof("create webvpn server: public: %t  admin: %t", o.Public, o.Admin)

	// create server
	s := &server{options: o}

	// initialize database
	db, err := initDatabase()
	if err != nil {
		return nil, err
	}
	s.DB = db

	// init server mux
	if o.Public {
		s.publicMux = initPublicServer()
	}
	if o.Admin {
		s.adminMux = initAdminServer()
	}

	return s, nil
}

func initDatabase() (*gorm.DB, error) {
	g.Log().Info("connecting to database")

	dsnURL := viper.GetString("dsn")
	dsn := strings.Split(dsnURL, "://")

	if dsn[0] != "mysql" {
		return nil, gerror.Newf("unsupported database engine: %s", dsn[0])
	}

	db, err := gorm.Open(mysql.Open(dsn[1]), &gorm.Config{})
	if err != nil {
		return nil, gerror.Newf("connect to database failed: %s", err.Error())
	}

	/*err = db.AutoMigrate(&model.Service{})
	if err != nil {
		return nil, err
	}*/

	g.Log().Infof("connected to database: dsn=%s://%s", dsn[0], regexp.MustCompile(`(:).*(@)`).ReplaceAllString(dsn[1], ":****@"))

	return db, nil
}

func initPublicServer() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/*", Handler)

	return r
}

func initAdminServer() *chi.Mux {
	return nil
}

func (s *server) Start() error {
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
