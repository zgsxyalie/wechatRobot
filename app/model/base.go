package model

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"path/filepath"
	"wechatRobot/app/helper"
)

var DB *gorm.DB

func init() {
	dbname := "./storage/data/wechat.db3"
	if !helper.Exists(dbname) {
		_ = os.MkdirAll(filepath.Dir(dbname), os.ModePerm)
	}

	db, err := gorm.Open(sqlite.Open(dbname), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatal("数据库连接失败：", err)
	}

	// 迁移 schema
	_ = db.AutoMigrate(&Contact{}, &ChatRoom{}, &Message{})

	// 开启外键约束
	db.Exec("PRAGMA foreign_keys=ON;")

	DB = db
}
