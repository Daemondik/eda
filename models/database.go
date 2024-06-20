package models

import (
	zapLogger "eda/logger"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

var DB *gorm.DB

func ConnectDb() error {
	dsn := fmt.Sprintf(
		"host=db user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Europe/moscow",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		zapLogger.Log.Error("Failed to connect to database", zap.Error(err))
		return err
	}

	zapLogger.Log.Info("Successfully connected to the database")
	DB.Logger = logger.Default.LogMode(logger.Info)

	zapLogger.Log.Info("Running migrations")
	err = DB.AutoMigrate(&User{}, &Message{})
	if err != nil {
		zapLogger.Log.Error("Failed to migrate", zap.Error(err))
		return err
	}

	return nil
}
