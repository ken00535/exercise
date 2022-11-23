package gorm

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func New(cfg Config) *gorm.DB {
	config, err := cfg.clone()
	if err != nil {
		log.Panic().Msg(err.Error())
	}
	config.setDSN()
	config.setConnection()
	db, err := connectDB(config)
	if err != nil {
		log.Panic().Msg(err.Error())
	}
	return db
}
