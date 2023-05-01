package models

import (
	"book_management_system_backend/config"
	"book_management_system_backend/utils"
	"errors"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

var DB *gorm.DB

var LockClause = clause.Locking{Strength: "UPDATE"}

var gormConfig = &gorm.Config{
	NamingStrategy: schema.NamingStrategy{
		SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
	},
	Logger: logger.New(
		zap.NewStdLog(utils.Logger),
		logger.Config{
			SlowThreshold:             time.Second,  // 慢 SQL 阈值
			LogLevel:                  logger.Error, // 日志级别
			IgnoreRecordNotFoundError: true,         // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,        // 禁用彩色打印
		},
	),
}

func InitDB() {
	var err error
	DB, err = gorm.Open(postgres.Open(config.Config.PostgresDSN.String()), gormConfig)
	if err != nil {
		panic(err)
	}

	err = DB.AutoMigrate(User{}, Book{}, UserJwtSecret{})
	if err != nil {
		panic(err)
	}

	if config.Config.Debug {
		DB = DB.Debug()
	}

	utils.Logger.Info("database connected")

	// initialize admin user
	var firstUser User
	err = DB.Take(&firstUser, 1).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			panic(err)
		}

		firstUser = User{
			Username:       "admin",
			HashedPassword: utils.MakePassword("admin"),
			IsAdmin:        true,
		}

		err = DB.Create(&firstUser).Error
	}
}
