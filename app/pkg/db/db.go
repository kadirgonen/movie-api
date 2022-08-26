package db

import (
	"github.com/kadirgonen/movie-api/app/pkg/config"
	"gorm.io/gorm"
)

type DBSelector interface {
	Create(config *config.Config) (*gorm.DB, error)
}

type DBBase struct {
	DbType DBSelector
}
