package database

import (
	"fmt"
	"os"
	"time"

	"opengeo/pkg/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Config 数据库配置
type Config struct {
	DSN             string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

// DefaultConfig 默认配置
func DefaultConfig(envPrefix string) *Config {
	dsn := os.Getenv(envPrefix + "_MYSQL_DSN")
	if dsn == "" {
		dsn = os.Getenv("MYSQL_DSN")
	}
	if dsn == "" {
		dsn = "root:root@tcp(127.0.0.1:3306)/opengeo?charset=utf8mb4&parseTime=True&loc=Local"
	}

	cfg := config.GetConfig()
	return &Config{
		DSN:             dsn,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
	}
}

// New 创建数据库连接
func New(cfg *Config) (*gorm.DB, error) {
	if cfg == nil {
		cfg = DefaultConfig("")
	}

	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return db, nil
}

// MustNew 创建数据库连接（失败则 panic）
func MustNew(cfg *Config) *gorm.DB {
	db, err := New(cfg)
	if err != nil {
		panic(fmt.Sprintf("database connection failed: %v", err))
	}
	return db
}
