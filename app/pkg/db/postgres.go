package db

import (
	"fmt"
	"sync"
	"time"

	"github.com/kadirgonen/movie-api/app/pkg/config"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type POSTGRES struct{}

var (
	once sync.Once
	conn *gorm.DB
)

func (ps *POSTGRES) Create(config *config.Config) (*gorm.DB, error) {
	once.Do(func() {
		res, err := CreatePostgreSQL(config)
		conn = res
		if err != nil {
			zap.L().Fatal("Error while db Creating", zap.Error(err))
		}
	})
	return conn, nil
}

func CreatePostgreSQL(config *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.DBConfig.DataSourceName), &gorm.Config{})
	if err != nil {
		zap.L().Fatal("Cannot connect to database", zap.Error(err))
		return nil, fmt.Errorf("Cannot connect to database: %v", err)
	}
	origin, err := db.DB()
	if err != nil {
		zap.L().Fatal("Cannot get sql.DB from database", zap.Error(err))
		return nil, err
	}
	if err := origin.Ping(); err != nil {
		zap.L().Fatal("Cannot ping to the db", zap.Error(err))
		return nil, err
	}

	origin.SetMaxOpenConns(config.DBConfig.MaxOpen)
	origin.SetMaxIdleConns(config.DBConfig.MaxIdle)
	origin.SetConnMaxLifetime(time.Duration(config.DBConfig.MaxLifetime) * time.Second)

	return db, nil
}
