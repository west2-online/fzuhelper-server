/*
Copyright 2024 The west2-online Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package db

import (
	"context"

	"gorm.io/gorm/logger"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	logger2 "github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func InitMySQL() {
	dsn, err := utils.GetMysqlDSN()
	if err != nil {
		logger2.LoggerObj.Fatal("get mysql DSN error: " + err.Error())
	}
	DB, err = gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
	if err != nil {
		logger2.LoggerObj.Fatal("mysql connect error")
	} else {
		logger2.LoggerObj.Info("mysql connect access")
	}

	sqlDB, _ := DB.DB()
	sqlDB.SetMaxIdleConns(constants.MaxIdleConns)
	sqlDB.SetMaxOpenConns(constants.MaxConnections)
	sqlDB.SetConnMaxLifetime(constants.ConnMaxLifetime)
	DB = DB.Table(constants.UserTableName).WithContext(context.Background())
}
