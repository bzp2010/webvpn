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
	"regexp"
	"strings"

	"github.com/gogf/gf/errors/gerror"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/bzp2010/webvpn/internal/utils"
)

var database *gorm.DB

func Database() (*gorm.DB, error) {
	if database != nil {
		return database, nil
	}
	return newDatabase()
}

func newDatabase() (*gorm.DB, error) {
	utils.Log().Info("connecting to database")

	dsnURL := viper.GetString("dsn")
	dsn := strings.Split(dsnURL, "://")

	if dsn[0] != "mysql" {
		return nil, gerror.Newf("unsupported database engine: %s", dsn[0])
	}

	db, err := gorm.Open(mysql.Open(dsn[1]), &gorm.Config{})
	if err != nil {
		return nil, gerror.Newf("connect to database failed: %s", err.Error())
	}

	utils.Log().Infof("connected to database: dsn=%s://%s", dsn[0], regexp.MustCompile(`(:).*(@)`).ReplaceAllString(dsn[1], ":****@"))

	database = db

	return db, nil
}