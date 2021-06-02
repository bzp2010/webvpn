package core

import (
	"regexp"
	"strings"

	"github.com/gogf/gf/errors/gerror"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var database *gorm.DB

func Database() (*gorm.DB, error) {
	if database != nil {
		return database, nil
	}
	return newDatabase()
}

func newDatabase() (*gorm.DB, error) {
	Log().Info("connecting to database")

	dsnURL := viper.GetString("dsn")
	dsn := strings.Split(dsnURL, "://")

	if dsn[0] != "mysql" {
		return nil, gerror.Newf("unsupported database engine: %s", dsn[0])
	}

	db, err := gorm.Open(mysql.Open(dsn[1]), &gorm.Config{})
	if err != nil {
		return nil, gerror.Newf("connect to database failed: %s", err.Error())
	}

	Log().Infof("connected to database: dsn=%s://%s", dsn[0], regexp.MustCompile(`(:).*(@)`).ReplaceAllString(dsn[1], ":****@"))

	database = db

	return db, nil
}