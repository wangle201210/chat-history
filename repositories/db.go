package repositories

import (
	"fmt"
	"log"
	"time"

	"github.com/wangle201210/chat-history/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// InitDB 初始化数据库连接
func InitDB(dsn string) error {
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	var err error
	db, err = gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		return fmt.Errorf("failed to connect database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %v", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移数据库表结构
	if err = autoMigrateTables(); err != nil {
		return fmt.Errorf("failed to migrate database tables: %v", err)
	}

	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	if db == nil {
		log.Fatal("database connection not initialized")
	}
	return db
}

// autoMigrateTables 自动迁移数据库表结构
func autoMigrateTables() error {
	// 自动迁移会创建表、缺失的外键、约束、列和索引
	return db.AutoMigrate(
		&models.Conversation{},
		&models.Message{},
		&models.Attachment{},
		&models.MessageAttachment{},
	)
}
