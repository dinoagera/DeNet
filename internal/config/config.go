package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DBUser          string        `env:"DB_USER" env-required:"true"`
	DBPassword      string        `env:"DB_PASSWORD" env-required:"true"`
	DBName          string        `env:"DB_NAME" env-required:"true"`
	DBPath          string        `env:"DB_PATH" env-required:"true"`
	HTTPAddress     string        `env:"HTTPAddress"`
	HTTPIdleTimeout time.Duration `env:"HTTPIdleTimeout"`
	HTTPReadTimeout time.Duration `env:"HTTPReadTimeout"`
	SecretKey       string        `env:"SecretKey"`
	TTL             time.Duration `env:"TTL"`
}

var cfg *Config

func InitConfig(log *slog.Logger) *Config {
	pathCfg := ".env"
	// if _, err := os.Stat(pathCfg); os.IsNotExist(err) {
	// 	log.Info("not found config", "error", err)
	// 	os.Exit(1)
	// }
	var localconfig Config
	err := cleanenv.ReadConfig(pathCfg, &localconfig)
	if err != nil {
		log.Info("init config to failed", "error", err)
		os.Exit(1)
	}
	cfg = &localconfig
	return cfg
}

func GetConfig() *Config {
	if cfg == nil {
		panic("config is not init")
	}
	return cfg
}
