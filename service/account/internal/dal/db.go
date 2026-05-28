package dal

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Config struct {
	DSN string
}

func defaultConfig() *Config {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		dsn = "root:root@tcp(127.0.0.1:3306)/opengeo?charset=utf8mb4&parseTime=True&loc=Local"
	}
	return &Config{DSN: dsn}
}

func Init(cfg *Config) error {
	if cfg == nil {
		cfg = defaultConfig()
	}

	var err error
	DB, err = gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("get sql db: %w", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := DB.AutoMigrate(
		&Tenant{}, &User{}, &Role{}, &Permission{}, &UserRole{}, &RolePermission{},
		&Account{}, &AccountHealth{}, &AlertRecord{},
		&BrowserFingerprint{}, &ProxyIP{}, &AccountEnvironment{},
		&AccountGroup{}, &AccountGroupRelation{},
	); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}

	return nil
}

func InitTestDB() error {
	DB = &gorm.DB{}
	return nil
}
