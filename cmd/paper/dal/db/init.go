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
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/utils"

	gormopentracing "gorm.io/plugin/opentracing"
)

var (
	DB *gorm.DB
	SF *utils.Snowflake
)

func Init() {
	var err error

	dsn, err := utils.GetMysqlDSN()
	if err != nil {
		panic(err)
	}

	DB, err = gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,                                // 禁用默认事务
			Logger:                 logger.Default.LogMode(logger.Info), // 设置日志模式
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // 使用单数表名
			},
		})
	if err != nil {
		panic(err)
	}

	if err = DB.Use(gormopentracing.New()); err != nil {
		panic(err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(constants.MaxIdleConns)       // 最大闲置连接数
	sqlDB.SetMaxOpenConns(constants.MaxConnections)     // 最大连接数
	sqlDB.SetConnMaxLifetime(constants.ConnMaxLifetime) // 最大可复用时间

	DB = DB.Table(constants.CheckServiceTableName)

	if SF, err = utils.NewSnowflake(config.Snowflake.DatancenterID, config.Snowflake.WorkerID); err != nil {
		panic(err)
	}
}
