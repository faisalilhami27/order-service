package config

import (
	"fmt"
	"net/url"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"time"
)

func InitDatabase() (*gorm.DB, error) {
	cfg := Config
	encodedPassword := url.QueryEscape(cfg.Database.Password)
	uri := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Database.Username,
		encodedPassword,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)
	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConnection)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConnection)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.MaxLifetimeConnection) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.Database.MaxIdleTime) * time.Second)
	return db, nil
}
