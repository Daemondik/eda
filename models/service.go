package models

import "eda/logger"

func InitializeServices() error {
	if err := logger.InitializeZapCustomLogger(); err != nil {
		return err
	}

	if err := ConnectDb(); err != nil {
		return err
	}

	if err := NewRedis(); err != nil {
		return err
	}

	return nil
}
