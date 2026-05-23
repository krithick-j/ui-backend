package configs

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	DBUserName      string `mapstructure:"MYSQL_USER"`
	DBUserPassword  string `mapstructure:"MYSQL_PASSWORD"`
	DBName          string `mapstructure:"MYSQL_DB"`
	DBHost          string `mapstructure:"MYSQL_HOST"`
	DBPort          int    `mapstructure:"MYSQL_PORT"`
	ChequeDrawValue int    `mapstructure:"CHEQUE_DRAW_VALUE"`
	AppPort         int    `mapstructure:"APP_PORT"`
	JWTSecret       string `mapstructure:"JWT_SECRET"`
}

var GlobalConfig *Config

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")
	viper.SetDefault("CHEQUE_DRAW_VALUE", 4000)
	viper.SetDefault("APP_PORT", 8000)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	required := map[string]string{
		"MYSQL_USER":     strings.TrimSpace(cfg.DBUserName),
		"MYSQL_PASSWORD": strings.TrimSpace(cfg.DBUserPassword),
		"MYSQL_DB":       strings.TrimSpace(cfg.DBName),
		"MYSQL_HOST":     strings.TrimSpace(cfg.DBHost),
		"JWT_SECRET":     strings.TrimSpace(cfg.JWTSecret),
	}
	for key, value := range required {
		if value == "" {
			return nil, fmt.Errorf("missing required environment variable: %s", key)
		}
	}
	if cfg.DBPort == 0 {
		return nil, fmt.Errorf("missing required environment variable: MYSQL_PORT")
	}

	GlobalConfig = cfg
	return GlobalConfig, nil
}
