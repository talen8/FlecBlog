package database

import (
	"fmt"
	"strings"
	"time"

	"flec_blog/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database 数据库实例
type Database struct {
	DB *gorm.DB
}

func (d *Database) GetDB() *gorm.DB {
	return d.DB
}

// escapeConnInfoValue 按 libpq 规则转义 keyword/value 连接串中的值
func escapeConnInfoValue(value string) string {
	return "'" + strings.NewReplacer(`\`, `\\`, `'`, `\'`).Replace(value) + "'"
}

// NewDB 创建数据库连接
func NewDB(config *config.DatabaseConfig) (*Database, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		escapeConnInfoValue(config.Host),
		escapeConnInfoValue(config.User),
		escapeConnInfoValue(config.Password),
		escapeConnInfoValue(config.DBName),
		config.Port)

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	return &Database{DB: db}, nil
}

// Close 关闭数据库连接
func (d *Database) Close() error {
	if d.DB != nil {
		sqlDB, err := d.DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
