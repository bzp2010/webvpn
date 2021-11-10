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
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"github.com/undefinedlabs/go-mpatch"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func checkpatchErr(t *testing.T, err error) {
	if err != nil {
		t.Errorf("patch error: %v", err)
	}
}

func patchnewdatabase(t *testing.T, returnError error) *mpatch.Patch {
	var handler *mpatch.Patch
	var err error
	handler, err = mpatch.PatchMethod(newDatabase, func() (*gorm.DB, error) {
		return nil, returnError
	})

	checkpatchErr(t, err)
	return handler
}

func initdatabase() {
	database = &gorm.DB{}
}

func Test_Database(t *testing.T) {
	var p1 *mpatch.Patch
	cases := []struct {
		name    string
		error   error
		before  func()
		finally func()
	}{
		{
			name:  "Database already initialized",
			error: nil,
			before: func() {
				initdatabase()
			},
			finally: func() {
				database = nil
			},
		},
		{
			name:  "New database failed",
			error: errors.New("create_db_error"),
			before: func() {
				p1 = patchnewdatabase(t, errors.New("create_db_error"))
			},
			finally: func() {
				_ = p1.Unpatch()
			},
		},
	}

	for _, testcase := range cases {

		// en: pre_execute patch code
		if testcase.before != nil {
			testcase.before()
		}

		// en: run test code
		_, err := Database()

		if fmt.Sprint(err) != fmt.Sprint(testcase.error) {
			t.Errorf("Caught unexpected error: %v", err)
		}

		// finally
		if testcase.finally != nil {
			testcase.finally()
		}

	}
}

func patchvipergetstring(t *testing.T, ret string) *mpatch.Patch {
	var handler *mpatch.Patch
	var err error
	handler, err = mpatch.PatchMethod(viper.GetString, func(key string) string {
		return ret
	})

	checkpatchErr(t, err)

	return handler
}

func patchmysqlopendatabase(t *testing.T) *mpatch.Patch {
	var handler *mpatch.Patch
	var err error

	handler, err = mpatch.PatchMethod(mysql.Open, func(dsn string) gorm.Dialector {
		return nil
	})

	checkpatchErr(t, err)

	return handler
}

func patchgormopendatabase(t *testing.T, retErr error) *mpatch.Patch {
	var handler *mpatch.Patch
	var err error

	handler, err = mpatch.PatchMethod(gorm.Open, func(dialector gorm.Dialector, opts ...gorm.Option) (db *gorm.DB, err error) {
		_ = handler.Unpatch()
		return nil, retErr
	})

	checkpatchErr(t, err)

	return handler
}

func Test_newDatabase(t *testing.T) {
	var p1, p2, p3 *mpatch.Patch
	cases := []struct {
		name     string
		error    error
		expectDb *gorm.DB
		before   func()
		finally  func()
	}{
		{
			name:     "Unsupported database engine",
			error:    errors.New("unsupported database engine: enginex"),
			expectDb: nil,
			before: func() {
				p1 = patchvipergetstring(t, "enginex://demouser:demopassword@demohost")
			},
			finally: func() {
				_ = p1.Unpatch()
			},
		},
		{
			name:     "database is mysql but open failed",
			error:    errors.New("connect to database failed: dial_db_failed"),
			expectDb: nil,
			before: func() {
				p1 = patchvipergetstring(t, "mysql://demouser:demopassword@demohost")
				p2 = patchmysqlopendatabase(t)
				p3 = patchgormopendatabase(t, errors.New("dial_db_failed"))
			},
			finally: func() {
				_ = p1.Unpatch()
				_ = p2.Unpatch()
				_ = p3.Unpatch()
			},
		},
		{
			name:     "database is mysql and open succeed",
			error:    nil,
			expectDb: nil,
			before: func() {
				p1 = patchvipergetstring(t, "mysql://demouser:demopassword@demohost")
				p2 = patchmysqlopendatabase(t)
				p3 = patchgormopendatabase(t, nil)
			},
			finally: func() {
				_ = p1.Unpatch()
				_ = p2.Unpatch()
				_ = p3.Unpatch()
			},
		},
	}

	for _, testcase := range cases {
		if testcase.before != nil {
			testcase.before()
		}

		db, err := newDatabase()

		if testcase.expectDb == nil {
			if db != nil {
				t.Errorf("Arg expect_db should be nil, but this isn't the case.")
			}
		} else {
			if db == nil {
				t.Errorf("Arg expect_db should not be nil, but it is actually nil.")
			}
		}

		if fmt.Sprint(err) != fmt.Sprint(testcase.error) {
			t.Errorf("Caught unexpected error: %v", err)
		}

		if testcase.finally != nil {
			testcase.finally()
		}
	}
}
