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
	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/base/client"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

var (
	DB *gorm.DB
	SF *utils.Snowflake
)

func InitMySQL() {
	Db, err := client.InitMySQL(constants.LaunchScreenTableName)
	if err != nil {
		logger.Fatal(err)
	}
	sf, err := utils.NewSnowflake(config.Snowflake.DatancenterID, config.Snowflake.WorkerID)
	if err != nil {
		logger.Fatalf("init Snowflake object error: %v", err.Error())
	}

	DB = Db
	SF = sf
}