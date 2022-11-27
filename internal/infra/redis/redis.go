package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"strconv"
	"time"

	"github.com/cenk/backoff"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

type Config struct {
	Addr         string `yaml:"addr" mapstructure:"addr"`
	Port         int    `yaml:"port" mapstructure:"port"`
	Username     string `yaml:"username" mapstructure:"username"`
	Password     string `yaml:"password" mapstructure:"password"`
	Database     int    `yaml:"database" mapstructure:"database"`
	ReadTimeout  int    `yaml:"read_timeout" mapstructure:"read_timeout"`   // base on second
	WriteTimeout int    `yaml:"writ_etimeout" mapstructure:"writ_etimeout"` // base on second
	TLSDisable   bool   `yaml:"tls_disable" mapstructure:"tls_disable"`
}

func New(lc fx.Lifecycle, cfg Config) (*redis.Client, error) {

	readTimeout := time.Duration(cfg.ReadTimeout) * time.Second
	writeTimeout := time.Duration(cfg.WriteTimeout) * time.Second

	var tlsConfig *tls.Config
	if !cfg.TLSDisable {
		tlsConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}
	port := strconv.Itoa(cfg.Port)
	client := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    cfg.Addr + ":" + port,
		OnConnect: func(ctx context.Context, conn *redis.Conn) error {
			log.Info().Msgf("connect to redis success: %s", conn.String())
			return nil
		},
		Username:     cfg.Username,
		Password:     cfg.Password,
		DB:           cfg.Database,
		MaxRetries:   3,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		TLSConfig:    tlsConfig,
	})
	if err := checkConnection(client); err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Info().Msg("stopping redis connection")
			err := client.Close()
			return err
		},
	})
	return client, nil
}

func checkConnection(rdb *redis.Client) error {
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = time.Duration(180) * time.Second

	err := backoff.Retry(func() error {
		err := rdb.Ping(context.Background()).Err()
		if err != nil {
			log.Err(err).Msgf("redis client ping failed: %v", err)
			return err
		}
		return nil
	}, bo)

	if err != nil {
		return errors.WithStack(fmt.Errorf("connect redis failed: %v", err))
	}
	return nil
}
