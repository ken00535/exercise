package config

import (
	"os"
	"strings"

	"assignment/internal/infra/gin"
	"assignment/internal/infra/zerolog"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var (
	_cfg Config
)

// Config is config instance
type Config struct {
	fx.Out
	Http *gin.Config    `yaml:"http"`
	Log  zerolog.Config `yaml:"log"`
}

// Init config singleton
func Init(configPath, filename string) {
	if path, find := os.LookupEnv("CONFIG_PATH"); find {
		configPath = path
	}
	if name, find := os.LookupEnv("CONFIG_FILENAME"); find {
		filename = name
	}
	viper.AutomaticEnv()
	viper.SetEnvPrefix("")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AddConfigPath(configPath)
	viper.SetConfigName(filename)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&_cfg); err != nil {
		panic(err)
	}
}

// Get get config instance
func Get() Config {
	return _cfg
}
