package common

import (
	"github.com/glebarez/sqlite"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

var DB *gorm.DB

func InitDb(dsn string) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // 使用一个新的 logger
			logger.Config{
				SlowThreshold:             0,           // 不显示慢查询日志
				LogLevel:                  logger.Error, // 设置日志级别为 Info，只输出 Info 级别以上的日志
				IgnoreRecordNotFoundError: true,        // 忽略记录未找到的错误
				Colorful:                  true,       // 不使用彩色输出
			},
		),
	})
	if err == nil {
		DB = db
		err = db.AutoMigrate(&Model{})
		if err != nil {
			logrus.Error("db auto migrate fail")
		}
	} else {
		logrus.Error("db connection err: " + err.Error())
	}
}

func CloseDb() {
	sqlDB, err := DB.DB()
	if err != nil {
		logrus.Error("fetch db fail")
		return
	}

	err = sqlDB.Close()
	if err != nil {
		logrus.Error("db close error: " + err.Error())
		return
	}
}
