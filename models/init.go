package models

import (
	"book_management_system_backend/config"
	"book_management_system_backend/utils"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
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
		utils.StdOutLogger,
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
	switch config.Config.Mode {
	case config.ModeTest:
		fallthrough
	case config.ModeBench:
		DB, err = gorm.Open(sqlite.Open("file::memory:"), gormConfig)
	case config.ModeDev:
		DB, err = gorm.Open(sqlite.Open("data.db"), gormConfig)
	case config.ModeProduction:
		DB, err = gorm.Open(postgres.Open(config.Config.PostgresDSN.String()), gormConfig)
	default:
		panic("unknown mode")
	}
	if err != nil {
		panic(err)
	}

	err = DB.AutoMigrate(User{}, Book{}, UserJwtSecret{}, Balance{}, Purchase{}, Sale{})
	if err != nil {
		panic(err)
	}

	if config.Config.Debug || config.Config.Mode == config.ModeTest {
		DB = DB.Debug()
	}

	utils.Logger.Info("database connected")

	// initialize admin user
	var firstUser User
	err = DB.Where(User{ID: 1}).Attrs(User{
		Username:       "admin",
		HashedPassword: utils.MakePassword("adminadmin"),
		IsAdmin:        true,
	}).FirstOrCreate(&firstUser).Error
	if err != nil {
		panic(err)
	}

	// initialize balance
	var firstBalance Balance
	err = DB.Where(Balance{ID: 1}).Attrs(Balance{
		UserID:        firstUser.ID,
		OperationType: OperationTypeInitialize,
	}).FirstOrCreate(&firstBalance).Error
	if err != nil {
		panic(err)
	}
}
