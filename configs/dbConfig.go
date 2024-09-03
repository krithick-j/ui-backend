package configs

import (
	"context"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func DbConnect(config *Config) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.DBUserName, config.DBUserPassword, config.DBHost, config.DBPort, config.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("\x1b[31mFailed to connect to the Database!\x1b[0m", err.Error())
		os.Exit(1)
	}

	Log.Infoln(context.Background(), "\x1b[32m🚀Successfully connected to Database\b[0m")
	// db.AutoMigrate(&models.User{})
	DB = db
}
