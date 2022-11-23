package gorm

import (
	"fmt"
	"strings"
	"time"

	"github.com/cenk/backoff"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseDriver string

const (
	Postgres       DatabaseDriver = "postgres"
	ConnectTimeout int            = 360
)

type Config struct {
	Host               string `yaml:"host" mapstructure:"host"`
	Port               int32  `yaml:"port" mapstructure:"port"`
	Username           string `yaml:"username" mapstructure:"username"`
	Password           string `yaml:"password" mapstructure:"password"`
	DBName             string `yaml:"dbname" mapstructure:"dbname"`
	Schema             string `yaml:"schema" mapstructure:"schema"`
	SSLEnable          bool   `yaml:"ssl" mapstructure:"ssl"`
	MaxIdleConns       int    `yaml:"max_idle_conns" mapstructure:"max_idle_conns"`
	MaxOpenConns       int    `yaml:"max_open_conns" mapstructure:"max_open_conns"`
	ConnMaxLifeTimeSec int    `yaml:"conn_max_life_time_sec" mapstructure:"conn_max_life_time_sec"`
	PrepareStmt        *bool  `yaml:"prepare_stmt" mapstructure:"prepare_stmt"`
	dsn                string
}

func (cfg *Config) clone() (*Config, error) {
	var config Config
	config = *cfg
	return &config, nil
}

func (cfg *Config) setConnection() {
	if cfg.MaxIdleConns == 0 {
		cfg.MaxIdleConns = 50
	}
	if cfg.MaxOpenConns == 0 {
		cfg.MaxOpenConns = 100
	}
	if cfg.ConnMaxLifeTimeSec == 0 {
		cfg.ConnMaxLifeTimeSec = ConnectTimeout - 20
	}
}

func (cfg *Config) setDSN() {
	dsn := fmt.Sprintf(`user=%s password=%s host=%s port=%d dbname=%s connect_timeout=%d`,
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		ConnectTimeout)
	if cfg.SSLEnable {
		dsn += " sslmode=require"
	} else {
		dsn += " sslmode=disable"
	}
	if strings.TrimSpace(cfg.Schema) != "" {
		dsn = fmt.Sprintf("%s search_path=%s", dsn, cfg.Schema)
	}
	cfg.dsn = dsn
}

func connectDB(cfg *Config) (*gorm.DB, error) {
	var err error
	dialector := postgres.Open(cfg.dsn)

	var isPrepareStmt bool
	if cfg.PrepareStmt != nil {
		isPrepareStmt = *cfg.PrepareStmt
	}

	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = time.Duration(180) * time.Second
	var db *gorm.DB

	err = backoff.Retry(func() error {
		db, err = gorm.Open(dialector, &gorm.Config{
			PrepareStmt: isPrepareStmt,
		})
		if err != nil {
			log.Printf("gorm open failed: %v", err)
			return err
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("get DB failed: %v", err)
			return err
		}

		err = sqlDB.Ping()
		return err
	}, bo)

	if err != nil {
		return nil, errors.WithStack(fmt.Errorf("connect db failed: %v", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifeTimeSec) * time.Second)
	log.Info().Msgf("connect to db success: %v:%v, dbname: %v", cfg.Host, cfg.Port, cfg.DBName)

	return db, nil
}
