package config

import (
	"bytes"
	"strings"

	"shorten/internal/app/usecase"
	"shorten/internal/infra/gin"
	"shorten/internal/infra/gorm"
	"shorten/internal/infra/zerolog"

	"github.com/spf13/viper"
	"go.uber.org/fx"
	"gopkg.in/yaml.v3"
)

var (
	_cfg Config
)

// Config is config instance
type Config struct {
	fx.Out
	Us   usecase.Config `yaml:"shorten" mapstructure:"shorten"`
	Gorm gorm.Config    `yaml:"db" mapstructure:"db"`
	Http gin.Config     `yaml:"http" mapstructure:"http"`
	Log  zerolog.Config `yaml:"log" mapstructure:"log"`
}

// Init config singleton
func Init() {
	b, _ := yaml.Marshal(&Config{})
	defaultConfig := bytes.NewReader(b)
	viper.SetConfigType("yaml")
	_ = viper.MergeConfig(defaultConfig)
	viper.AutomaticEnv()
	viper.SetEnvPrefix("")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.Unmarshal(&_cfg); err != nil {
		panic(err)
	}
}

// Get get config instance
func Get() Config {
	return _cfg
}
