package db

import (
	"context"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	logger2 "github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"gorm.io/gorm/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	DB *gorm.DB
)

func InitMySQL() {
	dsn, err := utils.GetMysqlDSN()
	if err != nil {

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
