package config

import (
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/go-playground/validator"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var ENVS = [][2]string{
	{"本地测试环境", "http://localhost:8080"},
}

type Config struct {
	Debug          bool `mapstructure:"debug"`
	HTTPPort       int  `mapstructure:"http_port"       validate:"required"`
	PrometheusPort int  `mapstructure:"prometheus_port" validate:"required"`

	Kafka struct {
		Endpoint string `mapstructure:"endpoint"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"kafka"`

	Databases struct {
		MySQL         string `mapstructure:"mysql"`
		Redis         string `mapstructure:"redis"`
		ElasticSearch struct {
			Endpoint string `mapstructure:"endpoint"`
			Username string `mapstructure:"username"`
			Password string `mapstructure:"password"`
		} `mapstructure:"elasticsearch"`
	} `mapstructure:"databases"`
	IsLocal bool
	Sentry  struct {
		Environment string
		DSN         string
	}
}

var (
	c    *Config
	once sync.Once
)

func Get() *Config {
	once.Do(setup)
	return c
}

func setup() {
	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(b), "../..")

	viper.SetConfigName("config")
	viper.AddConfigPath(projectRoot + "/configs")
	viper.AddConfigPath("/etc/{{ cookiecutter.project_slug }}/")
	viper.AddConfigPath(projectRoot)
	err := viper.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		log.Info().Msg("config files not found, ignoring")
	} else if err != nil {
		log.Warn().Err(err).Msg("read config failed")
	}

	// unmarshal it
	if err := viper.Unmarshal(&c); err != nil {
		log.Fatal().Err(err).Msg("unmarshal config failed")
	}

	if c.Debug {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	}

	// and validate it
	v := validator.New()
	if err := v.Struct(c); err != nil {
		log.Fatal().Err(err).Msg("validate config failed")
	}
	log.Info().Msg("load configs success")
}
