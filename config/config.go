package config

import (
	"fmt"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App      `yaml:"app"`
		HTTP     `yaml:"http"`
		Log      `yaml:"logger"`
		PG       `yaml:"postgres"`
		Security `yaml:"security"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	PG struct {
		Username string `env-required:"true" yaml:"username" env:"PG_USERNAME"`
		Password string `env-required:"true" yaml:"password" env:"PG_PASSWORD"`
		Host     string `env-required:"true" yaml:"host" env:"PG_HOST"`
		Port     string `env-required:"true" yaml:"port" env:"PG_PORT"`
		DBName   string `env-required:"true" yaml:"db_name" env:"PG_DB_NAME"`
		SSLMode  string `env-required:"true" yaml:"ssl_mode" env:"PG_SSL_MODE"`
	}

	Security struct {
		APIKey string `env-required:"true" yaml:"api_key" env:"API_KEY"`
	}
)

var config *Config
var once sync.Once

func GetConfig() (*Config, error) {
	var err error
	once.Do(func() {
		config, err = newConfig()
	})
	if err != nil {
		return nil, err
	}
	return config, nil
}

func newConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
