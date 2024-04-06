package configs

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBUserName      string `mapstructure:"MYSQL_USER"`
	DBUserPassword  string `mapstructure:"MYSQL_PASSWORD"`
	DBName          string `mapstructure:"MYSQL_DB"`
	DBHost          string `mapstructure:"MYSQL_HOST"`
	DBPort          int    `mapstructure:"MYSQL_PORT"`
	ChequeDrawValue int    `mapstructure:"CHEQUE_DRAW_VALUE"`
}

var GlobalConfig Config

func LoadConfig(path string) (GlobalConfig *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")
	viper.SetDefault("CHEQUE_DRAW_VALUE", 4000)
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return 
	}

	// config.ChequeDrawValue = viper.GetInt("CHEQUE_DRAW_VALUE")

	err = viper.Unmarshal(&GlobalConfig)
	return
}
