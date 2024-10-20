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

package client

import (
	"context"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

// InitMySQL 通用初始化mysql函数，传入tableName指定表
func InitMySQL(tableName string) (db *gorm.DB, sf *utils.Snowflake, err error) {
	dsn, err := utils.GetMysqlDSN()
	if err != nil {
		return nil, nil, fmt.Errorf("dal.InitMySQL %s:get mysql DSN error: %v", tableName, err.Error())
	}
	db, err = gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
	if err != nil {
		return nil, nil, fmt.Errorf("dal.InitMySQL %s:mysql connect error: %v", tableName, err.Error())
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(constants.MaxIdleConns)
	sqlDB.SetMaxOpenConns(constants.MaxConnections)
	sqlDB.SetConnMaxLifetime(constants.ConnMaxLifetime)
	db = db.Table(tableName).WithContext(context.Background())

	if sf, err = utils.NewSnowflake(config.Snowflake.DatancenterID, config.Snowflake.WorkerID); err != nil {
		return nil, nil, fmt.Errorf("dal.InitMySQL %s:Snowflake init error: %v", tableName, err.Error())
	}
	return db, sf, nil
}
